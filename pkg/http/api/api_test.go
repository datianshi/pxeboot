package api_test

import (
	"bytes"
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/http/api"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var data string = `
dhcp_interface: ens224
bind_ip: 172.16.100.2
gateway: 10.65.101.1
netmask: 255.255.255.0
dns: 10.192.2.10
password: VMware1!
nics:
  00-50-56-82-70-2a:
    dhcp_ip: 172.16.100.100
    ip: 10.65.101.10
    hostname: vc-01.example.org
  00-50-56-82-61-7c:
    dhcp_ip: 172.16.100.101
    ip: 10.65.101.11
    hostname: vc-02.example.org
boot_file: efi/boot/bootx64.efi
lease_time: 500
root_path: ./fixture/image
boot_config_file: efi/boot/boot.cfg
ntp_server: time.svc.pivotal.io`

func TestUpdateNicConfig(t *testing.T) {
	requestBody := `{
	"ip": "10.65.101.31", 
	"dhcp_ip": "172.16.100.102",
	"hostname": "test-host"
}`
	router := mux.NewRouter()
	var buf bytes.Buffer
	buf.WriteString(data)
	cfg, err := config.LoadConfig(&buf)
	if err != nil {
		t.Fatalf("Can not process the config %v", err)
	}
	a := api.NewAPI(cfg)
	router.HandleFunc("/api/conf/nic/{mac_address}", a.UpdateNicConfig()).Methods("PUT")
	r, err := http.NewRequest("PUT", "/api/conf/nic/00-50-56-82-70-2a", bytes.NewBufferString(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	ret, err := ioutil.ReadAll(w.Result().Body)
	fmt.Println(string(ret))

	if w.Result().StatusCode != http.StatusAccepted {
		t.Errorf("Expected the http status code %d, but got %d", http.StatusAccepted, w.Result().StatusCode)
	}
	server, found := cfg.Nics["00-50-56-82-70-2a"]
	if !found {
		t.Errorf("New nic config is not created")
	}
	if server.Ip != "10.65.101.31" {
		t.Errorf("Expected server update to %s, but got %s", "10.65.101.31", server.Ip)
	}

}

func TestCreateNicConfig(t *testing.T) {
	requestBody := `
{
	"mac_address": 
	"00:50:A6:83:70:98", 
	"ip": "10.65.101.31" , 
	"dhcp_ip": "172.16.100.102", 
	"hostname": "test-host" 
}
`
	router := mux.NewRouter()
	var buf bytes.Buffer
	buf.WriteString(data)
	cfg, err := config.LoadConfig(&buf)
	if err != nil {
		t.Fatalf("Can not process the config %v", err)
	}
	a := api.NewAPI(cfg)
	router.HandleFunc("/api/conf/nic", a.CreateNicConfig()).Methods("POST")
	r, err := http.NewRequest("POST", "/api/conf/nic", bytes.NewBufferString(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusAccepted {
		t.Errorf("Expected the http status code %d, but got %d", http.StatusAccepted, w.Result().StatusCode)
	}
	_, found := cfg.Nics["00-50-a6-83-70-98"]
	if !found {
		t.Errorf("New nic config is not created")
	}
}