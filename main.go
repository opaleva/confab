package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	pattern  *template.Template
}

func (t *templateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	t.once.Do(func() {
		t.pattern = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	err := t.pattern.Execute(writer, request)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var url = flag.String("url", ":8080", "App URL")
	flag.Parse()
	nc := newChat()
	http.Handle("/", &templateHandler{filename: "base.html"})
	http.Handle("/chat", nc)
	go nc.start()

	log.Println("Starting web server on", *url)
	log.Fatal(http.ListenAndServe(*url, nil))
}
