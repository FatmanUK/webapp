package main

import (
	"net/http"
	"text/template"
	"regexp"
)

// https://go.dev/doc/articles/wiki/

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))

const titleRe = "[a-zA-Z0-9]+"
var validPath = regexp.MustCompile("^/(edit|save|view)/(" + titleRe + ")$")
var linkRe = regexp.MustCompile("\\[(" + titleRe + ")\\]")

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
