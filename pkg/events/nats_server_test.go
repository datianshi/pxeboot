package events_test

import (
	stand "github.com/nats-io/nats-streaming-server/server"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type NatServer struct {
	stop chan bool
	stand *stand.StanServer
	nc *nats.Conn
}

func NewNatServer() *NatServer {
	return &NatServer{
		stop: make(chan bool, 1),
	}
}

func (nat *NatServer) Start(){
	//go func(){
	var err error
	nat.stand, err = stand.RunServer("nat_server")
	if err != nil {
		log.Fatal(err)
	}
	nat.nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	go func(done chan bool){
		<- done
		nat.stand.Shutdown()
	}(nat.stop)
	//Give some time for server to start
	time.Sleep(200 * time.Millisecond)
}

func (nat *NatServer) Stop() {
	nat.stop <- true
	//Give some time for server to stop
	time.Sleep(200 * time.Millisecond)
}

func (nat *NatServer) Publish(subj, msg string) {
	nat.nc.Publish(subj, []byte(msg))
}

func (nat *NatServer) Subscribe(subj string, fn func(*nats.Msg)) {
	go func() {
		nat.nc.Subscribe(subj, func(msg *nats.Msg) {
			fn(msg)
		})
	}()
}