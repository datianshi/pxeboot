package db

import (
	"database/sql"
	"fmt"

	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/model"
	_ "github.com/lib/pq"
)

//Database Object
type Database struct {
	cfg config.Database
}

//NewDatabase Create a New Database Object
func NewDatabase(cfg config.Database) *Database {
	return &Database{
		cfg: cfg,
	}
}

//Open connection
func (db *Database) openConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.cfg.Host, db.cfg.Port, db.cfg.Username, db.cfg.Password, db.cfg.DatabaseName)

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

//GetServers Retrieve All the servers
func (db *Database) GetServers() ([]model.ServerConfig, error) {
	var err error
	var database *sql.DB
	var rows *sql.Rows
	if database, err = db.openConnection(); err != nil {
		return nil, err
	}
	if rows, err = database.Query("select gateway, hostname, ip, netmask, mac_address from server"); err != nil {
		return nil, err
	}
	var servers []model.ServerConfig
	for rows.Next() {
		var server model.ServerConfig
		if err = rows.Scan(&server.Gateway, &server.Hostname, &server.Ip, &server.Netmask, &server.MacAddress); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, nil

}
