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
	mac1, _ := net.ParseMAC("00:50:56:82:75:e5")
	mac2, _ := net.ParseMAC("00:50:56:82:75:e6")
	mac3, _ := net.ParseMAC("00:50:56:82:75:e7")
	mac4, _ := net.ParseMAC("00:50:56:82:75:e8")
	ip, err := serverPool.AssignIP(mac1)
	if err!= nil || !ip.Equal(net.ParseIP("172.16.34.5")){
		t.Errorf("Error %s, Expect IP %s got assigned", err.Error(), "172.16.34.5")
	}
	ip, err = serverPool.AssignIP(mac2)
	if err != nil || !ip.Equal(net.ParseIP("172.16.34.6")) {
		t.Errorf("Error %s, Expect IP %s got assigned", err.Error(), "172.16.34.6")
	}
	ip, err = serverPool.AssignIP(mac1)
	if !ip.Equal(net.ParseIP("172.16.34.5")){
		t.Errorf("Expect IP %s got reassigned", "172.16.34.5")
	}
	ip, err = serverPool.AssignIP(mac3)
	if err != nil || !ip.Equal(net.ParseIP("172.16.34.7")) {
		t.Errorf("Error %s, Expect IP %s got assigned", err.Error(), "172.16.34.7")
	}
	ip, err = serverPool.AssignIP(mac4)
	if err == nil {
		t.Errorf("Expect IP Pool Exausted, but got %s", ip.String())
	}
	time.Sleep(2 * time.Second)
	ip, err = serverPool.AssignIP(mac2)
	if err != nil || !ip.Equal(net.ParseIP("172.16.34.5")) {
		t.Errorf("Error %s, Expect IP %s got assigned", err.Error(), "172.16.34.5")
	}
	ip, err = serverPool.AssignIP(mac3)
	if err != nil || !ip.Equal(net.ParseIP("172.16.34.6")) {
		t.Errorf("Error %s, Expect IP %s got assigned", err.Error(), "172.16.34.6")
	}
	ip, err = serverPool.AssignIP(mac1)
	if err!= nil || !ip.Equal(net.ParseIP("172.16.34.7")){
		t.Errorf("Error %s, Expect IP %s got assigned", err.Error(), "172.16.34.7")
	}
}