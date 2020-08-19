package api

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func (api *API) RegisterUITemplate(r *mux.Router) error{
	var indexHtml string
	var err error
	var tmpl *template.Template = template.New("Html Template")
	if indexHtml, err = api.htmlBox.FindString("templates/index.html.tpl"); err != nil {
		return err
	}
	if tmpl, err = tmpl.Parse(indexHtml); err != nil {
		return err
	}
	homeHander := func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
	r.HandleFunc("/", homeHander)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(api.htmlBox)))
	return nil
}