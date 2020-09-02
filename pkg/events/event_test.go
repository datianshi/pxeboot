package events_test

import (
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/events"
	"github.com/nats-io/nats.go"
	"testing"
	"time"
)

var (
	natServer *NatServer = NewNatServer()
)

var test_data1 = `
{
	"aggregation_id" : "1",
	"servers": 	[
	   {
		  "mac_address":"00:50:56:82:75:e5",
		  "hostname":"server1",
		  "ip":"10.65.144.51",
		  "gateway":"10.65.144.1",
		  "netmask":"255.255.255.0"
	   },
	   {
		  "mac_address":"00:50:56:82:75:e6",
		  "hostname":"server2",
		  "ip":"10.65.144.52",
		  "gateway":"10.65.144.1",
		  "netmask":"255.255.255.0"
	   }
	]
}
`

var test_data2 = `
{
	"aggregation_id" : "2",
	"servers": 	[
	   {
		  "mac_address":"00:50:56:82:75:e5",
		  "hostname":"server1",
		  "ip":"10.65.144.51",
		  "gateway":"10.65.144.1",
		  "netmask":"255.255.255.0"
	   }
	]
}
`

var test_data3 = `
{
	"aggregation_id" : "2",
	"servers": 	[
	   {
		  "mac_address":"00:50:56:82:75:e7",
		  "hostname":"server3",
		  "ip":"10.65.144.53",
		  "gateway":"10.65.144.1",
		  "netmask":"255.255.255.0"
	   },
	   {
		  "mac_address":"00:50:56:82:75:e8",
		  "hostname":"server4",
		  "ip":"10.65.144.54",
		  "gateway":"10.65.144.1",
		  "netmask":"255.255.255.0"
	   }
	]
}
`

func TestEvent(t *testing.T) {
	natServer.Start()

	c := &config.Config{
		Nics: make(map[string]config.ServerConfig, 0),
	}

	processor, err := events.NewEventProcessor("127.0.0.1", "", "", 4222, c, events.TimeToImage(1*time.Second))
	if err != nil {
		t.Errorf("Can not start the process %s", err)
	}
	processor.Start()
	time.Sleep(100 * time.Millisecond)
	natServer.Publish("build.imaging", test_data1)
	time.Sleep(100 * time.Millisecond)
	natServer.Publish("build.imaging", test_data2)
	time.Sleep(100 * time.Millisecond)
	natServer.Publish("build.imaging", test_data3)
	time.Sleep(100 * time.Millisecond)

	var countDone int = 0
	natServer.Subscribe("build.imaging.done", func(msg *nats.Msg) {
		countDone++
		fmt.Println(string(msg.Data))
	})
	if len(c.Nics) != 4 {
		t.Error("Expect the config having 4 nic config, but got", len(c.Nics))
	}
	//The check should set imaging as done
	//fmt.Println(processor.GetImagesByAggregationID("1"))
	time.Sleep(4000 * time.Millisecond)
	images := processor.GetImagesByStatus("done")
	if len(images) != 2 {
		t.Errorf("Expect both of imagings are done, but got %d", len(images))
	}
	if len(c.Nics) != 0 {
		t.Error("Expect the two nics were removed in the config", len(c.Nics))
	}
	if countDone != 2 {
		t.Error("Expect two imaging done report on channel build.imaging.done, but got", countDone)
	}
	natServer.Stop()
}
