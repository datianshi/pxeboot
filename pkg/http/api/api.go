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
	imageUploader *ImageUploader
}

func NewAPI(c *config.Config) *API {
	return &API{
		r: mux.NewRouter(),
		cfg: c,
		imageUploader: &ImageUploader{
			c,
		},
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
	a.r.HandleFunc("/api/conf/nics", a.GetNics())
	a.r.HandleFunc("/api/conf/nic/{mac_address}", a.GetNic()).Methods("GET")
	a.r.HandleFunc("/api/conf/nic/{mac_address}", a.UpdateNicConfig()).Methods("PUT")
	a.r.HandleFunc("/api/conf/nic/{mac_address}", a.DeleteNic()).Methods("DELETE")
	a.r.HandleFunc("/api/conf/deletenics", a.DeleteAllNics()).Methods("DELETE")
	a.r.HandleFunc("/api/conf/nic", a.CreateNicConfig()).Methods("POST")
	a.r.HandleFunc("/api/image", a.imageUploader.UploadHandler()).Methods("POST")
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

func (api *API) GetNics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		js, err := json.Marshal(convertToServerItems(api.cfg.Nics))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (api *API) GetNic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var serverConfig config.ServerConfig
		vars := mux.Vars(r)
		mac_address := vars["mac_address"]
		serverConfig , found := api.cfg.Nics[mac_address]
		if !found {
			w.WriteHeader(404) // unprocessable entity
			w.Write([]byte(fmt.Sprintf("nic %s does not exists", mac_address)))
		}
		item := ServerItem{
			serverConfig.Ip,
			serverConfig.DhcpIp,
			serverConfig.Hostname,
			mac_address,
		}
		js, err := json.Marshal(item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func convertToServerItems(nics map[string]config.ServerConfig) []ServerItem {
	var items []ServerItem
	for k,v := range nics {
		item := ServerItem{
			v.Ip,
			v.DhcpIp,
			v.Hostname,
			k,
		}
		items = append(items, item)
	}
	return items
}


func (api *API) UpdateNicConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("update")
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

func (api *API) DeleteNic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		mac_address := vars["mac_address"]
		_ , found := api.cfg.Nics[mac_address]
		if !found {
			w.WriteHeader(422) // unprocessable entity
			w.Write([]byte(fmt.Sprintf("nic %s does not exists", mac_address)))
			panic(errors.New(fmt.Sprintf("nic %s does not exists", mac_address)))
		} else {
			delete(api.cfg.Nics, mac_address)
			w.WriteHeader(http.StatusAccepted)
		}
	}
}

func (api *API) DeleteAllNics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api.cfg.Nics = map[string]config.ServerConfig{}
		w.WriteHeader(http.StatusAccepted)
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
			if err:= serverItem.Validate(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
			} else {
				fmt.Println("still come here")
				serverConfig := config.ServerConfig{serverItem.Ip, serverItem.DhcpIp, serverItem.Hostname}
				api.cfg.Nics[convertLowerCaseDash(serverItem.MacAddress)] = serverConfig
				w.WriteHeader(http.StatusAccepted)
			}
		}
	}
}

func convertLowerCaseDash(mac_address string) string{
	return strings.ToLower(util.Colon_To_Dash(mac_address))
}
