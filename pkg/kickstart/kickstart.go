package kickstart

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"time"
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
install --firstdisk="DELLBOSS VD",Hypervisor_0,HV,usb,IDSDM --overwritevmfs --novmfsondisk

# Set the network to DHCP on the first network adapter
network --bootproto=static --addvmportgroup=1 --ip={{.IP}} --netmask=255.255.255.0 --gateway=10.65.64.1 --nameserver=10.192.2.10 --hostname=pao-esx-01.thor.pvd.pez.pivotal.io
reboot

%firstboot --interpreter=busybox
vim-cmd hostsvc/enable_ssh
vim-cmd hostsvc/start_ssh
vim-cmd hostsvc/enable_esx_shell
vim-cmd hostsvc/start_esx_shell
cat > /etc/ntp.conf << __NTP_CONFIG__
restrict default kod nomodify notrap noquerynopeer
restrict 127.0.0.1
server 10.201.194.5
__NTP_CONFIG__

/sbin/chkconfig ntpd on

reboot
`

type Kickstart struct {
	R *mux.Router
}



func (k *Kickstart) Start(){
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: k.R, // Pass our instance of gorilla/mux in.
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
}

func (k *Kickstart) RegisterServerEndpoint(ip string) {
	go func(){
		k.R.HandleFunc("/{ip_address}/ks.cfg", kickstartHandler)
	}()
}

type info struct {
	IP string
}

func kickstartHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ip := vars["ip_address"]
	i := info{
		IP: ip,
	}
	t, _ := template.New("").Parse(kickstart_template)
	err := t.ExecuteTemplate(w, "", i)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}