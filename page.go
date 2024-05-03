package main

import (
//	"database/sql"
	"os"
//	"fmt"
//	"text/template"
)

type Page struct {
	Title   string
	Body    []byte
	Body_P  string
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile("data/" + filename, p.Body, 0600)
}

func makeLinks(m []byte) []byte {
	match := string(m)
	match = match[1:len(match) - 1]
	//fmt.Println(match)
	return []byte("<a href=\"/view/" + match + "\">" + match + "</a>")
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile("data/" + filename)
	if err != nil {
		return nil, err
	}
	body_p := string(linkRe.ReplaceAllFunc(body, makeLinks))
	return &Page{Title: title, Body: body, Body_P: body_p}, nil
}
