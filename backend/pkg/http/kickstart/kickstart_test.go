package kickstart_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/datianshi/pxeboot/pkg/http/kickstart"
	test_util "github.com/datianshi/pxeboot/pkg/testutil"
	"github.com/gorilla/mux"
)

func TestKickStart(t *testing.T) {
	router := mux.NewRouter()
	cfg, err := test_util.GetConfig("../../testutil/fixture/cfg.yaml")
	if err != nil {
		t.Fatalf("Can not process the config %v", err)
	}

	k := kickstart.NewKickStart(cfg.DHCPInterfaces[0])
	router.HandleFunc("/kickstart/{mac_address}/ks.cfg", k.Handler())
	r, err := http.NewRequest("GET", "/kickstart/00-50-56-82-61-7c/ks.cfg", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	real := string(body)
	expected := `vmaccepteula
clearpart --overwritevmfs --alldrives

# Set the root password for the DCUI and Tech Support Mode
rootpw VMware1!

# Install on the first local disk available on machine
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
server time.svc.pivotal.io
__NTP_CONFIG__

/sbin/chkconfig ntpd on

reboot
`
	if strings.Compare(real, expected) != 0 {
		t.Errorf("\nstart---%s---end\n not equal to \nstart---%s---end\n", real, expected)
	}
}
