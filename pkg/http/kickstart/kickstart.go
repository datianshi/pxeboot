package kickstart

import (
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
	"time"
)

type Kickstart struct {
	r *mux.Router
	C *config.Config
}

func NewKickStart(c *config.Config) *Kickstart {
	return &Kickstart{
		r : mux.NewRouter(),
		C : c,
	}
}

func (k *Kickstart) Start() {
		var port int
		if k.C.HTTPPort != 0 {
			port = k.C.HTTPPort
		} else {
			port = 80
		}
		srv := &http.Server{
			Addr:         fmt.Sprintf("%s:%d", k.C.BindIP, port),
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler: k.r, // Pass our instance of gorilla/mux in.
		}
		k.r.HandleFunc("/kickstart/{mac_address}/ks.cfg", k.Handler())
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
}

type info struct {
	IP string
	Gateway string
	NetMask string
	NameServer string
	HostName string
	NTPServer string
	Password string
}

func (k *Kickstart) Handler() http.HandlerFunc{
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
			Password: k.C.Password,
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