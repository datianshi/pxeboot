package http

import (
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/gorilla/mux"
	"net/http"
	"text/template"
)

type Kickstart struct {
	R *mux.Router
	C *config.Config
}


func (k *Kickstart) RegisterServerEndpoint() {
	k.R.HandleFunc("/kickstart/{mac_address}/ks.cfg", k.handler())
}

type info struct {
	IP string
	Gateway string
	NetMask string
	NameServer string
	HostName string
	NTPServer string
}

func (k *Kickstart) handler() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		mac_address := vars["mac_address"]
		serverConfig, found := k.C.Nics[mac_address]
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		i := info{
			IP: serverConfig.Ip,
			HostName: serverConfig.Hostname,
			Gateway: k.C.Gateway,
			NetMask: k.C.Netmask,
			NameServer: k.C.DNS,
			NTPServer: k.C.NTPServer,
		}
		t, _:= template.New("").Parse(k.C.KickStartTemplate)
		err := t.ExecuteTemplate(w, "", i)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
}