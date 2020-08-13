package util_test

import (
	"bytes"
	"fmt"
	"github.com/datianshi/pxeboot/pkg/util"
	"os"
	"strings"
	"testing"
)

func TestDotToDash(t *testing.T) {
	ret := util.Dot_To_Dash("10.20.76.700")
	if ret != "10-20-76-700" {
		t.Errorf("want 10-20-76-700, but got %s", ret)
	}
}

func TestColonToDash(t *testing.T) {
	ret := util.Colon_To_Dash("00:50:56:82:70:2a")
	if ret != "00-50-56-82-70-2a" {
		t.Errorf("want 00-50-56-82-70-2a, but got %s", ret)
	}
}

func TestDashToDot(t *testing.T) {
	ret := util.Dash_To_Dot("10-20-76-700")
	if ret != "10.20.76.700" {
		t.Errorf("want 10.20.76.700, but got %s", ret)
	}
}


func TestBootConfigFile(t *testing.T) {
	file, err := os.Open("./fixture/boot.cfg")
	if err != nil {
		t.Errorf("Can not open the file %s", "./fixture/boot.cfg")
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	util.BootConfigFile(file, buf, "http://172.16.100.2/static/ks.cfg")
	real := buf.String()
	expected := `bootstate=0
title=Loading ESXi installer
timeout=5
prefix=images
kernel=b.b00
kernelopt=http://172.16.100.2/static/ks.cfg
modules=/jumpstrt.gz --- /useropts.gz --- /features.gz --- /k.b00 --- /uc_intel.b00 --- /uc_amd.b00 --- /uc_hygon.b00 --- /procfs.b00 --- /vmx.v00 --- /vim.v00 --- /tpm.v00 --- /sb.v00 --- /s.v00 --- /bnxtnet.v00 --- /bnxtroce.v00 --- /brcmfcoe.v00 --- /brcmnvme.v00 --- /elxiscsi.v00 --- /elxnet.v00 --- /i40en.v00 --- /i40iwn.v00 --- /iavmd.v00 --- /igbn.v00 --- /iser.v00 --- /ixgben.v00 --- /lpfc.v00 --- /lpnic.v00 --- /lsi_mr3.v00 --- /lsi_msgp.v00 --- /lsi_msgp.v01 --- /lsi_msgp.v02 --- /mtip32xx.v00 --- /ne1000.v00 --- /nenic.v00 --- /nfnic.v00 --- /nhpsa.v00 --- /nmlx4_co.v00 --- /nmlx4_en.v00 --- /nmlx4_rd.v00 --- /nmlx5_co.v00 --- /nmlx5_rd.v00 --- /ntg3.v00 --- /nvme_pci.v00 --- /nvmerdma.v00 --- /nvmxnet3.v00 --- /nvmxnet3.v01 --- /pvscsi.v00 --- /qcnic.v00 --- /qedentv.v00 --- /qedrntv.v00 --- /qfle3.v00 --- /qfle3f.v00 --- /qfle3i.v00 --- /qflge.v00 --- /rste.v00 --- /sfvmk.v00 --- /smartpqi.v00 --- /vmkata.v00 --- /vmkfcoe.v00 --- /vmkusb.v00 --- /vmw_ahci.v00 --- /crx.v00 --- /elx_esx_.v00 --- /btldr.v00 --- /esx_dvfi.v00 --- /esx_ui.v00 --- /esxupdt.v00 --- /tpmesxup.v00 --- /weaselin.v00 --- /loadesx.v00 --- /lsuv2_hp.v00 --- /lsuv2_in.v00 --- /lsuv2_ls.v00 --- /lsuv2_nv.v00 --- /lsuv2_oe.v00 --- /lsuv2_oe.v01 --- /lsuv2_oe.v02 --- /lsuv2_sm.v00 --- /native_m.v00 --- /qlnative.v00 --- /vdfs.v00 --- /vmware_e.v00 --- /vsan.v00 --- /vsanheal.v00 --- /vsanmgmt.v00 --- /tools.t00 --- /xorg.v00 --- /imgdb.tgz --- /imgpayld.tgz
build=7.0.0-1.0.15843807
updated=0
`
	fmt.Printf(expected)
	fmt.Printf(real)
	if strings.Compare(real, expected) != 0 {
		t.Errorf("%s\n not equal to %s\n", real, expected)
	}
}
