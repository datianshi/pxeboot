package main

import (
	"github.com/datianshi/pxeboot/pkg/dhcp"
	"github.com/datianshi/pxeboot/pkg/tftp"
	"time"
)

func main() {
	go func() {
		dhcp.Start()
	}()

	go func() {
		tftp.Start()
	}()

	time.Sleep(time.Hour * 1)


}
