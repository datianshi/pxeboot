package config

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
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

//Config configuration for pxeboot server
type Config struct {
	BindIP              string   `yaml:"bind_ip" json:"bind_ip"`
	ManagementInterface string   `yaml:"management_interface" json:"management_interface"`
	DNS                 string   `yaml:"dns" json:"dns"`
	DHCPRange           string   `yaml:"dhcp_range" json:"dhcp_range"`
	DHCPServerPort      int      `yaml:"dhcp_server_port" json:"dhcp_server_port"`
	TFTPServerPort      int      `yaml:"tftp_server_port" json:"tftp_server_port"`
	HTTPPort            int      `yaml:"http_port" json:"http_port"`
	NTPServer           string   `yaml:"ntp_server" json:"ntp_server"`
	BootFile            string   `yaml:"boot_file" json:"boot_file"`
	LeaseTime           int      `yaml:"lease_time" json:"lease_time"`
	RootPath            string   `yaml:"root_path" json:"root_path"`
	DHCPInterface       string   `yaml:"dhcp_interface" json:"dhcp_interface"`
	Password            string   `yaml:"password" json:"password"`
	BootConfigFile      string   `yaml:"boot_config_file" json:"boot_config_file"`
	KickStartTemplate   string   `yaml:"kickstart_template" json:"kickstart_template"`
	Database            Database `yaml:"database" json:"database"`
}

//Database Database Config
type Database struct {
	Username     string `yaml:"username" json:"username"`
	Password     string `yaml:"password" json:"password"`
	Host         string `yaml:"host" json:"host"`
	Port         int    `yaml:"port" json:"port"`
	DatabaseName string `yaml:"database_name" json:"database_name"`
}

//LoadConfig load config from io.Reader
func LoadConfig(reader io.Reader) (Config, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, err
}
