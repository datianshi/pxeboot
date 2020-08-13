package dhcp

import (
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/http"
	"github.com/datianshi/pxeboot/pkg/util"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"log"
	"net"
)

type DhcpServer struct{
	k http.Kickstart
	c config.Config
}

func handleDHCP(cfg *config.Config) server4.Handler {
	return func(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
		// this function will just print the received DHCPv4 message, without replying
		log.Print(m.Summary())
		client_mac := m.ClientHWAddr.String()
		serverConfig, found := cfg.Nics[util.Colon_To_Dash(client_mac)]
		if !found {
			log.Println(fmt.Sprintf("No DHCP Offering for mac address: %s", client_mac))
			return
		}
		ret, _ := dhcpv4.New(
			dhcpv4.WithReply(m),
			dhcpv4.WithNetmask(net.IPv4Mask(255, 255, 255, 0)),
			dhcpv4.WithLeaseTime(uint32(cfg.LeaseTime)),
			dhcpv4.WithYourIP(net.ParseIP(serverConfig.Ip)),
		)
		//TFTP Server
		ret.Options[66] = []byte(cfg.BindIP)
		//BootFileName
		ret.Options[67] = []byte(cfg.BootFile)
		//TFTP Server
		log.Print(ret.Summary())
		conn.WriteTo(ret.ToBytes(), peer)
	}
}
func Start(cfg *config.Config) {
	laddr := net.UDPAddr{
		IP:   net.ParseIP(cfg.BindIP),
		Port: 67,
	}
	server, err := server4.NewServer("ens224", &laddr, handleDHCP(cfg))
	if err != nil {
		log.Fatal(err)
	}
	// This never returns. If you want to do other stuff, dump it into a
	// goroutine.
	server.Serve()
}