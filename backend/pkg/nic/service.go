package nic

import (
	"github.com/datianshi/pxeboot/pkg/db"
	"github.com/datianshi/pxeboot/pkg/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Service

//Service service to control servers
type Service interface {
	GetServers() ([]model.ServerConfig, error)
	FindServer(string) (model.ServerConfig, error)
	DeleteServer(string) error
	UpdateServer(model.ServerConfig) (model.ServerConfig, error)
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
	return nil, nil
}

//FindServer based on the mac address
func (service *DefaultService) FindServer(string) (model.ServerConfig, error) {
	return model.ServerConfig{}, nil
}

//DeleteServer based on the mac address
func (service *DefaultService) DeleteServer(string) error {
	return nil
}

//UpdateServer Patch the config
func (service *DefaultService) UpdateServer(model.ServerConfig) (model.ServerConfig, error) {
	return model.ServerConfig{}, nil
}

//CreateServer create a server
func (service *DefaultService) CreateServer(model.ServerConfig) (model.ServerConfig, error) {
	return model.ServerConfig{}, nil
}

//DeleteAll delete all the server config
func (service *DefaultService) DeleteAll() error {
	return nil
}
