package main

import (
	"net/http"
	"text/template"
	"bytes"
)

// https://go.dev/doc/articles/wiki/

var templates = template.Must(template.ParseFiles(
	"templates/edit.html",
	"templates/view.html",
	"templates/header.html",
	"templates/footer.html"))

func captureTemplate(tmpl string, p *Page) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := templates.ExecuteTemplate(buf, tmpl + ".html", p)
	return buf.Bytes(), err
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	output, err := captureTemplate(tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if tmpl == "view" {
		output, err = markupOutput(output)
	}
	// add a header and footer which are not marked up
	// so we can have actual HTML tags in there
	var capped bytes.Buffer
	sh, err := captureTemplate("header", p)
	capped.Write(sh)
	capped.Write(output)
	sf, err := captureTemplate("footer", p)
	capped.Write(sf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(capped.Bytes())
}
