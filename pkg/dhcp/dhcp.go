package dhcp

import (
	"log"
	"net"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)
func handler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	// this function will just print the received DHCPv4 message, without replying
	log.Print(m.Summary())
	ret, _ := dhcpv4.New(
		dhcpv4.WithReply(m),
		dhcpv4.WithNetmask(net.IPv4Mask(255, 255, 255, 0)),
		dhcpv4.WithLeaseTime(uint32(5200)),
		dhcpv4.WithRouter(net.ParseIP("172.16.100.1")),
		dhcpv4.WithYourIP(net.ParseIP("172.16.100.101")),
	)
	//TFTP Server
	ret.Options[66] = []byte("172.16.100.2")
	//BootFileName
	ret.Options[67] = []byte("mboot.efi")
	//TFTP Server
	log.Print(ret.Summary())
	conn.WriteTo(ret.ToBytes(), peer)
}
func Start() {
	laddr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 67,
	}
	server, err := server4.NewServer("ens224", &laddr, handler)
	if err != nil {
		log.Fatal(err)
	}
	// This never returns. If you want to do other stuff, dump it into a
	// goroutine.
	server.Serve()
}