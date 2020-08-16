package api

import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/util"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)


type API struct {
	r *mux.Router
	cfg *config.Config
}

type ServerItem struct {
	Ip string `yaml:"ip" json:"ip"`
	DhcpIp string `yaml:"dhcp_ip" json:"dhcp_ip"`
	Hostname string `yaml:"hostname" json:"hostname"`
	MacAddress string `json:"mac_address"`
}

func NewAPI(c *config.Config) *API {
	return &API{
		r: mux.NewRouter(),
		cfg: c,
	}
}

func (a *API) Start() {
	var port int
	if a.cfg.HTTPPort != 0 {
		port = a.cfg.HTTPPort
	} else {
		port = 80
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", a.cfg.ManagementIp, port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: a.r, // Pass our instance of gorilla/mux in.
	}
	a.r.HandleFunc("/api/conf", a.GetConfigHandler())
	a.r.HandleFunc("/api/conf/nic/{mac_address}", a.UpdateNicConfig()).Methods("PUT")
	a.r.HandleFunc("/api/conf/nic", a.CreateNicConfig()).Methods("POST")
	if err := RegisterUITemplate(a.r); err != nil {
		log.Fatal(err)
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (api *API) GetConfigHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		js, err := json.Marshal(api.cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}


func (api *API) UpdateNicConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var serverConfig config.ServerConfig
		vars := mux.Vars(r)
		mac_address := vars["mac_address"]
		_ , found := api.cfg.Nics[mac_address]
		if !found {
			w.WriteHeader(422) // unprocessable entity
			w.Write([]byte(fmt.Sprintf("nic %s does not exists", mac_address)))
			panic(errors.New(fmt.Sprintf("nic %s does not exists", mac_address)))
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(body, &serverConfig); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		} else {
			api.cfg.Nics[mac_address] = serverConfig
			w.WriteHeader(http.StatusAccepted)
		}
	}
}

func (api *API) CreateNicConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var serverItem ServerItem
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(body, &serverItem); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		} else {
			serverConfig := config.ServerConfig{serverItem.Ip, serverItem.DhcpIp, serverItem.Hostname}
			api.cfg.Nics[convertLowerCaseDash(serverItem.MacAddress)] = serverConfig
			w.WriteHeader(http.StatusAccepted)
		}
	}
}

func convertLowerCaseDash(mac_address string) string{
	return strings.ToLower(util.Colon_To_Dash(mac_address))
}
