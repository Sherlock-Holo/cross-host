package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sherlock-Holo/cross-host/plugin/my_bridge"
)

func HandleDelete() error {
	containerId := os.Getenv("CNI_CONTAINERID")
	if containerId == "" {
		return NoContainerId
	}

	netns := os.Getenv("CNI_NETNS")
	if netns == "" {
		return NoNetns
	}

	ifaceName := os.Getenv("CNI_IFNAME")
	if ifaceName == "" {
		return NoIfaceName
	}

	extraArgs := make(map[string]string)

	if args, ok := os.LookupEnv("CNI_ARGS"); ok {
		for _, arg := range strings.Split(args, ";") {
			splitArg := strings.Split(arg, "=")
			if len(splitArg) != 2 {
				extraArgs[splitArg[0]] = splitArg[1]
			}
		}
	}

	return handleDelete(containerId, netns, ifaceName, my_bridge.ConfigDir, extraArgs, os.Stdin)
}

func handleDelete(
	containerId, netns, ifaceName, cfgDir string,
	extraArgs map[string]string,
	cniCfgReader io.Reader,
) error {
	var (
		cniCfg       my_bridge.CniConfig
		containerCfg my_bridge.ContainerConfig
		pluginCfg    my_bridge.PluginConfig
	)

	if err := json.NewDecoder(cniCfgReader).Decode(&cniCfg); err != nil {
		return err
	}

	containerCfgPath := filepath.Join(cfgDir, cniCfg.Name, containerId+".json")
	defer func() {
		_ = os.Remove(containerCfgPath)
	}()

	containerCfgFile, err := os.Open(containerCfgPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	defer func() {
		_ = containerCfgFile.Close()
	}()

	if err := json.NewDecoder(containerCfgFile).Decode(&containerCfg); err != nil {
		return err
	}

	if containerCfg.Bridge != cniCfg.Name {
		return fmt.Errorf("container bridge is %s, but delete argument bridge is %s", containerCfg.Bridge, cniCfg.Name)
	}

	pluginCfgPath := filepath.Join(cfgDir, cniCfg.Name, my_bridge.ConfigFilename)

	pluginCfgFile, err := os.OpenFile(pluginCfgPath, os.O_RDWR, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	defer func() {
		_ = pluginCfgFile.Close()
	}()

	if err := json.NewDecoder(pluginCfgFile).Decode(&pluginCfg); err != nil {
		return err
	}

	containerIp := containerCfg.IP.To4()

	index := binary.BigEndian.Uint16(containerIp[2:])

	if !my_bridge.Release(&pluginCfg.Bitmap, int(index)) {
		return fmt.Errorf("container IP %s is not used", containerIp)
	}

	if _, err := pluginCfgFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if err := json.NewEncoder(pluginCfgFile).Encode(pluginCfg); err != nil {
		return err
	}

	return nil
}
