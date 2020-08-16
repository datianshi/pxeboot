package http

import (
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/http/api"
	"github.com/datianshi/pxeboot/pkg/http/ui"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func Start(conf *config.Config, router *mux.Router){
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:8089", conf.BindIP),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: router, // Pass our instance of gorilla/mux in.
	}
	a := api.NewAPI(conf)
	if err := ui.RegisterUITemplate(router); err != nil {
		log.Fatal(err)
	}
	router.HandleFunc("/api/conf", a.GetConfigHandler())
	router.HandleFunc("/api/conf/nic/{mac_address}", a.UpdateNicConfig()).Methods("PUT")
	router.HandleFunc("/api/conf/nic", a.CreateNicConfig()).Methods("POST")
	if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
	}
}
