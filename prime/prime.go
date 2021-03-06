package prime

import (
	"appengine"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/prime", primeHandler)
	http.HandleFunc("/nextprime", nextprimeHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	type PageInfo struct {
		Title, Author, BaseUrl, BlogUrl, RepoUrl string
	}
	path := r.URL.Path[1:]
	title := "Prime Number Web App"
	var page, scheme string

	// Figure out where our templates are located.
	templateDir := "templates"

	switch path {
	default:
		page = filepath.Join(templateDir, "index.html")
	case "usage":
		page = filepath.Join(templateDir, "usage.html")
		title += " - Usage"
	}

	// For calculating the base URL
	if nil == r.TLS {
		scheme = "http"
	} else {
		scheme = "https"
	}

	t, err := template.ParseFiles(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageInfo := PageInfo{
		Title:   title,
		Author:  "Erik L. Arneson",
		BaseUrl: scheme + "://" + r.Host + "/",
		RepoUrl: "https://github.com/pymander/prime-json-service",
		BlogUrl: "http://arnesonium.com/2015/02/prime-number-google-app-engine/",
	}

	t.Execute(w, pageInfo)
}

func primeHandler(w http.ResponseWriter, r *http.Request) {
	var numberstring string
	var output []byte
	var result *Result
	var err error
	c := appengine.NewContext(r)

	numberstring = r.FormValue("number")
	if len(numberstring) > 300 {
		http.Error(w, "Please only test integers of less than 300 digits.", http.StatusNotAcceptable)
		return
	}

	// Check primality.
	result, err = IsPrime(c, numberstring)

	output, err = json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(output))
}

func nextprimeHandler(w http.ResponseWriter, r *http.Request) {
	var output []byte
	c := appengine.NewContext(r)

	result, err := GetNextPrime(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err = json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(output))
}
