package main

import (
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/dhcp"
	"github.com/datianshi/pxeboot/pkg/http"
	"github.com/datianshi/pxeboot/pkg/tftp"
	"github.com/datianshi/pxeboot/pkg/util"
	"github.com/gorilla/mux"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("can not load the config", err)
	}

	//Create directories for each nic
	for k, _ := range cfg.Nics {
		//$root_path/nic_mac_address
		serverDir := fmt.Sprintf("%s/%s", cfg.RootPath, k)
		err = os.Mkdir(serverDir, 0755)
		if err != nil {
			log.Println(err)
		}
		//$root_path/nic_mac_address/images symlink -> $root_path
		err = os.Symlink(cfg.RootPath, fmt.Sprintf("%s/%s", serverDir, util.SYMLINK_PER_SERVER_DIR))
		_ , err := os.Create(fmt.Sprintf("%s/boot.cfg", serverDir))
		if err != nil {
			log.Println(err)
		}
		fileWrite, err := os.OpenFile(fmt.Sprintf("%s/boot.cfg", serverDir), os.O_RDWR, 666)
		if err != nil {
			log.Println(err)
		}
		fileRead, err := os.Open(fmt.Sprintf("%s/efi/boot/boot.cfg", cfg.RootPath))
		kickstartUrl := fmt.Sprintf("http://%s/kickstart/%s/ks.cfg", cfg.BindIP, k)
		util.BootConfigFile(fileRead, fileWrite, kickstartUrl)
		err = fileWrite.Close()
		if err != nil {
			log.Println(err)
		}
	}

	//Load Http Endpoint
	router := mux.NewRouter()
	k := http.Kickstart{
		R: router,
		C: cfg,
	}
	k.RegisterServerEndpoint()
	go func(){
		http.Start(cfg, router)
	}()

	//Start TFTP
	go func() {
		tftp.Start(cfg)
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("clean shutdown")
		for k, _ := range cfg.Nics {
			os.RemoveAll(fmt.Sprintf("%s/%s", cfg.RootPath, k))
		}
		os.Exit(1)
	}()

	//Start dhcp server and block
	dhcp.Start(cfg)


}
