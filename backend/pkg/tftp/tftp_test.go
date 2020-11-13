package tftp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/model"
	"github.com/datianshi/pxeboot/pkg/nic"
	"github.com/datianshi/pxeboot/pkg/nic/nicfakes"
)

var data = `dhcp_interface: ens224
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
ntp_server: time.svc.pivotal.io
boot_config_file: efi/boot/boot.cfg
`

func TestReadBootConfigFile(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(data)
	cfg, err := config.LoadConfig(&buf)
	if err != nil {
		t.Errorf("Can not process the config")
	}

	fileName := "/10-00-50-56-82-70-2a/boot.cfg"

	var remoteBuf bytes.Buffer
	pxeReadHandler(cfg, getFakeNicService())(fileName, &remoteBuf)

	real := remoteBuf.String()
	expected := `bootstate=0
title=Loading ESXi installer
timeout=5
prefix=
kernel=/b.b00
kernelopt=ks=http://172.16.100.2/kickstart/00-50-56-82-70-2a/ks.cfg
modules=/jumpstrt.gz --- /useropts.gz --- /features.gz --- /k.b00 --- /uc_intel.b00 --- /uc_amd.b00 --- /uc_hygon.b00 --- /procfs.b00 --- /vmx.v00 --- /vim.v00 --- /tpm.v00 --- /sb.v00 --- /s.v00 --- /bnxtnet.v00 --- /bnxtroce.v00 --- /brcmfcoe.v00 --- /brcmnvme.v00 --- /elxiscsi.v00 --- /elxnet.v00 --- /i40en.v00 --- /i40iwn.v00 --- /iavmd.v00 --- /igbn.v00 --- /iser.v00 --- /ixgben.v00 --- /lpfc.v00 --- /lpnic.v00 --- /lsi_mr3.v00 --- /lsi_msgp.v00 --- /lsi_msgp.v01 --- /lsi_msgp.v02 --- /mtip32xx.v00 --- /ne1000.v00 --- /nenic.v00 --- /nfnic.v00 --- /nhpsa.v00 --- /nmlx4_co.v00 --- /nmlx4_en.v00 --- /nmlx4_rd.v00 --- /nmlx5_co.v00 --- /nmlx5_rd.v00 --- /ntg3.v00 --- /nvme_pci.v00 --- /nvmerdma.v00 --- /nvmxnet3.v00 --- /nvmxnet3.v01 --- /pvscsi.v00 --- /qcnic.v00 --- /qedentv.v00 --- /qedrntv.v00 --- /qfle3.v00 --- /qfle3f.v00 --- /qfle3i.v00 --- /qflge.v00 --- /rste.v00 --- /sfvmk.v00 --- /smartpqi.v00 --- /vmkata.v00 --- /vmkfcoe.v00 --- /vmkusb.v00 --- /vmw_ahci.v00 --- /crx.v00 --- /elx_esx_.v00 --- /btldr.v00 --- /esx_dvfi.v00 --- /esx_ui.v00 --- /esxupdt.v00 --- /tpmesxup.v00 --- /weaselin.v00 --- /loadesx.v00 --- /lsuv2_hp.v00 --- /lsuv2_in.v00 --- /lsuv2_ls.v00 --- /lsuv2_nv.v00 --- /lsuv2_oe.v00 --- /lsuv2_oe.v01 --- /lsuv2_oe.v02 --- /lsuv2_sm.v00 --- /native_m.v00 --- /qlnative.v00 --- /vdfs.v00 --- /vmware_e.v00 --- /vsan.v00 --- /vsanheal.v00 --- /vsanmgmt.v00 --- /tools.t00 --- /xorg.v00 --- /imgdb.tgz --- /imgpayld.tgz
build=7.0.0-1.0.15843807
updated=0
`
	if strings.Compare(real, expected) != 0 {
		t.Errorf("%s\n not equal to %s\n", real, expected)
	}
}

func TestReadRegularFile(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(data)
	cfg, err := config.LoadConfig(&buf)
	if err != nil {
		t.Errorf("Can not process the config")
	}

	fileName := "afile"

	var remoteBuf bytes.Buffer
	pxeReadHandler(cfg, getFakeNicService())(fileName, &remoteBuf)

	real := remoteBuf.String()
	expected := `this is a file`
	if strings.Compare(real, expected) != 0 {
		t.Errorf("%s\n not equal to %s\n", real, expected)
	}

}

func getFakeNicService() nic.Service {
	nicService := nicfakes.FakeService{}
	nicService.GetServersStub = func() ([]model.ServerConfig, error) {
		return []model.ServerConfig{
			{
				MacAddress: "00-50-56-82-70-2a",
			},
			{
				MacAddress: "00-50-56-82-61-7c",
			},
		}, nil
	}
	return &nicService
}
