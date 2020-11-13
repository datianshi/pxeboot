package tftp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/model"
	"github.com/datianshi/pxeboot/pkg/nic"
	"github.com/pin/tftp"
)

// writeHandler is called when client starts file upload to server
func writeHandler(filename string, wt io.WriterTo) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	n, err := wt.WriteTo(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes received\n", n)
	return nil
}

type hackReader struct {
	bytes.Buffer
	size int64
}

//hack to match Seek interface. TFTP server library has a type cast
func (reader *hackReader) Seek(offset int64, whence int) (int64, error) {
	return reader.size, nil
}

func pxeReadHandler(cfg config.Config, nicService nic.Service) func(filename string, rf io.ReaderFrom) error {
	//Not read from the file. Generate it dynamically
	return func(filename string, rf io.ReaderFrom) error {
		servers, err := nicService.GetServers()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not get list of servers %s", err)
			return err
		}
		exist, nic := getNic(filename, servers)
		var r io.Reader
		if exist {
			r, err = bootConfig(cfg, nic)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Can not generate boot config %s", err)
			}
		} else {
			r, err = os.Open(fmt.Sprintf("%s/%s", cfg.RootPath, filename))
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return err
			}
		}
		_, err = rf.ReadFrom(r)
		if err != nil {
			//fmt.Fprintf(os.Stderr, "%v\n", err)
			return err
		}
		//fmt.Printf("%d bytes sent\n", n)
		return nil
	}
}

//Generate our own boot config file
func bootConfig(cfg config.Config, nic string) (*hackReader, error) {
	//fmt.Println(fmt.Sprintf("%s/%s", cfg.RootPath, cfg.BootConfigFile))
	originTemplate, err := os.Open(fmt.Sprintf("%s/%s", cfg.RootPath, cfg.BootConfigFile))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	kickstartURL := fmt.Sprintf("http://%s/kickstart/%s/ks.cfg", cfg.BindIP, nic)
	scanner := bufio.NewScanner(originTemplate)
	s := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "kernelopt=") {
			re := regexp.MustCompile(`^(kernelopt=).*$`)
			t := re.ReplaceAllString(line, fmt.Sprintf("${1}ks=%s", kickstartURL))
			write, err := buf.Write([]byte(fmt.Sprintf("%s\n", t)))
			s = s + write
			if err != nil {
				return nil, err
			}
		} else {
			write, _ := buf.Write([]byte(fmt.Sprintf("%s\n", line)))
			s = s + write
		}

	}
	ret := &hackReader{
		buf,
		int64(s),
	}
	return ret, nil
}

func getNic(filename string, servers []model.ServerConfig) (bool, string) {
	for _, v := range servers {
		if strings.Contains(filename, v.MacAddress) {
			return true, v.MacAddress
		}
	}
	return false, ""
}

func pxeWriteHandler(cfg config.Config) func(filename string, rf io.WriterTo) error {
	return func(filename string, rf io.WriterTo) error {
		fmt.Printf("Write File %s\n", filename)
		return writeHandler(fmt.Sprintf("%s/%s", cfg.RootPath, filename), rf)
	}
}

//Start start the tftp server
func Start(cfg config.Config, nicService nic.Service) {
	var port int
	if cfg.TFTPServerPort != 0 {
		port = cfg.TFTPServerPort
	} else {
		port = 69
	}
	s := tftp.NewServer(pxeReadHandler(cfg, nicService), pxeWriteHandler(cfg))
	s.SetTimeout(5 * time.Second)
	err := s.ListenAndServe(fmt.Sprintf("%s:%d", cfg.BindIP, port))
	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}

}
