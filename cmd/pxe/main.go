package main

import (
	"github.com/datianshi/pxeboot/pkg/dhcp"
	"github.com/datianshi/pxeboot/pkg/kickstart"
	"github.com/datianshi/pxeboot/pkg/tftp"
	"github.com/gorilla/mux"
	"time"
)

func main() {
	go func() {
		tftp.Start()
	}()

	k := kickstart.Kickstart{
		R: mux.NewRouter(),
	}

	k.Start()

	go func() {
		dhcp.Start(k)
	}()

	time.Sleep(time.Hour * 1)


}
