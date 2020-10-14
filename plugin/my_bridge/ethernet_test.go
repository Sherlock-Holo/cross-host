package my_bridge

import (
	"net"
	"testing"

	"github.com/vishvananda/netlink"
)

func TestCreateBridge(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantErr  bool
	}{
		{
			name:     "test1",
			args:     args{name: "test1"},
			wantName: "test1",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateBridge(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBridge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.Attrs().Name != tt.wantName {
				t.Errorf("CreateBridge() got = %v, want name %s", got, tt.wantName)
			}
		})
	}
}

/*func TestAssociateIpToBridge(t *testing.T) {
	_, ipNet1, err := net.ParseCIDR("192.168.1.1/24")
	if err != nil {
		t.Fatal(err)
	}

	_, ipNet2, err := net.ParseCIDR("10.200.0.1/16")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		ip *net.IPNet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				ip: ipNet1,
			},
		},
		{
			name: "test2",
			args: args{
				ip: ipNet2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bridge, err := CreateBridge(tt.name)
			if err != nil {
				t.Fatalf("create bridge %s failed: %+v", tt.name, err)
			}

			if err := AssociateIpToBridge(bridge, tt.args.ip); (err != nil) != tt.wantErr {
				t.Errorf("AssociateIpToBridge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}*/

func TestCreateVethPair(t *testing.T) {
	type args struct {
		name1 string
		name2 string
	}
	tests := []struct {
		name         string
		args         args
		wantVethName string
		wantPeerName string
		wantErr      bool
	}{
		{
			name: "test1-test2",
			args: args{
				name1: "test1",
				name2: "test2",
			},
			wantVethName: "test1",
			wantPeerName: "test2",
			wantErr:      false,
		},
		{
			name: "test11-test22",
			args: args{
				name1: "test11",
				name2: "test22",
			},
			wantVethName: "test11",
			wantPeerName: "test22",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVeth, gotPeer, err := CreateVethPair(tt.args.name1, tt.args.name2)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateVethPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotVeth.Attrs().Name != tt.wantVethName {
				t.Errorf("CreateVethPair() gotVeth name = %v, want %v", gotVeth.Attrs().Name, tt.wantVethName)
			}

			if gotPeer.Attrs().Name != tt.wantPeerName {
				t.Errorf("CreateVethPair() gotPeer name = %v, want %v", gotPeer.Attrs().Name, tt.wantPeerName)
			}
		})
	}
}

func TestAssociateAddrAndSetRoute(t *testing.T) {
	bridge, err := CreateBridge("test1")
	if err != nil {
		t.Fatal(err)
	}

	if err := netlink.LinkSetUp(bridge); err != nil {
		t.Fatal(err)
	}

	veth, _, err := CreateVethPair("veth1", "veth2")
	if err != nil {
		t.Fatal(err)
	}

	if err := netlink.LinkSetUp(veth); err != nil {
		t.Fatal(err)
	}

	gateway := net.ParseIP("10.10.1.1")

	type args struct {
		link    netlink.Link
		addr    net.IP
		mask    net.IPMask
		gateway *net.IP
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "associate bridge",
			args: args{
				link:    bridge,
				addr:    net.ParseIP("192.168.1.1"),
				mask:    net.CIDRMask(24, 32),
				gateway: nil,
			},
		},
		{
			name: "associate veth",
			args: args{
				link:    veth,
				addr:    net.ParseIP("10.10.1.10"),
				mask:    net.CIDRMask(16, 32),
				gateway: &gateway,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AssociateAddrAndSetRoute(tt.args.link, tt.args.addr, tt.args.mask, tt.args.gateway); (err != nil) != tt.wantErr {
				t.Errorf("AssociateAddrAndSetRoute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
