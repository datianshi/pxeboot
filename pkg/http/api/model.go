package api

import (
	"errors"
	"net"
)

type ServerItem struct {
	Ip string `yaml:"ip" json:"ip"`
	DhcpIp string `yaml:"dhcp_ip" json:"dhcp_ip"`
	Hostname string `yaml:"hostname" json:"hostname"`
	MacAddress string `json:"mac_address"`
}

func (s ServerItem) Validate() error {
	if _, err := net.ParseMAC(s.MacAddress); err != nil {
		return errors.New("error : Mac address format is not valid")
	}
	if ip := net.ParseIP(s.Ip); ip == nil {
		return errors.New("IP address is not valid")
	}
	if ip := net.ParseIP(s.DhcpIp); ip == nil {
		return errors.New("DHCP assigned IP address is not valid")
	}
	return nil
}