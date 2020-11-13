package api_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/http/api"
	"github.com/datianshi/pxeboot/pkg/model"
	"github.com/datianshi/pxeboot/pkg/nic"
	"github.com/datianshi/pxeboot/pkg/nic/nicfakes"
	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/util/json"
)

var data string = `
dhcp_interface: ens224
bind_ip: 172.16.100.2
dns: 10.192.2.10
password: VMware1!
nics:
  00-50-56-82-70-2a:
    dhcp_ip: 172.16.100.100
    ip: 10.65.101.10
    hostname: vc-01.example.org
    gateway: 10.65.101.1
    netmask: 255.255.255.0
  00-50-56-82-61-7c:
    dhcp_ip: 172.16.100.101
    ip: 10.65.101.11
    hostname: vc-02.example.org
    gateway: 10.65.101.1
    netmask: 255.255.255.0
boot_file: efi/boot/bootx64.efi
lease_time: 500
root_path: ./fixture/image
boot_config_file: efi/boot/boot.cfg
ntp_server: time.svc.pivotal.io`

func TestUpdateNicConfig(t *testing.T) {
	requestBody := `{
	"ip": "10.65.101.31", 
	"dhcp_ip": "172.16.100.102",
	"hostname": "test-host",
    "gateway": "10.65.101.1",
	"netmask": "255.255.255.0"
}`

	nicService := &nicfakes.FakeService{}
	nicService.UpdateServerStub = func(n model.ServerConfig) error {
		return nil
	}

	router, a := setupAPI(t, nicService)
	router.HandleFunc("/api/conf/nic/{mac_address}", a.UpdateNicConfig()).Methods("PUT")
	r, _ := http.NewRequest("PUT", "/api/conf/nic/00-50-56-82-70-2a", bytes.NewBufferString(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if nicService.UpdateServerCallCount() != 1 {
		t.Errorf("Expect update server to be called once")
	}

	if w.Result().StatusCode != http.StatusAccepted {
		t.Errorf("Expected the http status code %d, but got %d", http.StatusAccepted, w.Result().StatusCode)
	}

}

func TestDeleteNicConfig(t *testing.T) {
	nicService := &nicfakes.FakeService{}
	nicService.DeleteServerStub = func(string) error {
		return nil
	}

	router, a := setupAPI(t, nicService)
	router.HandleFunc("/api/conf/nic/{mac_address}", a.DeleteNic()).Methods("DELETE")
	r, _ := http.NewRequest("DELETE", "/api/conf/nic/00-50-56-82-70-2a", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusAccepted {
		t.Errorf("Expected the http status code %d, but got %d", http.StatusAccepted, w.Result().StatusCode)
	}

	if nicService.DeleteServerCallCount() != 1 {
		t.Errorf("Expecte Delete Server called once")
	}
}

func TestDeleteAllNics(t *testing.T) {
	nicService := &nicfakes.FakeService{}
	nicService.DeleteAllStub = func() error {
		return nil
	}

	router, a := setupAPI(t, nicService)
	router.HandleFunc("/api/conf/deletenics", a.DeleteAllNics()).Methods("DELETE")
	r, _ := http.NewRequest("DELETE", "/api/conf/deletenics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusAccepted {
		t.Errorf("Expected the http status code %d, but got %d", http.StatusAccepted, w.Result().StatusCode)
	}

	if nicService.DeleteAllCallCount() != 1 {
		t.Errorf("Expecte Delete Server called once")
	}
}

func TestCreateNicWithRoundMacAddr(t *testing.T) {
	requestBodyWithWrongMac := `
{
	"mac_address": "00:50:A6:83:7:98",
	"ip": "10.65.101.31",
	"dhcp_ip": "172.16.100.102",
	"hostname": "test-host"
}
`
	nicService := &nicfakes.FakeService{}
	nicService.CreateServerStub = func(server model.ServerConfig) (model.ServerConfig, error) {
		return server, nil
	}
	router, a := setupAPI(t, nicService)
	router.HandleFunc("/api/conf/nic", a.CreateNicConfig()).Methods("POST")
	r, _ := http.NewRequest("POST", "/api/conf/nic", bytes.NewBufferString(requestBodyWithWrongMac))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Expected the http status code %d, but got %d", http.StatusBadRequest, w.Result().StatusCode)
	}
}

func TestCreateNicConfig(t *testing.T) {
	requestBody := `
{
	"mac_address": "00:50:A6:83:70:98",
	"ip": "10.65.101.31" ,
	"dhcp_ip": "172.16.100.102",
	"hostname": "test-host",
    "gateway": "10.65.101.1",
	"netmask": "255.255.255.0"
}
`
	nicService := &nicfakes.FakeService{}
	nicService.CreateServerStub = func(server model.ServerConfig) (model.ServerConfig, error) {
		server.ID = 100
		return server, nil
	}
	router, a := setupAPI(t, nicService)
	router.HandleFunc("/api/conf/nic", a.CreateNicConfig()).Methods("POST")
	r, _ := http.NewRequest("POST", "/api/conf/nic", bytes.NewBufferString(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusAccepted {
		t.Errorf("Expect the http status code %d, but got %d", http.StatusAccepted, w.Result().StatusCode)
	}

	if nicService.CreateServerCallCount() != 1 {
		t.Errorf("Expect Create Server called once")
	}

	ret, _ := ioutil.ReadAll(w.Result().Body)
	var server model.ServerConfig
	if err := json.Unmarshal(ret, &server); err != nil {
		t.Errorf("Invalid Response")
	}

	if server.ID != 100 {
		t.Errorf("Expect server id is generated by the server side, but got %d", server.ID)
	}

}

func setupAPI(t *testing.T, nicService nic.Service) (*mux.Router, *api.API) {
	router := mux.NewRouter()
	var buf bytes.Buffer
	buf.WriteString(data)
	cfg, err := config.LoadConfig(&buf)
	if err != nil {
		t.Fatalf("Can not process the config %v", err)
	}
	a := api.NewAPI(cfg, nicService)
	return router, a
}
