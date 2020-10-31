package api

import (
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/hooklift/iso9660"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type ImageUploader struct {
	Cfg *config.DHCPInterface
}

func (iu *ImageUploader) UploadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//max 2GB
		r.ParseMultipartForm(2048)
		file, _, err := r.FormFile("image")
		if err != nil {
			iu.handleError(w, err)
		}
		defer file.Close()
		reader, err := iso9660.NewReader(file)
		if err != nil {
			iu.handleError(w, err)
		}
		if err := iu.Upload(reader); err != nil {
			iu.handleError(w, err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (iu *ImageUploader) handleError(w http.ResponseWriter, err error) {
	w.Write([]byte(err.Error()))
	w.WriteHeader(http.StatusBadRequest)
	panic(err)
}

func (iu *ImageUploader) Upload(r *iso9660.Reader) error{
	for {
		f, err := r.Next()
		if err == io.EOF {
			break
		}

		fp := filepath.Join(iu.Cfg.RootPath, f.Name())
		if f.IsDir() {
			if err := os.MkdirAll(fp, f.Mode()); err != nil {
				return err
			}
			continue
		}

		parentDir, _ := filepath.Split(fp)
		if err := os.MkdirAll(parentDir, f.Mode()); err != nil {
			return err
		}

		fmt.Printf("Extracting %s...\n", fp)
		freader := f.Sys().(io.Reader)
		ff, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(ff, freader)

		cerr := ff.Close() // With or without an error from Copy, we want to attempt Close.

		if err != nil { // Panic with Copy's err.
			return err
		} else if cerr != nil {
			return err
		}
	}
	return nil
}
