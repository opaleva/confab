package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type TemplateHandler struct {
	once     sync.Once
	filename string
	temp     *template.Template
}

func (t *TemplateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	t.once.Do(func() {
		t.temp = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	err := t.temp.Execute(writer, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.Handle("/", &TemplateHandler{filename: "base.html"})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
