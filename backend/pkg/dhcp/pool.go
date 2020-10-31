package dhcp

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

type ServerPool struct {
	lease int
	ipRange string
	items []*item
}

func (p ServerPool) AssignIP(addr net.HardwareAddr) (net.IP, error) {
	var ip net.IP
	//checking if there is already a lease with this mac address
	for _, item := range p.items {
		if addr.String() == item.hardwareAddr.String() {
			return item.ip, nil
		}
	}
	for _, item := range p.items {
		if ip = item.take(addr); ip != nil {
			return ip, nil
		}
	}
	return nil, errors.New("ALl the IPs are assigned")
}

func NewServerPool(lease int, ipRange string) (*ServerPool, error){
	ips, err := getIps(ipRange)
	items := make([]*item, 0)
	if err != nil {
		return nil, err
	}
	ip1 := ips[0]
	ip2 := ips[1]
	var currentIp = ip1.To4()
	for bytes.Compare(currentIp, ip2.To4()) <= 0 {
		ip := make([]byte, 4)
		copy(ip, currentIp)
		i := &item{
			lease : lease,
			ip: ip,
		}
		items = append(items, i)
		currentIp[3]++
	}
	return &ServerPool{
		lease: lease,
		items: items,
		ipRange: ipRange,
	}, nil
 }

func getIps(ipRange string) ([]net.IP, error){
	regx := regexp.MustCompile(`-`)
	ret := regx.Split(ipRange, -1)
	if len(ret) != 2  {
		return nil, errors.New(fmt.Sprintf("Can not parse the ip range %s", ipRange))
	}
	var ip1, ip2 net.IP
	if ip1 = net.ParseIP(strings.TrimSpace(ret[0])); ip1 == nil {
		return nil, errors.New(fmt.Sprintf("ip  %s is not valid", ret[0]))
	}
	if ip2 = net.ParseIP(strings.TrimSpace(ret[1])); ip2 == nil {
		return nil, errors.New(fmt.Sprintf("ip  %s is not valid", ret[2]))
	}
	return []net.IP{ip1, ip2}, nil
}

type item struct {
	lease int
	taken bool
	ip net.IP
	hardwareAddr net.HardwareAddr
}

func (i *item) isTaken() bool {
	return i.taken
}

func (i *item) take(addr net.HardwareAddr) net.IP{
	if i.taken {
		return nil
	} else {
		i.taken = true
		go func() {
			timer := time.NewTimer(time.Duration(i.lease) * time.Second)
			<-timer.C
			i.hardwareAddr = nil
			i.taken = false
		}()
		i.hardwareAddr = addr
		return i.ip
	}
}
