package kickstart

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/nic"
	"github.com/gorilla/mux"
)

//Kickstart Kick Start Struct
type Kickstart struct {
	r *mux.Router
	C config.Config
	N nic.Service
}

//NewKickStart Initialize Kick Start
func NewKickStart(c config.Config, N nic.Service) *Kickstart {
	return &Kickstart{
		r: mux.NewRouter(),
		C: c,
		N: N,
	}
}

//Start start kick start server
func (k *Kickstart) Start() {
	var port int
	if k.C.HTTPPort != 0 {
		port = k.C.HTTPPort
	} else {
		port = 80
	}
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", k.C.BindIP, port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      k.r, // Pass our instance of gorilla/mux in.
	}
	k.r.HandleFunc("/kickstart/{mac_address}/ks.cfg", k.Handler())
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

type info struct {
	IP         string
	Gateway    string
	NetMask    string
	NameServer string
	HostName   string
	NTPServer  string
	Password   string
}

//Handler Kick Start Handler
func (k *Kickstart) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		macAddress := vars["mac_address"]
		serverConfig, err := k.N.FindServer(macAddress)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		i := info{
			IP:         serverConfig.Ip,
			HostName:   serverConfig.Hostname,
			Gateway:    serverConfig.Gateway,
			NetMask:    serverConfig.Netmask,
			NameServer: k.C.DNS,
			NTPServer:  k.C.NTPServer,
			Password:   k.C.Password,
		}
		t, _ := template.New("").Parse(k.C.KickStartTemplate)
		err = t.ExecuteTemplate(w, "", i)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
}
