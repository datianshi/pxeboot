package tftp

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/pin/tftp"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
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

func pxeReadHandler(cfg *config.Config) func(filename string, rf io.ReaderFrom) error {
		//Not read from the file. Generate it dynamically
		return func(filename string, rf io.ReaderFrom) error{
			exist, nic := getNic(filename, cfg)
			var r io.Reader
			var err error
			if exist {
				r, err = bootConfig(cfg, nic)
				if (err != nil) {
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
func bootConfig(cfg *config.Config, nic string) (*hackReader, error){
	//fmt.Println(fmt.Sprintf("%s/%s", cfg.RootPath, cfg.BootConfigFile))
	origin_template, err:= os.Open(fmt.Sprintf("%s/%s", cfg.RootPath, cfg.BootConfigFile))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	kickstartUrl := fmt.Sprintf("http://%s/kickstart/%s/ks.cfg", cfg.BindIP, nic)
	scanner := bufio.NewScanner(origin_template)
	s := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "kernelopt=") {
			re := regexp.MustCompile(`^(kernelopt=).*$`)
			t := re.ReplaceAllString(line, fmt.Sprintf("${1}ks=%s", kickstartUrl))
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

func getNic(filename string, cfg *config.Config) (bool, string) {
	for k, _ := range cfg.Nics {
		if strings.Contains(filename, k) {
			return true, k
		}
	}
	return false, ""
}

func pxeWriteHandler(cfg *config.Config) func(filename string, rf io.WriterTo) error {
	return func(filename string, rf io.WriterTo) error{
		fmt.Printf("Write File %s\n", filename)
		return writeHandler(fmt.Sprintf("%s/%s", cfg.RootPath, filename), rf)
	}
}



func Start(cfg *config.Config) {
	var port int
	if cfg.TFTPServerPort != 0 {
		port = cfg.TFTPServerPort
	} else {
		port = 67
	}
	s := tftp.NewServer(pxeReadHandler(cfg), pxeWriteHandler(cfg))
	s.SetTimeout(5 * time.Second)
	err := s.ListenAndServe(fmt.Sprintf("%s:%d", cfg.BindIP, port))
	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}

}