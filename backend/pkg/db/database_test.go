package db_test

import (
	"bytes"
	"fmt"
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
		t.Errorf("Can not connect to the database %v, the error is %v", cfg.Database, err)
	}

	if len(servers) == 0 {
		t.Errorf("Expect more than 0 records return")
	}

	fmt.Println(servers[0].Hostname)
	fmt.Println(servers[0].MacAddress)

}

func connect() (*sql.DB, error) {
	return sql.Open("postgres", "host=localhost port=5432 user=postgres password=mysecretpassword dbname=pxeboot sslmode=disable")
}
