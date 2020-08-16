package http_test

import (
	"bytes"
	"github.com/datianshi/pxeboot/pkg/config"
	h "github.com/datianshi/pxeboot/pkg/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var data string = `
dhcp_interface: ens224
bind_ip: 172.16.100.2
gateway: 10.65.101.1
netmask: 255.255.255.0
dns: 10.192.2.10
password: VMware1!
nics:
  00-50-56-82-70-2a:
    dhcp_ip: 172.16.100.100
    ip: 10.65.101.10
    hostname: vc-01.example.org
  00-50-56-82-61-7c:
    dhcp_ip: 172.16.100.101
    ip: 10.65.101.11
    hostname: vc-02.example.org
boot_file: efi/boot/bootx64.efi
lease_time: 500
root_path: ./fixture/image
boot_config_file: efi/boot/boot.cfg
ntp_server: time.svc.pivotal.io
kickstart_template: |
    vmaccepteula
    clearpart --overwritevmfs --alldrives

    # Set the root password for the DCUI and Tech Support Mode
    rootpw {{.Password}}

    # Install on the first local disk available on machine
    #install --firstdisk="DELLBOSS VD",Hypervisor_0,HV,usb,IDSDM --overwritevmfs --novmfsondisk
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
    restrict default kod nomodify notrap noquerynopeer
    restrict 127.0.0.1
    server {{.NTPServer}}
    __NTP_CONFIG__

    /sbin/chkconfig ntpd on

    reboot`

func TestKickStart(t *testing.T) {
	router := mux.NewRouter()
	var buf bytes.Buffer
	buf.WriteString(data)
	cfg, err := config.LoadConfig(&buf)
	if err != nil {
		t.Fatalf("Can not process the config %v", err)
	}

	k := h.Kickstart{
		R: router,
		C: cfg,
	}

	r, err := http.NewRequest("GET", "/kickstart/00-50-56-82-61-7c/ks.cfg", nil)
	w := httptest.NewRecorder()
	k.RegisterServerEndpoint()
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	real := string(body)
	expected := `# Accept the VMware End User License Agreement
  vmaccepteula
  clearpart --overwritevmfs --alldrives

  # Set the root password for the DCUI and Tech Support Mode
  rootpw VMware1!

  # Install on the first local disk available on machine
  #install --firstdisk="DELLBOSS VD",Hypervisor_0,HV,usb,IDSDM --overwritevmfs --novmfsondisk
  install --firstdisk --overwritevmfs

  # Set the network to DHCP on the first network adapter
  network --bootproto=static --addvmportgroup=1 --ip=10.65.101.11 --netmask=255.255.255.0 --gateway=10.65.101.1 --nameserver=10.192.2.10 --hostname=vc-02.example.org
  reboot

  %firstboot --interpreter=busybox
  vim-cmd hostsvc/enable_ssh
  vim-cmd hostsvc/start_ssh
  vim-cmd hostsvc/enable_esx_shell
  vim-cmd hostsvc/start_esx_shell
  cat > /etc/ntp.conf << __NTP_CONFIG__
  restrict default kod nomodify notrap noquerynopeer
  restrict 127.0.0.1
  server time.svc.pivotal.io
  __NTP_CONFIG__

  /sbin/chkconfig ntpd on

  reboot`
	if strings.Compare(real, expected) != 0 {
		t.Errorf("\n%s\n not equal to \n%s\nabc", real, expected)
	}
}