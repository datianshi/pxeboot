package test_util

import (
	"github.com/datianshi/pxeboot/pkg/config"
	"os"
)

func GetConfig(filePath string) ( *config.Config, error){
	configfile, err := os.Open(filePath)
	defer configfile.Close()
	if err != nil {
		return nil, err
	}
	cfg, err := config.LoadConfig(configfile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}