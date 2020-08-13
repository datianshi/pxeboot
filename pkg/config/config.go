package config

import (
	"fmt"
	"github.com/datianshi/pxeboot/pkg/util"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	DHCPInterface string `yaml:"dhcp_interface"`
	Password string `yaml:"password"`
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

//SetupPxeDirectory Each nic will have a directory
func (cfg *Config) SetupPxeDirectory() {
	//Create directories for each nic
	for k, _ := range cfg.Nics {
		//$root_path/01-nic_mac_address
		serverDir := fmt.Sprintf("%s/01-%s", cfg.RootPath, k)
		err := os.Mkdir(serverDir, 0755)
		if err != nil {
			log.Println(err)
		}
		//$root_path/01-nic_mac_address/images symlink -> $root_path
		//01 means ethernet
		_ , err = os.Create(fmt.Sprintf("%s/boot.cfg", serverDir))
		if err != nil {
			log.Println(err)
		}
		fileWrite, err := os.OpenFile(fmt.Sprintf("%s/boot.cfg", serverDir), os.O_RDWR, 666)
		if err != nil {
			log.Println(err)
		}
		fileRead, err := os.Open(fmt.Sprintf("%s/efi/boot/boot.cfg", cfg.RootPath))
		util.BootConfigFile(fileRead, fileWrite, cfg.BindIP, k)
		err = fileWrite.Close()
		if err != nil {
			log.Println(err)
		}
	}
}

func (cfg *Config) RemovePxeDirectory() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("clean shutdown")
		for k, _ := range cfg.Nics {
			os.RemoveAll(fmt.Sprintf("%s/01-%s", cfg.RootPath, k))
		}
		os.Exit(1)
	}()
}

