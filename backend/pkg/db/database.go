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
	if rows, err = database.Query("select id, gateway, hostname, ip, netmask, mac_address from server"); err != nil {
		return nil, err
	}
	var servers []model.ServerConfig
	for rows.Next() {
		var server model.ServerConfig
		if err = rows.Scan(&server.ID, &server.Gateway, &server.Hostname, &server.Ip, &server.Netmask, &server.MacAddress); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}

//FindServer Retrieve the server from database
func (db *Database) FindServer(macAddress string) (model.ServerConfig, error) {
	sqlStatement := `select id, gateway, hostname, ip, netmask, mac_address FROM server WHERE mac_address=$1;`
	var err error
	var server model.ServerConfig
	var database *sql.DB
	if database, err = db.openConnection(); err != nil {
		return model.ServerConfig{}, err
	}
	row := database.QueryRow(sqlStatement, macAddress)
	if err = row.Scan(&server.ID, &server.Gateway, &server.Hostname, &server.Ip, &server.Netmask, &server.MacAddress); err != nil {
		return model.ServerConfig{}, err
	}
	return server, nil
}
