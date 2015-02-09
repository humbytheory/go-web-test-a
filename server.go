package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Page struct {
	Title string
}

var templates = template.Must(template.ParseFiles("static/header.html"))

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	templates.ExecuteTemplate(w, "header.html", Page{Title: "Home"})
	log.Println("get")
	for i, p := range r.Form {
		log.Printf("i:%-20v  p:%v", i, p)
	}

}

func PostList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	log.Println("post")
	for i, p := range r.Form {
		log.Printf("i:%-20v  p:%v", i, p)
	}
	http.Redirect(w, r, "/", http.StatusFound)

}

func main() {
	var hostport = flag.String("hostport", "localhost:8000", "host:port to server on")
	var staticPath = flag.String("staticPath", "static/", "Path to static files")

	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/", RootHandler).Methods("GET")
	router.HandleFunc("/", PostList).Methods("POST")

	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(
				http.Dir(*staticPath),
			),
		),
	)

	addr := fmt.Sprintf("%s", *hostport)
	log.Printf("listing on http://%s/", addr)

	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}

}
