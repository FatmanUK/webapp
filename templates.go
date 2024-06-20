package main

import (
	"net/http"
	"text/template"
	"bytes"
)

type Template struct {
	Name string
}

var errNoLoadTemplates = "Loading templates failed"

var templates *template.Template

func (*Template) Init() {
	glob, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		panic(err.Error())
	}
	templates = template.Must(glob, err)
	if templates == nil {
		panic(errNoLoadTemplates)
	}
}

func captureTemplate(tmpl string, p interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := templates.ExecuteTemplate(buf, tmpl + ".html", p)
	return buf.Bytes(), err
}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	output, err := captureTemplate(tmpl, p)
	herr_500 := http.StatusInternalServerError
	if err != nil {
		http.Error(w, err.Error(), herr_500)
	}
	if tmpl == "view" || tmpl == "debug" {
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
		http.Error(w, err.Error(), herr_500)
	}
	w.Write(capped.Bytes())
}
