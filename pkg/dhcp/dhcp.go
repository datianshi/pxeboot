package dhcp

import (
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/http/kickstart"
	"github.com/datianshi/pxeboot/pkg/util"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"log"
	"net"
)

type DhcpServer struct{
	k kickstart.Kickstart
	c config.Config
}

func handleDHCP(cfg *config.Config, pool *ServerPool) server4.Handler {
	return func(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
		// this function will just print the received DHCPv4 message, without replying
		log.Print(m.Summary())
		client_mac := m.ClientHWAddr.String()
		_, found := cfg.Nics[util.Colon_To_Dash(client_mac)]
		if !found {
			log.Println(fmt.Sprintf("No DHCP Offering for mac address: %s", client_mac))
			return
		}
		replyIp, err := pool.AssignIP(m.ClientHWAddr)
		if err != nil {
			log.Printf("Can not assign IP error %s", err.Error())
			return
		}
		ret, _ := dhcpv4.New(
			dhcpv4.WithReply(m),
			dhcpv4.WithNetmask(net.IPv4Mask(255, 255, 255, 0)),
			dhcpv4.WithServerIP(net.ParseIP(cfg.BindIP)),
			dhcpv4.WithLeaseTime(uint32(cfg.LeaseTime)),
			dhcpv4.WithYourIP(replyIp),
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
	var port int
	if cfg.DHCPServerPort != 0 {
		port = cfg.DHCPServerPort
	} else {
		port = 67
	}
	laddr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: port,
	}
	pool, err := NewServerPool(cfg.LeaseTime, cfg.DHCPRange)
	if err != nil {
		log.Fatal(err)
	}
	server, err := server4.NewServer(cfg.DHCPInterface, &laddr, handleDHCP(cfg, pool))
	if err != nil {
		log.Fatal(err)
	}
	server.Serve()
}