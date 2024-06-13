package main

import (
	"gorm.io/gorm"
	"github.com/glebarez/sqlite" // pure Go?
)

type Page struct {
	gorm.Model
	Title   string
	Body    []byte
}

func (re Page) debugOutput() string {
	output := `
## Pages
  
___`
	return output
}

var db *gorm.DB

func openDatabase() {
	dbfile := c.GetString("db.file")
	d, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	db = d
	db.AutoMigrate(&Page{})
}

func loadPage(title string) (*Page, error) {
	var p = &Page{}
	result := db.Where("title = ?", title).Last(&p)
	return p, result.Error
}

func (p *Page) save() error {
	db.Create(p)
	return nil
}
