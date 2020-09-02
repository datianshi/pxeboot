package events

import (
	"encoding/json"
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/model"
	"github.com/datianshi/pxeboot/pkg/util"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type EventProcessor struct {
	client * nats.Conn
	config *config.Config
	repository Repository
	checker CheckImagingDone
}

type CheckImagingDone func(repository Repository) []Imaging


func TimeToImage(duration time.Duration) CheckImagingDone {
	return func(repository Repository) []Imaging {
		needReconcil := make([]Imaging, 0)
		IDs, err := repository.GetAll()
		if err != nil {
			log.Printf("Can not get IDs in repository %s", err)
			return needReconcil
		}
		for _, id := range IDs {
			var err error
			var image Imaging
			if image, err = repository.Get(id); err != nil {
				log.Printf("Can not get ID %s in repository ", id)
				continue
			}
			if image.Status != "processing" {
				continue
			}
			if time.Now().After(image.StartTime.Add(duration)) {
				image.Status = "done"
				repository.Save(image)
				needReconcil = append(needReconcil, image)
			}
		}
		return needReconcil
	}
}

//func PingIPChecker(repository Repository) []Imaging{
//	needReconcil := make([]Imaging, 0)
//	IDs, err := repository.GetAll()
//	if err != nil {
//		log.Printf("Can not get IDs in repository %s", err)
//		return needReconcil
//	}
//	for _, id := range IDs {
//		var err error
//		var image Imaging
//		if image, err = repository.Get(id); err != nil {
//			log.Printf("Can not get ID %s in repository ", id)
//		}
//		//300 seconds time out
//		if time.Now().After(image.StartTime.Add(300 * time.Second)){
//			image.Status = "failure"
//			repository.Save(image)
//			needReconcil = append(needReconcil, image)
//		}else if image.Status == "processing" {
//			var done bool = true
//			for _, item := range image.Servers {
//				pinger, _ := ping.NewPinger(item.Ip)
//				pinger.Timeout = 2 * time.Second
//				pinger.Count = 1
//				pinger.Run()
//				if pinger.PacketsRecv != 1 {
//					done = false
//					break
//				}
//			}
//			if done {
//				image.Status = "done"
//				repository.Save(image)
//				needReconcil = append(needReconcil, image)
//			}
//		}
//	}
//	return needReconcil
//}

type EventModel struct {
	AggregationID string `json:"aggregation_id"`
	Servers []model.ServerItem `json:"servers"`
}

type Imaging struct {
	AggregationID string
	ID string
	Servers []model.ServerItem
	Status string
	StartTime time.Time
}

type Repository interface {
	Get(string) (Imaging, error)
	GetAll() ([]string, error)
	Save(Imaging) error
	AnyImageRunning(items []model.ServerItem) bool
}

type InMemoryRepository struct {
	data map[string]Imaging
}

func NewInMemoryRepository() *InMemoryRepository{
	return &InMemoryRepository{
		data : make(map[string]Imaging, 0),
	}
}

func(r *InMemoryRepository) Get(id string) (Imaging, error){
	return r.data[id], nil
}

func(r *InMemoryRepository) GetAll() ([]string, error){
	keys := make([]string, 0, len(r.data))
	for k := range r.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func(r *InMemoryRepository) Save(imaging Imaging) error {
	if imaging.ID == "" {
		uid, err := util.UUID()
		if err != nil {
			return err
		}
		imaging.ID = uid
	}
	if imaging.Status == "" {
		imaging.Status = "processing"
	}
	uid, err := util.UUID()
	if err != nil {
		return err
	}
	r.data[uid] = imaging
	return nil
}

func NewEventProcessor(host, username, password string, port int, config *config.Config, checker CheckImagingDone)  (*EventProcessor, error) {
	var connectString string
	if username == "" {
		connectString = fmt.Sprintf("nats://%s:%d", host, port)
	}else {
		connectString = fmt.Sprintf("nats://%s:%s@%s:%d", username, password, host, port)
	}
	client, err := nats.Connect(connectString)
	if err != nil {
		return nil, err
	}
	return &EventProcessor{
		client: client,
		config: config,
		repository: NewInMemoryRepository(),
		checker: checker,
	}, nil
}

func (ep *EventProcessor) Start() {
	ep.client.Subscribe("build.imaging", func(msg *nats.Msg) {
		go func(){
			var event EventModel
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				log.Printf("Error %s. Can not imaging request %s", err.Error(), string(msg.Data))
				return
			}
			servers := event.Servers
			if ep.repository.AnyImageRunning(servers){
				log.Printf("Servers are already under imaging %s",  string(msg.Data))
				return
			}
			for index, _ := range servers {
				c := config.ServerConfig{servers[index].Ip, servers[index].Hostname, servers[index].Gateway, servers[index].Netmask}
				ep.config.Nics[util.Colon_To_Dash(servers[index].MacAddress)] = c
			}
			ep.repository.Save(Imaging{
				AggregationID: event.AggregationID,
				Servers: servers,
				StartTime: time.Now(),
			})
		}()
	})
	ep.check()
}

func (ep *EventProcessor) GetImagesByStatus(status string) []Imaging{
	images := make([]Imaging, 0)
	IDs, _ := ep.repository.GetAll()
	for _, id := range IDs {
		item, _ := ep.repository.Get(id)
		if item.Status == status{
			images = append(images, item)
		}
	}
	return images

}

func (ep *EventProcessor) GetImagesByAggregationID(aggregationId string) []Imaging{
	images := make([]Imaging, 0)
	IDs, _ := ep.repository.GetAll()
	for _, id := range IDs {
		item, _ := ep.repository.Get(id)
		if item.AggregationID == aggregationId{
			images = append(images, item)
		}
	}
	return images
}
//Reconcile every 3 seconds. Remove the nics from imaging server
//publish done message
func (ep *EventProcessor) check() {
	go func(){
		for {
			needReconcile := ep.checker(ep.repository)
			for _, image := range needReconcile {
				for _, server := range image.Servers {
					delete(ep.config.Nics, util.Colon_To_Dash(server.MacAddress))
				}
				msg := fmt.Sprintf(`{"aggregation_id" : "%s", "status" : "%s"}`, image.AggregationID, image.Status)
				ep.client.Publish("build.imaging.done", []byte(msg))
			}
			time.Sleep(3 * time.Second)
		}
	}()
}


func (repository *InMemoryRepository) AnyImageRunning(items []model.ServerItem) bool {
	for _, image := range repository.data {
		if image.Status == "processing" {
			for _ , item := range items {
				for _, server:= range image.Servers {
					if item.MacAddress == server.MacAddress {
						return true
					}
				}
			}
		}
	}
	return false
}