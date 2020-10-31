package api

import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/util"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

type API struct {
	r             *mux.Router
	cfg           *config.DHCPInterface
	imageUploader *ImageUploader
	htmlBox       *packr.Box
}

func NewAPI(c *config.DHCPInterface, r *mux.Router) *API {
	return &API{
		r:   r,
		cfg: c,
		imageUploader: &ImageUploader{
			c,
		},
		htmlBox: packr.New("html", "./public"),
	}
}

func (a *API) RegisterEndpoint(interfaceName string) {
	a.r.HandleFunc(fmt.Sprintf("/api/%s/conf", interfaceName), a.GetConfigHandler())
	a.r.HandleFunc(fmt.Sprintf("/api/%s/conf/nics", interfaceName), a.GetNics())
	a.r.HandleFunc(fmt.Sprintf("/api/%s/conf/nic/{mac_address}", interfaceName), a.GetNic()).Methods("GET")
	a.r.HandleFunc(fmt.Sprintf("/api/%s/conf/nic/{mac_address}", interfaceName), a.UpdateNicConfig()).Methods("PUT")
	a.r.HandleFunc(fmt.Sprintf("/api/%s/conf/nic/{mac_address}", interfaceName), a.DeleteNic()).Methods("DELETE")
	a.r.HandleFunc(fmt.Sprintf("/api/%s/conf/deletenics", interfaceName), a.DeleteAllNics()).Methods("DELETE")
	a.r.HandleFunc(fmt.Sprintf("/api/%s/conf/nic", interfaceName), a.CreateNicConfig()).Methods("POST")
	a.r.HandleFunc(fmt.Sprintf("/api/%s/image", interfaceName), a.imageUploader.UploadHandler()).Methods("POST")
	if err := a.RegisterUITemplate(a.r); err != nil {
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
		serverConfig, found := api.cfg.Nics[mac_address]
		if !found {
			w.WriteHeader(404) // unprocessable entity
			w.Write([]byte(fmt.Sprintf("nic %s does not exists", mac_address)))
		}
		item := ServerItem{serverConfig.Ip, serverConfig.Hostname, mac_address, serverConfig.Gateway, serverConfig.Netmask}
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
	for k, v := range nics {
		item := ServerItem{v.Ip, v.Hostname, k, v.Gateway, v.Netmask}
		items = append(items, item)
	}
	return items
}

func (api *API) UpdateNicConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var serverConfig config.ServerConfig
		vars := mux.Vars(r)
		mac_address := vars["mac_address"]
		_, found := api.cfg.Nics[mac_address]
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
		_, found := api.cfg.Nics[mac_address]
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
			if err := serverItem.Validate(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
			} else {
				serverConfig := config.ServerConfig{serverItem.Ip, serverItem.Hostname, serverItem.Gateway, serverItem.Netmask}
				api.cfg.Nics[convertLowerCaseDash(serverItem.MacAddress)] = serverConfig
				w.WriteHeader(http.StatusAccepted)
			}
		}
	}
}

func convertLowerCaseDash(mac_address string) string {
	return strings.ToLower(util.Colon_To_Dash(mac_address))
}
