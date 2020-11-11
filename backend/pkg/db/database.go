package db

import (
	"database/sql"
	"fmt"

	"github.com/datianshi/pxeboot/pkg/config"
	_ "github.com/lib/pq"
)

type Database struct {
	cfg *config.Config
}

func NewDatabase(cfg *config.Config) *Database {
	return &Database{
		cfg: cfg,
	}
}

func (db *Database) openConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.cfg.Database.Host, db.cfg.Database.Port, db.cfg.Database.Username, db.cfg.Database.Password, db.cfg.Database.DatabaseName)

	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	return database, nil
}

func (db *Database) GetServers() ([]config.ServerConfig, error) {
	var err error
	var database *sql.DB
	var rows *sql.Rows
	if database, err = db.openConnection(); err != nil {
		return nil, err
	}
	if rows, err = database.Query("select gateway, hostname, ip, netmask, mac_address from server"); err != nil {
		return nil, err
	}
	var servers []config.ServerConfig
	for rows.Next() {
		var server config.ServerConfig
		if err = rows.Scan(&server.Gateway, &server.Hostname, &server.Ip, &server.Netmask, &server.MacAddress); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, nil

}
