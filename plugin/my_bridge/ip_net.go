package my_bridge

import (
	"net"
	"strings"
)

type IPNet net.IPNet

func ParseIPNet(s string) (*IPNet, error) {
	ip, netIpNet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}

	netIpNet.IP = ip

	ipNet := IPNet(*netIpNet)

	return &ipNet, nil
}

func (i *IPNet) String() string {
	ipNet := net.IPNet(*i)

	return ipNet.String()
}

func (i *IPNet) UnmarshalJSON(bytes []byte) error {
	ip, ipNet, err := net.ParseCIDR(strings.ReplaceAll(string(bytes), "\"", ""))
	if err != nil {
		return err
	}

	*i = IPNet(*ipNet)
	i.IP = ip

	return nil
}

func (i *IPNet) MarshalJSON() ([]byte, error) {
	ipNet := net.IPNet(*i)

	return []byte("\"" + ipNet.String() + "\""), nil
}
