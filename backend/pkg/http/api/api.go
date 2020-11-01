package api

import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/util"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

type API struct {
	r             *mux.Router
	cfg           *config.Config
	imageUploader *ImageUploader
	htmlBox       *packr.Box
}

func NewAPI(c *config.Config) *API {
	return &API{
		r:   mux.NewRouter(),
		cfg: c,
		imageUploader: &ImageUploader{
			c,
		},
		htmlBox: packr.New("html", "./public"),
	}
}

func (a *API) Start() {
	var port int
	if a.cfg.HTTPPort != 0 {
		port = a.cfg.HTTPPort
	} else {
		port = 80
	}
	ipv4Adr, err := getInterfaceIpv4Addr(a.cfg.ManagementInterface)
	if err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", ipv4Adr, port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.r, // Pass our instance of gorilla/mux in.
	}
	a.r.HandleFunc("/conf", a.GetConfigHandler())
	a.r.HandleFunc("/conf/nics", a.GetNics())
	a.r.HandleFunc("/conf/nic/{mac_address}", a.GetNic()).Methods("GET")
	a.r.HandleFunc("/conf/nic/{mac_address}", a.UpdateNicConfig()).Methods("PUT")
	a.r.HandleFunc("/conf/nic/{mac_address}", a.DeleteNic()).Methods("DELETE")
	a.r.HandleFunc("/conf/deletenics", a.DeleteAllNics()).Methods("DELETE")
	a.r.HandleFunc("/conf/nic", a.CreateNicConfig()).Methods("POST")
	a.r.HandleFunc("/image", a.imageUploader.UploadHandler()).Methods("POST")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func getInterfaceIpv4Addr(interfaceName string) (addr string, err error) {
	var (
		ief      *net.Interface
		addrs    []net.Addr
		ipv4Addr net.IP
	)
	if ief, err = net.InterfaceByName(interfaceName); err != nil { // get interface
		return
	}
	if addrs, err = ief.Addrs(); err != nil { // get addresses
		return
	}
	for _, addr := range addrs { // get ipv4 address
		if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			break
		}
	}
	if ipv4Addr == nil {
		return "", fmt.Errorf("interface %s don't have an ipv4 address", interfaceName)
	}
	return ipv4Addr.String(), nil
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
		item := ServerItem{
			serverConfig.Ip,
			serverConfig.Hostname,
			mac_address,
			serverConfig.Gateway,
			serverConfig.Netmask,
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
	for k, v := range nics {
		item := ServerItem{
			v.Ip,
			v.Hostname,
			k,
			v.Gateway,
			v.Netmask,
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
