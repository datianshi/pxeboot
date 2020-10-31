package main

import (
	"flag"
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/dhcp"
	"github.com/datianshi/pxeboot/pkg/http/api"
	"github.com/datianshi/pxeboot/pkg/http/kickstart"
	"github.com/datianshi/pxeboot/pkg/tftp"
	"log"
	"os"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Config File for pxe boot")
	flag.Parse()
	if "" == configPath {
		log.Fatalln("[Usage] - pxe -config pxe.yaml")
	}

	//Load Config
	configfile, err := os.Open(configPath)
	defer configfile.Close()
	if err != nil {
		log.Fatalln(fmt.Sprintf("can not open the file %s", configPath))
	}
	cfg, err := config.LoadConfig(configfile)
	if err != nil {
		log.Fatalln("can not load the config", err)
	}

	//start kickstart http server
	k := kickstart.NewKickStart(cfg)
	go func(){
		k.Start()
	}()
	//Start management api server
	if cfg.ManagementIp != "" {
		a := api.NewAPI(cfg)
		go func(){
			a.Start()
		}()
	}

	//Start TFTP
	go func() {
		tftp.Start(cfg)
	}()



	//Start dhcp server and block
	dhcp.Start(cfg)
}
