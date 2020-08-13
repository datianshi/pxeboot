package tftp

import (
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/pin/tftp"
	"io"
	"os"
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

func rootPathReadHandler(rootPath string) func(filename string, rf io.ReaderFrom) error {
	return func(filename string, rf io.ReaderFrom) error{
		return readHandler(fmt.Sprintf("%s/%s", rootPath, filename), rf)
	}
}

func rootPathWriteHandler(rootPath string) func(filename string, rf io.WriterTo) error {
	return func(filename string, rf io.WriterTo) error{
		return writeHandler(fmt.Sprintf("%s/%s", rootPath, filename), rf)
	}
}



func Start(cfg *config.Config) {
	s := tftp.NewServer(rootPathReadHandler(cfg.RootPath), rootPathWriteHandler(cfg.RootPath))
	s.SetTimeout(5 * time.Second)
	err := s.ListenAndServe(fmt.Sprintf("%s:69", cfg.BindIP))
	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}

}