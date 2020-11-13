-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE SERVER (
	id serial PRIMARY KEY,
	mac_address VARCHAR ( 50 ) UNIQUE NOT NULL,
	ip VARCHAR ( 50 ) NOT NULL,
	hostname VARCHAR ( 50 ) NOT NULL,
    netmask VARCHAR ( 50 ) NOT NULL,
    gateway VARCHAR ( 50 ),
	created_on TIMESTAMP NOT NULL
);
