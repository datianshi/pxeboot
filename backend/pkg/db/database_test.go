package db_test

import (
	"bytes"
	"testing"
	"time"

	"database/sql"

	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/db"
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

	// Do test logic

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

	//Test
	database := db.NewDatabase(cfg.Database)

	servers, err := database.GetServers()

	if err != nil {
		t.Errorf("Can not get servers %v", err)
	}

	if len(servers) != 2 {
		t.Errorf("Expect 2 records return but got", len(servers))
	}

	server, err := database.FindServer("00-50-56-82-78-2a")
	if err != nil {
		t.Errorf("Can not get server with mac address 00-50-56-82-78-2a %v", err)
	}
	if server.Ip != "10.65.123.21" {
		t.Errorf("Can not retrieve the correct server with mac address 00-50-56-82-78-2a err %v", err)
	}

}

func connect() (*sql.DB, error) {
	return sql.Open("postgres", "host=localhost port=5432 user=postgres password=mysecretpassword dbname=pxeboot sslmode=disable")
}
