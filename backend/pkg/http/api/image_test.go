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
	c := &config.Config{
		RootPath: "./fixture",
	}
	imageUploader := &api.ImageUploader{
		Cfg: c,
	}
	src, err := os.Open("./fixture/test.iso")
	if err != nil {
		fmt.Println(err)
	}
	r, err := iso9660.NewReader(src)
	if err != nil {
		fmt.Println(err)

	}
	defer cleanUp()
	imageUploader.Upload(r)
}

func cleanUp() {
	os.Remove("./fixture/file1.")
	os.Remove("./fixture/file2.")
}