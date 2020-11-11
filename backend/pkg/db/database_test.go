package db_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/db"
)

var data string = `
database:
  username: postgres
  password: mysecretpassword
  host: 192.168.86.103
  port: 5432
  database_name: pxeboot
`

func TestDBConnection(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(data)
	cfg, err := config.LoadConfig(&buf)
	if err != nil {
		t.Fatalf("Can not process the config %v", err)
	}

	database := db.NewDatabase(cfg)
	servers, err := database.GetServers()
	if err != nil {
		t.Errorf("Can not connect to the database %v, the error is %v", cfg.Database, err)
	}

	if len(servers) == 0 {
		t.Errorf("Expect more than 0 records return")
	}

	fmt.Println(servers[0].Hostname)
	fmt.Println(servers[0].MacAddress)

}
