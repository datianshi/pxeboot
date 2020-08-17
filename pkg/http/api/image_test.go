package api_test

import (
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/http/api"
	"github.com/hooklift/iso9660"
	"os"
	"testing"
)

func TestImageUpload(t *testing.T){
	t.Skip("skipping test.")
	c := &config.Config{
		RootPath: "./fixture",
	}
	imageUploader := &api.ImageUploader{
		Cfg: c,
	}
	src, err := os.Open("/Users/sding/Downloads/VMware-VMvisor-Installer-7.0.0-15843807.x86_64.iso")
	if err != nil {
		fmt.Println(err)
	}
	r, err := iso9660.NewReader(src)
	if err != nil {
		fmt.Println(err)
	}
	imageUploader.Upload(r)

}