package model

import (
	"errors"
	"net"

	"github.com/spf13/pflag"
)

//ServerConfig Server config
type ServerConfig struct {
	ID         int    `yaml:"id" json:"id"`
	Ip         string `yaml:"ip" json:"ip"`
	Hostname   string `yaml:"hostname" json:"hostname"`
	Gateway    string `yaml:"gateway" json:"gateway"`
	Netmask    string `yaml:"netmask" json:"netmask"`
	MacAddress string `yaml:"mac_address" json:"mac_address"`
}

//Validate Validate
func (s ServerConfig) Validate() error {
	if _, err := net.ParseMAC(s.MacAddress); err != nil {
		return errors.New("error : Mac address format is not valid")
	}
	if ip := net.ParseIP(s.Ip); ip == nil {
		return errors.New("IP address is not valid")
	}
	if ip := net.ParseIP(s.Gateway); ip == nil {
		return errors.New("the gateway is not valid")
	}
	if mask := pflag.ParseIPv4Mask(s.Netmask); mask == nil {
		return errors.New("net Mask is not valid")
	}
	return nil
}
