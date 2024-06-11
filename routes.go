package main

import (
	"net/http"
	"regexp"
)

const titleRe = "[a-zA-Z0-9]+"
var validPath = regexp.MustCompile("^/(edit|save|view)/(" + titleRe + ")$")

func createRoutes() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	fs := http.FileServer(http.Dir("static"))
	t := http.StripPrefix("/static/", fs)
	http.Handle("/static/", t)

	http.HandleFunc("/", handler)
}

// https://go.dev/doc/articles/wiki/

func run() error {
	return http.ListenAndServeTLS(":8443", c.GetString("tls.crt"), c.GetString("tls.key"), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/" + c.GetString("web.first_page"), http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
