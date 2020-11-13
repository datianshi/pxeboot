package dhcp

import (
	"fmt"
	"log"
	"net"

	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/nic"
	"github.com/datianshi/pxeboot/pkg/util"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

//Server Server to handle dhcp request
type Server struct {
	nicService nic.Service
	cfg        config.Config
	pool       *ServerPool
}

func (s *Server) handleDHCP() server4.Handler {
	return func(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
		log.Print(m.Summary())
		clientMAC := m.ClientHWAddr.String()
		nicConfig, err := s.nicService.FindServer(util.Colon_To_Dash(clientMAC))
		if err != nil {
			log.Println(fmt.Sprintf("No DHCP Offering for mac address: %s", clientMAC))
			return
		}
		replyIP, err := s.pool.AssignIP(m.ClientHWAddr)
		if err != nil {
			log.Printf("Can not assign IP error %s", err.Error())
			return
		}
		ret, _ := dhcpv4.New(
			dhcpv4.WithReply(m),
			dhcpv4.WithNetmask(net.IPv4Mask(255, 255, 255, 0)),
			dhcpv4.WithServerIP(net.ParseIP(nicConfig.Ip)),
			dhcpv4.WithLeaseTime(uint32(s.cfg.LeaseTime)),
			dhcpv4.WithYourIP(replyIP),
		)
		//TFTP Server
		ret.Options[66] = []byte(s.cfg.BindIP)
		//BootFileName
		ret.Options[67] = []byte(s.cfg.BootFile)
		//TFTP Server
		log.Print(ret.Summary())
		conn.WriteTo(ret.ToBytes(), peer)
	}
}

//Start Start DHCP server
func Start(cfg config.Config, nicService nic.Service) {
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
	dhcpServer := Server{
		nicService: nicService,
		cfg:        cfg,
		pool:       pool,
	}
	server, err := server4.NewServer(cfg.DHCPInterface, &laddr, dhcpServer.handleDHCP())
	if err != nil {
		log.Fatal(err)
	}
	server.Serve()
}
