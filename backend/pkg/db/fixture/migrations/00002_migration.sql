-- +goose Up
-- SQL in this section is executed when the migration is applied.
insert into server(mac_address, ip, hostname, netmask, gateway, created_on)
values('00-50-56-82-70-2a', '10.65.123.20', 'vc-01.example.org', '255.255.255.0', '10.65.123.1', current_timestamp);