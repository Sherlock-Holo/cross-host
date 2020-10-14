package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/Sherlock-Holo/cross-host/plugin/my_bridge"
)

func prepareSuccess() (successTempDir string, successCniCfg []byte, err error) {
	successTempDir, err = ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		return "", nil, err
	}

	successCniCfg, err = json.Marshal(my_bridge.CniConfig{
		CniVersion: currentVersion,
		Name:       "success",
		Type:       "my_bridge",
		Args:       nil,
		IpMasq:     false,
		Ipam:       nil,
		Dns:        nil,
	})
	if err != nil {
		return "", nil, err
	}

	if err := os.MkdirAll(filepath.Join(successTempDir, "success"), 0755); err != nil {
		return "", nil, err
	}

	containerIp := net.ParseIP("192.168.1.1")

	containerCfg, err := json.Marshal(my_bridge.ContainerConfig{
		IP:     containerIp,
		Bridge: "success",
	})
	if err != nil {
		return "", nil, err
	}

	if err := ioutil.WriteFile(filepath.Join(successTempDir, "success", "success-id.json"), containerCfg, 0644); err != nil {
		return "", nil, err
	}

	ipNet, err := my_bridge.ParseIPNet("192.168.0.0/16")
	if err != nil {
		return "", nil, err
	}

	pluginCfg := my_bridge.PluginConfig{
		IP:     ipNet,
		Bitmap: [8192]byte{},
	}

	index := int(binary.BigEndian.Uint16(containerIp.To4()[2:]))

	for i := 0; i <= index; i++ {
		my_bridge.Acquire(&pluginCfg.Bitmap)
	}

	pluginCfgJsonData, err := json.Marshal(pluginCfg)
	if err != nil {
		return "", nil, err
	}

	if err := ioutil.WriteFile(filepath.Join(successTempDir, "success", my_bridge.ConfigFilename), pluginCfgJsonData, 0644); err != nil {
		return "", nil, err
	}

	return successTempDir, successCniCfg, nil
}

func Test_handleDelete(t *testing.T) {
	successTempDir, successCniCfg, err := prepareSuccess()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(successTempDir)
	}()

	type args struct {
		containerId  string
		netns        string
		ifaceName    string
		cfgDir       string
		extraArgs    map[string]string
		cniCfgReader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty cni cfg reader",
			args: args{
				cniCfgReader: bytes.NewReader(nil),
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				containerId:  "success-id",
				netns:        "test-netns",
				ifaceName:    "test-iface",
				cfgDir:       successTempDir,
				extraArgs:    nil,
				cniCfgReader: bytes.NewReader(successCniCfg),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := handleDelete(tt.args.containerId, tt.args.netns, tt.args.ifaceName, tt.args.cfgDir, tt.args.extraArgs, tt.args.cniCfgReader); (err != nil) != tt.wantErr {
				t.Errorf("handleDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
