package primejson

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

// Adapted from http://rosettacode.org/wiki/Happy_numbers#Go
// This is a good example of me being uncomfortable with pointers in Go yet.
func happy(arg big.Int) bool {
	var zero = big.NewInt(0)
	var one  = big.NewInt(1)
	var ten  = big.NewInt(10)
	var n big.Int

	n.Set(&arg)
	m := make(map[string]bool)
	for n.Cmp(one) > 0 {
		m[n.String()] = true
		var d, x big.Int
		for x, n = n, *zero; x.Cmp(zero) > 0; x.Div(&x, ten) {
			d.Mod(&x, ten)
			n.Add(&n, d.Mul(&d, &d))
		}
		if m[n.String()] {
			return false
		}
	}
	return true
}

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
	result.Happy = happy(result.Number)

	output,err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(output))
}

