package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	// "unicode"

	"github.com/astaxie/beego/validation"
	"github.com/gorilla/mux"
)

type Page struct {
	Title string
}

var data = map[string]string{
	"BaseUrl": "zzz",
	"Options": "---opt---",
	"Version": "123",
	"Title":   "yo",
}

var tmpl = make(map[string]*template.Template)

func init() {
	// nas nfs         -- add delete modify
	// nas smb share   -- add delete modify
	// nas smb dir     -- add delete modify
	// nas smb ad user -- add
	tmpl["main"] = template.Must(template.ParseFiles("views/footer.html", "views/navbar.html", "views/index.html", "views/head.html", "views/base.html"))
	tmpl["san"] = template.Must(template.ParseFiles("views/footer.html", "views/navbar.html", "views/index.html", "views/head.html", "views/base.html"))
	tmpl["nas"] = template.Must(template.ParseFiles("views/nas/smb-start.html", "views/footer.html", "views/navbar.html", "views/head.html", "views/base.html"))
	tmpl["nas-new"] = template.Must(template.ParseFiles("views/nas/nas-smb-form-new.html", "views/nas/smb-share-add.html", "views/footer.html", "views/navbar.html", "views/head.html", "views/base.html"))
	tmpl["backups"] = template.Must(template.ParseFiles("views/footer.html", "views/navbar.html", "views/index.html", "views/head.html", "views/base.html"))

}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	data["NavActive"] = "main"
	log.Printf("get for %v from %s", data["NavActive"], r.RemoteAddr)

	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), http.StatusInternalServerError)
	}
	tmpl["main"].ExecuteTemplate(w, "base", data)
}

func SanHandler(w http.ResponseWriter, r *http.Request) {
	data["NavActive"] = "san"
	log.Printf("get for %v from %s", data["NavActive"], r.RemoteAddr)

	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), http.StatusInternalServerError)
	}
	tmpl["index"].ExecuteTemplate(w, "base", data)
}

func NasHandler(w http.ResponseWriter, r *http.Request) {
	data["NavActive"] = "nas"
	log.Printf("get for %v from %s", data["NavActive"], r.RemoteAddr)

	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), http.StatusInternalServerError)
	}
	tmpl["nas"].ExecuteTemplate(w, "base", data)
}

func NasNewHandler(w http.ResponseWriter, r *http.Request) {
	data["NavActive"] = "nas"
	data["FormType"] = "nas-smb-form-new"
	log.Printf("get for %v from %s", data["NavActive"], r.RemoteAddr)
	w.Header().Set("Content-type", "text/html")

	tmpl["nas-new"].ExecuteTemplate(w, "base", data)
}

func BackupsHandler(w http.ResponseWriter, r *http.Request) {
	data["NavActive"] = "backups"
	log.Printf("get for %v from %s", data["NavActive"], r.RemoteAddr)

	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), http.StatusInternalServerError)
	}
	tmpl["index"].ExecuteTemplate(w, "base", data)
}

func RtSearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	log.Printf("get from %s", r.RemoteAddr)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), http.StatusInternalServerError)
	}
	tmpl["index"].ExecuteTemplate(w, "base", data)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}

	log.Printf("post for %v from %s  on path: %v", "PostHandler", r.RemoteAddr, r.URL.Path)
	type filter struct {
		Name string
		Type string
		Max  int
	}

	valid := validation.Validation{}

	valid.MaxSize(r.PostFormValue("storageWorm"), 1, "Worm")
	valid.MaxSize(r.PostFormValue("storageBackup"), 2, "Backups")
	valid.MaxSize(r.PostFormValue("storageSize"), 5000, "Size")
	valid.MaxSize(r.PostFormValue("storageHidden"), 1, "Hidden")

	valid.AlphaNumeric(r.PostFormValue("storageOwner"), "Owner")
	valid.AlphaDash(r.PostFormValue("storageName"), "Requested Name")
	valid.Alpha(r.PostFormValue("storageClass"), "Classification")
	valid.Alpha(r.PostFormValue("storageAccessList"), "Access List")

	if r.PostFormValue("storageAlertGroup") != "" {
		valid.Email(r.PostFormValue("storageAlertGroup"), "Email Alert")
	}

	if valid.HasErrors() {
		errormap := []string{}
		for _, err := range valid.Errors {
			errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
		}
		for _, e := range errormap {
			log.Printf("%v", e)
		}
	}

	// log.Printf("form: %v", r.Form)
	http.Redirect(w, r, "/nas/new/", http.StatusFound)

}

func main() {
	var staticPath = flag.String("staticPath", "static/", "Path to static files")

	port := os.Getenv("PORT")

	flag.Parse()

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/", RootHandler).Methods("GET")
	router.HandleFunc("/san/", SanHandler).Methods("GET")
	router.HandleFunc("/nas/", NasHandler).Methods("GET")
	router.HandleFunc("/nas/new/", NasNewHandler).Methods("GET")
	router.HandleFunc("/nas/new/", PostHandler).Methods("POST")
	router.HandleFunc("/backups/", BackupsHandler).Methods("GET")
	router.HandleFunc("/rtsearch/", RtSearchHandler).Methods("POST")

	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(
				http.Dir(*staticPath),
			),
		),
	)

	addr := fmt.Sprintf("localhost:%s", port)
	log.Printf("listening on http://%s/", addr)

	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}

}
