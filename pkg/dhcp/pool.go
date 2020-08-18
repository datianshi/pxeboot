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

func (p ServerPool) AssignIP() (net.IP, error) {
	var ip net.IP
	for _, item := range p.items {
		if ip = item.take(); ip != nil {
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
}

func (i *item) isTaken() bool {
	return i.taken
}

func (i *item) take() net.IP{
	if i.taken {
		return nil
	} else {
		i.taken = true
		go func() {
			timer := time.NewTimer(time.Duration(i.lease) * time.Second)
			<-timer.C
			i.taken = false
		}()
		return i.ip
	}
}
