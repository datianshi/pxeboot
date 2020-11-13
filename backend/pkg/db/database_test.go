package db_test

import (
	"bytes"
	"testing"
	"time"

	"database/sql"

	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/db"
	"github.com/datianshi/pxeboot/pkg/model"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/pressly/goose"
)

var data string = `
database:
  username: postgres
  password: mysecretpassword
  host: localhost
  port: 5432
  database_name: pxeboot
`

func TestDBConnection(t *testing.T) {
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Username("postgres").
		Password("mysecretpassword").
		Database("pxeboot").
		StartTimeout(45 * time.Second))

	err := postgres.Start()
	defer func() {
		if err := postgres.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	var buf bytes.Buffer
	buf.WriteString(data)
	cfg, err := config.LoadConfig(&buf)
	if err != nil {
		t.Fatalf("Can not process the config %v", err)
	}
	//Load Data
	loader, err := connect()
	if err != nil {
		t.Fatal(err)
	}

	if err := goose.Up(loader, "./fixture/migrations"); err != nil {
		t.Fatal(err)
	}

	//Test Get All servers
	database := db.NewDatabase(cfg.Database)

	servers, err := database.GetServers()

	if err != nil {
		t.Errorf("Can not get servers %v", err)
	}

	if len(servers) != 2 {
		t.Errorf("Expect 2 records return but got %d", len(servers))
	}

	//Test Find Server
	server, err := database.FindServer("00-50-56-82-78-2a")
	if err != nil {
		t.Errorf("Can not get server with mac address 00-50-56-82-78-2a %v", err)
	}
	if server.Ip != "10.65.123.21" {
		t.Errorf("Can not retrieve the correct server with mac address 00-50-56-82-78-2a err %v", err)
	}
	//Test Delete Server
	err = database.DeleteServer("00-50-56-82-78-2a")
	if err != nil {
		t.Errorf("Can not delete server with mac address 00-50-56-82-78-2a %v", err)
	}
	_, err = database.FindServer("00-50-56-82-78-2a")
	if err != sql.ErrNoRows {
		t.Error("Expect server with mac address 00-50-56-82-78-2a already deleted")
	}
	//Test Create Server

	nic := model.ServerConfig{
		Ip:         "192.168.0.200",
		Gateway:    "192.168.0.1",
		Netmask:    "255.255.255.0",
		MacAddress: "00-50-56-82-78-2a",
		Hostname:   "test-server.org",
	}
	createdServer, err := database.CreateServer(nic)
	if err != nil {
		t.Errorf("Expect no error when create a server in database but got err %v", err)
	}
	if createdServer.ID != 3 {
		t.Errorf("Expect server id to be 2 but got %d", createdServer.ID)
	}
	servers, _ = database.GetServers()
	if len(servers) != 2 {
		t.Errorf("Expect 2 records return but got %d", len(servers))
	}

	//Update Server
	updated := model.ServerConfig{
		Ip:         "192.168.0.200",
		Gateway:    "192.168.0.1",
		Netmask:    "255.255.255.0",
		MacAddress: "00-50-56-82-70-2a",
		Hostname:   "updated.org",
	}
	if err = database.UpdateServer(updated); err != nil {
		t.Errorf("Expect no error when update a server in database but got err %v", err)
	}
	server, _ = database.FindServer("00-50-56-82-70-2a")
	if server.Hostname != "updated.org" {
		t.Errorf("Expect the host name is updated but got %s", server.Hostname)
	}

	//Delete All Servers
	if err = database.DeleteAll(); err != nil {
		t.Errorf("Failed to delete all the recrods err %v", err)
	}
	servers, _ = database.GetServers()
	if len(servers) != 0 {
		t.Errorf("Expect no row, but got %d rows", len(servers))
	}
}

func connect() (*sql.DB, error) {
	return sql.Open("postgres", "host=localhost port=5432 user=postgres password=mysecretpassword dbname=pxeboot sslmode=disable")
}
