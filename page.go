package main

import (
	"os"
)

type Page struct {
	Title   string
	Body    []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(webRoot + "/data/" + filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(webRoot + "/data/" + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
