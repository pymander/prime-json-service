package hello

import (
	"fmt"
	"net/http"
	"html/template"
	"math/big"
	"encoding/json"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/prime", prime)
}

func handler(w http.ResponseWriter, r *http.Request) {
	type PageInfo struct {
		Title, Author, BlogUrl, RepoUrl string
	}

	t, _ := template.ParseFiles("templates/index.html")
	page := PageInfo{
		Title: "Prime Number Testing App",
		Author: "Erik L. Arneson",
	}

	t.Execute(w, page)
}

// Prime number stuff
type Result struct {
	Count int
	Number big.Int
	Prime bool
	Happy bool
}

// @TODO Happy number test
// @TODO Store primes in data store - including count
// @TODO Lookup primes from data store

func prime(w http.ResponseWriter, r *http.Request) {
	var numberstring string
	result := new(Result)

	// Obvious prime testing things
	// 1. It ends in an even number.
	// 2. It ends in a 5.


	numberstring = r.FormValue("number")

	result.Number.SetString(numberstring, 10)
	result.Prime = result.Number.ProbablyPrime(50)

	output,err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(output))
}

