package tftp

import (
	"fmt"
	"io"
	"github.com/pin/tftp"
	"os"
	"time"
)


// readHandler is called when client starts file download from server
var root_path = "/netboot/tftp/"
func readHandler(filename string, rf io.ReaderFrom) error {
	file, err := os.Open(fmt.Sprintf("%s/%s", root_path, filename))
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
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", root_path, filename), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
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

func Start() {
	// use nil in place of handler to disable read or write operations
	s := tftp.NewServer(readHandler, writeHandler)
	s.SetTimeout(5 * time.Second) // optional
	err := s.ListenAndServe("172.16.100.2:69") // blocks until s.Shutdown() is called
	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}
}