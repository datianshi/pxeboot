package api_test

import (
	"bytes"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/http/api"
	test_util "github.com/datianshi/pxeboot/pkg/testutil"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestUpdateNicConfig(t *testing.T) {
	requestBody := `{
	"ip": "10.65.101.31",  
	"hostname": "test-host",
	"netmask": "255.255.255.0",
	"gateway": "10.65.101.1"
}`
	router, cfg, a := setupAPI(t)
	router.HandleFunc("/api/conf/nic/{mac_address}", a.UpdateNicConfig()).Methods("PUT")
	r, _ := http.NewRequest("PUT", "/api/conf/nic/00-50-56-82-70-2a", bytes.NewBufferString(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	ioutil.ReadAll(w.Result().Body)

	if w.Result().StatusCode != http.StatusAccepted {
		t.Errorf("Expected the http status code %d, but got %d", http.StatusAccepted, w.Result().StatusCode)
	}
	server, found := cfg.Nics["00-50-56-82-70-2a"]
	if !found {
		t.Errorf("New nic config is not updated")
	}
	if server.Ip != "10.65.101.31" {
		t.Errorf("Expected server update to %s, but got %s", "10.65.101.31", server.Ip)
	}

}

func TestDeleteNicConfig(t *testing.T) {
	router, cfg, a := setupAPI(t)
	router.HandleFunc("/api/conf/nic/{mac_address}", a.DeleteNic()).Methods("DELETE")
	r, _ := http.NewRequest("DELETE", "/api/conf/nic/00-50-56-82-70-2a", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusAccepted {
		t.Errorf("Expected the http status code %d, but got %d", http.StatusAccepted, w.Result().StatusCode)
	}
	_, found := cfg.Nics["00-50-56-82-70-2a"]
	if found {
		t.Errorf("Expect %s nic config is deleted", "00-50-56-82-70-2a")
	}
}

func TestDeleteAllNics(t *testing.T) {
	router, cfg, a := setupAPI(t)
	router.HandleFunc("/api/conf/deletenics", a.DeleteAllNics()).Methods("DELETE")
	r, _ := http.NewRequest("DELETE", "/api/conf/deletenics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusAccepted {
		t.Errorf("Expected the http status code %d, but got %d", http.StatusAccepted, w.Result().StatusCode)
	}
	if len(cfg.Nics) != 0 {
		t.Errorf("Expect to have 0 Nics in config")
	}
}

func TestCreateNicConfig(t *testing.T) {
	requestBody := `{
	"mac_address":"00:50:A6:83:76:98",
	"ip": "10.65.101.31",
	"hostname": "test-host",
	"netmask": "255.255.255.0",
	"gateway": "10.65.101.1"
}`
	router, cfg, a := setupAPI(t)
	router.HandleFunc("/api/conf/nic", a.CreateNicConfig()).Methods("POST")
	r, _ := http.NewRequest("POST", "/api/conf/nic", bytes.NewBufferString(requestBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Result().StatusCode != http.StatusAccepted {
		t.Errorf("Expected the http status code %d, but got %d", http.StatusAccepted, w.Result().StatusCode)
	}
	_, found := cfg.Nics["00-50-a6-83-76-98"]
	if !found {
		t.Errorf("New nic config is not created")
	}

	requestBodyWithWrongMac := `
{
	"mac_address":"00:50:A6:83:7:98", 
	"ip": "10.65.101.31" ,
	"hostname": "test-host" 
}
`
	r, _ = http.NewRequest("POST", "/api/conf/nic", bytes.NewBufferString(requestBodyWithWrongMac))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Expected the http status code %d, but got %d", http.StatusBadRequest, w.Result().StatusCode)
	}
}


func setupAPI(t *testing.T) (*mux.Router, *config.DHCPInterface, *api.API) {
	router := mux.NewRouter()
	cfg, err := test_util.GetConfig("../../testutil/fixture/cfg.yaml")
	if err != nil {
		t.Fatalf("Can not process the config %v", err)
	}
	a := api.NewAPI(cfg.DHCPInterfaces[0], router)
	return router, cfg.DHCPInterfaces[0], a
}
