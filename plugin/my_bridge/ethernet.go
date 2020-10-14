package my_bridge

import (
	"net"
	"os"

	"github.com/vishvananda/netlink"
)

func CreateBridge(name string) (netlink.Link, error) {
	linkAttrs := netlink.LinkAttrs{
		Name: name,
	}

	bridge := &netlink.Bridge{
		LinkAttrs: linkAttrs,
	}

	if err := netlink.LinkAdd(bridge); err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err
	}

	return link, nil
}

/*func AssociateIpToBridge(bridge netlink.Link, ip *net.IPNet) error {
	addr := &netlink.Addr{
		IPNet: ip,
		Label: "",
	}

	return netlink.AddrAdd(bridge, addr)
}*/

func CreateVethPair(name1, name2 string) (veth, peer netlink.Link, err error) {
	attrs := netlink.LinkAttrs{
		Name: name1,
	}

	veth = &netlink.Veth{
		LinkAttrs: attrs,
		PeerName:  name2,
	}

	if err := netlink.LinkAdd(veth); err != nil {
		if !os.IsExist(err) {
			return nil, nil, err
		}
	}

	veth, err = netlink.LinkByName(name1)
	if err != nil {
		return nil, nil, err
	}

	peer, err = netlink.LinkByName(name2)
	if err != nil {
		return nil, nil, err
	}

	return
}

func AssociateAddrAndSetRoute(link netlink.Link, addr net.IP, mask net.IPMask, gateway *net.IP) error {
	if err := netlink.AddrAdd(link, &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   addr,
			Mask: mask,
		},
	}); err != nil {
		return err
	}

	if gateway != nil {
		return netlink.RouteAdd(&netlink.Route{
			Dst: &net.IPNet{
				IP:   net.ParseIP("0.0.0.0").To4(),
				Mask: net.CIDRMask(0, 32),
			},
			Gw: *gateway,
		})
	}

	return nil
}
