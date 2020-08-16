package ui

import (
	"github.com/gorilla/mux"
	"net/http"
	"html/template"
)

func RegisterUITemplate(r *mux.Router) error{
	tmpl, err := template.ParseFiles("./public/templates/index.html")
	if err != nil {
		return err
	}
	r.HandleFunc("/", HomeHandler(tmpl))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("public/assets"))))
	r.HandleFunc("/static/", HomeHandler(tmpl))
	return nil
}

func HomeHandler(tmpl *template.Template) http.HandlerFunc{
	t, _ := template.ParseFiles("./public/templates/index.html")
	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, nil)
	}
}
