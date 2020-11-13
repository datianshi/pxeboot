package nic

import (
	"fmt"

	"github.com/datianshi/pxeboot/pkg/db"
	"github.com/datianshi/pxeboot/pkg/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Service

//Service service to control servers
type Service interface {
	GetServers() ([]model.ServerConfig, error)
	FindServer(string) (model.ServerConfig, error)
	DeleteServer(string) error
	UpdateServer(model.ServerConfig) error
	CreateServer(model.ServerConfig) (model.ServerConfig, error)
	DeleteAll() error
}

//NewService create a new service
func NewService(database *db.Database) Service {
	return &DefaultService{
		dao: database,
	}
}

//DefaultService service backed by postgres database
type DefaultService struct {
	dao *db.Database
}

//GetServers get servers
func (service *DefaultService) GetServers() ([]model.ServerConfig, error) {
	return service.dao.GetServers()
}

//FindServer based on the mac address
func (service *DefaultService) FindServer(macAddress string) (model.ServerConfig, error) {
	return service.dao.FindServer(macAddress)
}

//DeleteServer based on the mac address
func (service *DefaultService) DeleteServer(macAddress string) error {
	if _, err := service.dao.FindServer(macAddress); err != nil {
		return fmt.Errorf("Can not find the server with mac address %s", macAddress)
	}
	return service.dao.DeleteServer(macAddress)
}

//UpdateServer Patch the config
func (service *DefaultService) UpdateServer(s model.ServerConfig) error {
	if _, err := service.dao.FindServer(s.MacAddress); err != nil {
		return fmt.Errorf("Can not find the server with mac address %s", s.MacAddress)
	}
	return service.dao.UpdateServer(s)
}

//CreateServer create a server
func (service *DefaultService) CreateServer(s model.ServerConfig) (model.ServerConfig, error) {
	return service.dao.CreateServer(s)
}

//DeleteAll delete all the server config
func (service *DefaultService) DeleteAll() error {
	return service.dao.DeleteAll()
}
