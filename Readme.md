# Golang PXE Servers for ESXi

Goals:

* Able to do EFI pxe boot instead of legacy bios mode
* All the pieces in one program
* Able to bootstrap multiple servers concurrently
* Able to interpolate kick start and boot config according each server's nic

## What does this project cover 

* A DHCP Server ([DHCP Lib](https://github.com/insomniacslk/dhcp))
* A TFTP Server ([TFTP Lib](https://github.com/pin/tftp))
* An HTTP Server for kickstart

## What users need to provide

* A folder that extracts the esxi iso image
* [A Config File](./pxe-sample.yaml)
    
```yaml
#DHCP Bind Interface
dhcp_interface: vboxnet5
#DHCP Server IP
bind_ip: 192.168.99.1
#Management IP which sit on an management network
management_ip: 192.168.62.1
#DHCP Range
dhcp_range: 172.16.100.10 - 172.16.100.100
#Default 67
#dhcp_server_port: 10000
#Default 69
#tftp_server_port: 10001
#Default 80
#http_port: 8089
#ESXi Servers DNS
dns: 10.192.2.10
#ESXi Password
password: VMware1!
nics:
  #ESXi Nic Address
  00-50-56-82-70-2a:
    #ESXi Static IP Assignment
    ip: 10.65.101.10
    #ESXi Host Name Assignment
    hostname: vc-01.example.org
    #ESXi Servers Gateway
    gateway: 10.65.101.1
    #ESXi Servers Netmask
    netmask: 255.255.255.0
  00-50-56-82-d8-86:
    ip: 10.65.101.11
    hostname: vc-02.example.org
    #ESXi Servers Gateway
    gateway: 10.65.101.1
    #ESXi Servers Netmask
    netmask: 255.255.255.0
#Boot File Path. Relative to the extracted iso image folder
boot_file: efi/boot/bootx64.efi
#DHCP Lease time in second
lease_time: 500
#TFTP Root Path, which is the extracted iso image folder
root_path: /Users/sding/image
#ESXi ntp server
ntp_server: time.svc.pivotal.io
#boot config file, which comes with iso image. Used to interpolate values
boot_config_file: efi/boot/boot.cfg
#kick start template with variables -> password, ip, netmask, gateway, nameserver, ntpserver and hostname
kickstart_template: |
  #
  # Sample scripted installation file
  #

  # Accept the VMware End User License Agreement
    vmaccepteula
    clearpart --overwritevmfs --alldrives

    # Set the root password for the DCUI and Tech Support Mode
    rootpw {{.Password}}

    # Install on the first local disk available on machine
    install --firstdisk --overwritevmfs

    # Set the network to DHCP on the first network adapter
    network --bootproto=static --addvmportgroup=1 --ip={{.IP}} --netmask={{.NetMask}} --gateway={{.Gateway}} --nameserver={{.NameServer}} --hostname={{.HostName}}
    reboot

    %firstboot --interpreter=busybox
    vim-cmd hostsvc/enable_ssh
    vim-cmd hostsvc/start_ssh
    vim-cmd hostsvc/enable_esx_shell
    vim-cmd hostsvc/start_esx_shell
    cat > /etc/ntp.conf << __NTP_CONFIG__
    server {{.NTPServer}}
    __NTP_CONFIG__

    /sbin/chkconfig ntpd on

    reboot
```
## Run

**YES** Program needs to run in the same layer2 network with Servers

```go run cmd/pxe/main.go -config pxe.yaml```

## API  

### Add Mac Address

```
curl -X POST http://Management_IP/api/conf/nic --data '{"mac_address":"00:50:56:82:75:e5", "hostname":"server1", "ip":"10.65.101.10"}'
```
### Delete Mac Address

```
curl -X DELETE http://Management_IP/api/conf/nic/${mac_address}
```

## UI

```
http://Management_IP/
```

## A lot of future work :) Foundation exist !~ 