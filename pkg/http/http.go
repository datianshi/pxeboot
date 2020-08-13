package http

import (
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func Start(conf *config.Config, router *mux.Router){
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:80", conf.BindIP),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: router, // Pass our instance of gorilla/mux in.
	}
	if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
	}
}
