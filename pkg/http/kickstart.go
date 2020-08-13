package http

import (
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/gorilla/mux"
	"net/http"
	"text/template"
)

const kickstart_template = `
#
# Sample scripted installation file
#

# Accept the VMware End User License Agreement
vmaccepteula
clearpart --overwritevmfs --alldrives

# Set the root password for the DCUI and Tech Support Mode
rootpw VMware1!

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

reboot
`

type Kickstart struct {
	R *mux.Router
	C *config.Config
}


func (k *Kickstart) RegisterServerEndpoint() {
	k.R.HandleFunc("/kickstart/{mac_address}/ks.cfg", k.handler())
}

type info struct {
	IP string
	Gateway string
	NetMask string
	NameServer string
	HostName string
	NTPServer string
}

func (k *Kickstart) handler() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		mac_address := vars["mac_address"]
		serverConfig, found := k.C.Nics[mac_address]
		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		i := info{
			IP: serverConfig.Ip,
			HostName: serverConfig.Hostname,
			Gateway: k.C.Gateway,
			NetMask: k.C.Netmask,
			NameServer: k.C.DNS,
			NTPServer: k.C.NTPServer,
		}
		t, _:= template.New("").Parse(kickstart_template)
		err := t.ExecuteTemplate(w, "", i)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
}