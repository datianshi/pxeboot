package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/model"
	"github.com/datianshi/pxeboot/pkg/nic"
	"github.com/datianshi/pxeboot/pkg/util"
	"github.com/gorilla/mux"
)

//API api
type API struct {
	r             *mux.Router
	cfg           config.Config
	imageUploader *ImageUploader
	nicService    nic.Service
}

//NewAPI NewAPI
func NewAPI(c config.Config, nicService nic.Service) *API {
	return &API{
		r:   mux.NewRouter(),
		cfg: c,
		imageUploader: &ImageUploader{
			c,
		},
		nicService: nicService,
	}
}

//Start Start
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

//GetConfigHandler getConfig Handler
func (a *API) GetConfigHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		js, err := json.Marshal(a.cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//GetNics getNics
func (a *API) GetNics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := a.nicService.GetServers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		js, err := json.Marshal(items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//GetNic GetNic
func (a *API) GetNic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		macAddress := vars["mac_address"]
		item, err := a.nicService.FindServer(macAddress)
		if err != nil {
			w.WriteHeader(404) // unprocessable entity
			w.Write([]byte(fmt.Sprintf("nic %s does not exists", macAddress)))
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

//UpdateNicConfig UpdateNicConfig
func (a *API) UpdateNicConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var serverConfig model.ServerConfig
		vars := mux.Vars(r)
		macAddress := vars["mac_address"]
		_, err := a.nicService.FindServer(macAddress)
		if err != nil {
			w.WriteHeader(422) // unprocessable entity
			w.Write([]byte(fmt.Sprintf("nic %s does not exists", macAddress)))
			panic(fmt.Errorf("nic %s does not exists", macAddress))
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
			if err = a.nicService.UpdateServer(serverConfig); err != nil {
				w.WriteHeader(422)
				w.Write([]byte(fmt.Sprintf("Update failed with %v", err)))
				return
			}
			w.WriteHeader(http.StatusAccepted)
			if body, err = json.Marshal(&serverConfig); err != nil {
				w.WriteHeader(422)
				w.Write([]byte(fmt.Sprintf("Unknown return body %v", serverConfig)))
				return
			}
		}
	}
}

//DeleteNic DeleteNic
func (a *API) DeleteNic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		macAddress := vars["mac_address"]
		_, err := a.nicService.FindServer(macAddress)
		if err != nil {
			w.WriteHeader(422) // unprocessable entity
			w.Write([]byte(fmt.Sprintf("nic %s does not exists", macAddress)))
			panic(fmt.Errorf("nic %s does not exists", macAddress))
		} else {
			if err = a.nicService.DeleteServer(macAddress); err != nil {
				w.WriteHeader(422)
				return
			}
			w.WriteHeader(http.StatusAccepted)
		}
	}
}

//DeleteAllNics DeleteAllNics
func (a *API) DeleteAllNics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a.nicService.DeleteAll(); err != nil {
			w.WriteHeader(422)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

//CreateNicConfig CreateNicConfig
func (a *API) CreateNicConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var serverItem model.ServerConfig
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
				serverConfig := model.ServerConfig{0, serverItem.Ip, serverItem.Hostname, serverItem.Gateway, serverItem.Netmask, convertLowerCaseDash(serverItem.MacAddress)}
				if serverConfig, err = a.nicService.CreateServer(serverConfig); err != nil {
					w.WriteHeader(422)
					return
				}
				if body, err = json.Marshal(serverConfig); err != nil {
					w.WriteHeader(422)
					return
				}
				w.WriteHeader(http.StatusAccepted)
				w.Write(body)
			}
		}
	}
}

func convertLowerCaseDash(macAddress string) string {
	return strings.ToLower(util.Colon_To_Dash(macAddress))
}
