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

func (re *Page) Debug() string {
	output := `
## Pages
___`
	return output
}

var db *gorm.DB

func (re *Page) OpenDatabase(dbfile string) {
	d, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	db = d
	db.AutoMigrate(&Page{})
}

func (re *Page) LoadPage(title string) error {
	result := db.Where("title = ?", title).Last(re)
	return result.Error
}

func (re *Page) Save() {
	db.Create(re)
}
