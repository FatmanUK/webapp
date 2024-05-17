package main

import (
	"net/http"
	"text/template"
	"regexp"
	"bytes"
	"fmt"
	"bufio"
	"strings"
)

// https://go.dev/doc/articles/wiki/

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))
var statics = template.Must(template.ParseFiles("templates/header.html", "templates/footer.html"))
var staticHead bytes.Buffer
var staticFoot bytes.Buffer

const titleRe = "[a-zA-Z0-9]+"
var validPath = regexp.MustCompile("^/(edit|save|view)/(" + titleRe + ")$")
var linkRe = regexp.MustCompile("\\[(" + titleRe + ")\\]")
func initStaticTemplates() error {
	//fmt.Println("Loading static templates.")
	err := statics.ExecuteTemplate(&staticHead, "header.html", nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	//fmt.Println("HEAD: ", staticHead.String())
	err = statics.ExecuteTemplate(&staticFoot, "footer.html", nil)
	return err
}

func markupOutput(o []byte) ([]byte, error) {
	var err error
	var os bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(string(o)))
	for scanner.Scan() {
		err = scanner.Err()
		if err != nil {
			return []byte(""), err
		}
		line := []byte(scanner.Text())

		os.Write(line)
		os.Write([]byte("\n"))
	}

	// add a header and footer which are not marked up
	// so we can have actual HTML tags in there
	var capped bytes.Buffer
	capped.Write(staticHead.Bytes())
	capped.Write(smoothOs)
	capped.Write(staticFoot.Bytes())
	return capped.Bytes(), err
}

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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(output)
}
