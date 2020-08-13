package config

import (
	"gopkg.in/yaml.v2"
)

var data = `
bind_ip: 172.16.100.2	
gateway: 10.65.101.1
netmask: 255.255.255.0
dns: 10.192.2.10
nics:
  00-50-56-82-70-2a:
    dhcp_ip: 172.16.102.100
    ip: 10.65.102.2
    hostname: vc-01.example.org
  00-50-56-82-60-2a:
    dhcp_ip: 172.16.102.101
    ip: 10.65.102.3  
    hostname: vc-02.example.org
boot_file: efi/boot/bootx64.efi
lease_time: 500
root_path: /home/ubuntu/images
ntp_server: time.svc.pivotal.io
`

type ServerConfig struct {
	Ip string
	DhcpIp string `yaml:"dhcp_ip"`
	Hostname string
}
type Config struct {
	BindIP string `yaml:"bind_ip"`
	DNS string `yaml:"dns"`
	Gateway string `yaml:"gateway"`
	NTPServer string `yaml:"ntp_server"`
	Netmask string `yaml:"netmask"`
	Nics map[string]ServerConfig `yaml:"nics"`
	BootFile string `yaml:"boot_file"`
	LeaseTime int `yaml:"lease_time"`
	RootPath string `yaml:"root_path"`
}

func LoadConfig() (*Config, error){
	config := Config{}
	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}
	return &config, err
}
