package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

var data = `
dhcp_interface: ens224
bind_ip: 172.16.100.2	
gateway: 10.65.101.1
netmask: 255.255.255.0
dns: 10.192.2.10
nics:
  00-50-56-82-70-2a:
    dhcp_ip: 172.16.100.100
    ip: 10.65.101.10
    hostname: vc-01.example.org
  00-50-56-82-61-7c:
    dhcp_ip: 172.16.100.101
    ip: 10.65.101.11  
    hostname: vc-02.example.org
  00-50-56-82-d8-86:
    dhcp_ip: 172.16.100.102
    ip: 10.65.101.12
    hostname: vc-03.example.org
boot_file: efi/boot/bootx64.efi
lease_time: 500
root_path: /home/ubuntu/images
boot_config_file: efi/boot/boot.cfg
ntp_server: time.svc.pivotal.io
`

type ServerConfig struct {
	Ip string `yaml:"ip" json:"ip"`
	DhcpIp string `yaml:"dhcp_ip" json:"dhcp_ip"`
	Hostname string `yaml:"hostname" json:"hostname"`
}
type Config struct {
	BindIP string `yaml:"bind_ip" json:"bind_ip"`
	ManagementIp string `yaml:"management_ip" json:"management_ip"`
	DNS string `yaml:"dns" json:"dns"`
	DHCPServerPort int `yaml:"dhcp_server_port" json:"dhcp_server_port"`
	TFTPServerPort int `yaml:"tftp_server_port" json:"tftp_server_port"`
	HTTPPort int `yaml:"http_port" json:"http_port"`
	Gateway string `yaml:"gateway" json:"gateway"`
	NTPServer string `yaml:"ntp_server" json:"ntp_server"`
	Netmask string `yaml:"netmask" json:"netmask"`
	Nics map[string]ServerConfig `yaml:"nics" json:"nics"`
	BootFile string `yaml:"boot_file" json:"boot_file"`
	LeaseTime int `yaml:"lease_time" json:"lease_time"`
	RootPath string `yaml:"root_path" json:"root_path"`
	DHCPInterface string `yaml:"dhcp_interface" json:"dhcp_interface"`
	Password string `yaml:"password" json:"password"`
	BootConfigFile string `yaml:"boot_config_file" json:"boot_config_file"`
	KickStartTemplate string `yaml:"kickstart_template" json:"kickstart_template"`
}

func LoadConfig(reader io.Reader) (*Config, error){
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, err
}