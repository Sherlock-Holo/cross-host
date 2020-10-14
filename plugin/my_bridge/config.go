package my_bridge

import (
	"net"
)

const ConfigDir = "/run/cross-host"
const ConfigFilename = "config.json"

type PluginConfig struct {
	IP     *IPNet
	Bitmap [8192]byte // support xxx.xxx.0.0/16
}

type ContainerConfig struct {
	IP     net.IP
	Bridge string
}

type CniConfig struct {
	CniVersion string                 `json:"cniVersion"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Args       map[string]interface{} `json:"args"`
	IpMasq     bool                   `json:"ipMasq"`
	Ipam       *CniIpam               `json:"ipam"`
	Dns        *CniDns                `json:"dns"`
}

type CniDns struct {
	Nameservers []string `json:"nameservers"`
	Domain      string   `json:"domain"`
	Search      []string `json:"search"`
	Options     []string `json:"options"`
}

type CniIpam struct {
	Type    string     `json:"type"`
	Subnet  *IPNet     `json:"subnet"`
	Gateway net.IP     `json:"gateway"`
	Routes  []CniRoute `json:"routes"`
}

type CniRoute struct {
	Dst *IPNet `json:"dst"`
	Gw  net.IP `json:"gw"`
}
