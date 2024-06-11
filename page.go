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
	return os.WriteFile(c.GetString("web.root") + "/" + filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(c.GetString("web.root") + "/" + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
