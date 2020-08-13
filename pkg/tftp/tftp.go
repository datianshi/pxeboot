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


// readHandler is called when client starts file download from server
func readHandler(filename string, rf io.ReaderFrom) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	n, err := rf.ReadFrom(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes sent\n", n)
	return nil
}

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
				r, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
					return err
				}
			}
			n, err := rf.ReadFrom(r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return err
			}
			fmt.Printf("%d bytes sent\n", n)
			return nil
		}
	}

	//Generate our own boot config file
	func bootConfig(cfg *config.Config, nic string) (io.Reader, error){
		fmt.Println(fmt.Sprintf("%s/%s", cfg.RootPath, cfg.BootConfigFile))
		origin_template, err:= os.Open(fmt.Sprintf("%s/%s", cfg.RootPath, cfg.BootConfigFile))
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		kickstartUrl := fmt.Sprintf("http://%s/kickstart/%s/ks.cfg", cfg.BindIP, nic)
		//prefix := fmt.Sprintf("01-%s/%s", nic, SYMLINK_PER_SERVER_DIR)
		scanner := bufio.NewScanner(origin_template)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "kernelopt=") {
				re := regexp.MustCompile(`^(kernelopt=).*$`)
				t := re.ReplaceAllString(line, fmt.Sprintf("${1}ks=%s", kickstartUrl))
				_, err := buf.Write([]byte(fmt.Sprintf("%s\n", t)))
				if err != nil {
					return nil, err
				}
			} else {
			//} else if strings.Contains(line, "prefix=") {
			//	re := regexp.MustCompile(`^(prefix=).*$`)
			//	t := re.ReplaceAllString(line, fmt.Sprintf("${1}%s", prefix))
			//	//fmt.Println(t)
			//	dst.Write([]byte(fmt.Sprintf("%s\n", t)))
			//} else {
				buf.Write([]byte(fmt.Sprintf("%s\n", line)))
			}

		}
		return &buf, nil
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
		return writeHandler(fmt.Sprintf("%s/%s", cfg.RootPath, filename), rf)
	}
}



func Start(cfg *config.Config) {
	s := tftp.NewServer(pxeReadHandler(cfg), pxeWriteHandler(cfg))
	s.SetTimeout(5 * time.Second)
	err := s.ListenAndServe(fmt.Sprintf("%s:69", cfg.BindIP))
	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}

}