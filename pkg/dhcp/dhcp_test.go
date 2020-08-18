package dhcp_test

import (
	"github.com/datianshi/pxeboot/pkg/dhcp"
	"net"
	"testing"
	"time"
)

func TestNewServerPool(t *testing.T) {
	_, err := dhcp.NewServerPool(2, "172.16.34.5 - 172.16.34.10")
	if err!= nil {
		t.Errorf("Expect no error, but got %s", err.Error())
	}
}

func TestDHCPAssignment(t *testing.T) {
	serverPool, _ := dhcp.NewServerPool(1, "172.16.34.5 - 172.16.34.7")
	ip, err := serverPool.AssignIP()
	if err!= nil || !ip.Equal(net.ParseIP("172.16.34.5")){
		t.Errorf("Error %s, Expect IP %s got assigned", err.Error(), "172.16.34.5")
	}
	ip, err = serverPool.AssignIP()
	if err != nil || !ip.Equal(net.ParseIP("172.16.34.6")) {
		t.Errorf("Error %s, Expect IP %s got assigned", err.Error(), "172.16.34.6")
	}
	ip, err = serverPool.AssignIP()
	if err != nil || !ip.Equal(net.ParseIP("172.16.34.7")) {
		t.Errorf("Error %s, Expect IP %s got assigned", err.Error(), "172.16.34.7")
	}
	ip, err = serverPool.AssignIP()
	if err == nil {
		t.Errorf("Expect IP Pool Exausted, but got %s", ip.String())
	}
	time.Sleep(2 * time.Second)
	ip, err = serverPool.AssignIP()
	if err!= nil || !ip.Equal(net.ParseIP("172.16.34.5")){
		t.Errorf("Error %s, Expect IP %s got assigned", err.Error(), "172.16.34.5")
	}
	ip, err = serverPool.AssignIP()
	if err != nil || !ip.Equal(net.ParseIP("172.16.34.6")) {
		t.Errorf("Error %s, Expect IP %s got assigned", err.Error(), "172.16.34.6")
	}
	ip, err = serverPool.AssignIP()
	if err != nil || !ip.Equal(net.ParseIP("172.16.34.7")) {
		t.Errorf("Error %s, Expect IP %s got assigned", err.Error(), "172.16.34.7")
	}
	ip, err = serverPool.AssignIP()
	if err == nil {
		t.Errorf("Expect IP Pool Exausted, but got %s", ip.String())
	}
}