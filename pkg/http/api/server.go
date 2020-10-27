package api

import (
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Server struct {
	managementIp string
	httpPort int
	apis []*API
	r *mux.Router
}

func NewServer(config *config.Config) *Server {
	r := mux.NewRouter()
	apis := make([]*API, 0)
	for _, dhcpInterface := range config.DHCPInterfaces {
		api := NewAPI(dhcpInterface, r)
		apis = append(apis, api)
	}
	return &Server{
		managementIp: config.ManagementIp,
		httpPort: config.HTTPPort,
		r: r,
		apis: apis,
	}
}

func (s *Server) Start() {
	for _, api := range s.apis {
		api.RegisterEndpoint(api.cfg.DHCPInterface)
	}
	var port int
	if s.httpPort != 0 {
		port = s.httpPort
	} else {
		port = 80
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.managementIp, port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: s.r, // Pass our instance of gorilla/mux in.
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
