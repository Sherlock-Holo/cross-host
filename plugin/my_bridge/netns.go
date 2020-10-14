package my_bridge

import (
	"os"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func JoinNetns(veth netlink.Link, netnsPath string) (currentNetns netns.NsHandle, err error) {
	netnsFile, err := os.Open(netnsPath)
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = netnsFile.Close()
	}()

	netnsFd := int(netnsFile.Fd())

	if err := netlink.LinkSetNsFd(veth, netnsFd); err != nil {
		return 0, err
	}

	currentNetns, err = netns.Get()
	if err != nil {
		return 0, err
	}

	if err := netns.Set(netns.NsHandle(netnsFd)); err != nil {
		return 0, err
	}

	return currentNetns, nil
}

func ReturnOriginNetns(origin netns.NsHandle) error {
	return netns.Set(origin)
}
