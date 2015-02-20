package prime

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/prime", primeHandler)
	//http.HandleFunc("/nextprime", nextprimeHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	type PageInfo struct {
		Title, Author, BaseUrl, BlogUrl, RepoUrl string
	}
	path := r.URL.Path[1:]
	title := "Prime Number Web App"
	var page, scheme string

	switch path {
	default:
		page = "templates/index.html"
	case "usage":
		page = "templates/usage.html"
		title += " - Usage"
	}

	// For calculating the base URL
	if nil == r.TLS {
		scheme = "http"
	} else {
		scheme = "https"
	}

	t, _ := template.ParseFiles(page)
	pageInfo := PageInfo{
		Title:   title,
		Author:  "Erik L. Arneson",
		BaseUrl: scheme + "://" + r.Host + "/",
	}

	t.Execute(w, pageInfo)
}

// Prime number stuff
type Result struct {
	Count  int
	Number *big.Int
	Prime  bool
	Happy  bool
}

// Adapted from http://rosettacode.org/wiki/Happy_numbers#Go
// This is a good example of me being uncomfortable with pointers in Go yet.
func happy(arg *big.Int) bool {
	var zero = big.NewInt(0)
	var one = big.NewInt(1)
	var ten = big.NewInt(10)
	var n big.Int

	n.Set(arg)
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

func primeHandler(w http.ResponseWriter, r *http.Request) {
	var numberstring string
	var number big.Int
	var result *Result
	c := appengine.NewContext(r)

	numberstring = r.FormValue("number")
	if len(numberstring) > 300 {
		http.Error(w, "Please only test integers of less than 300 digits.", http.StatusNotAcceptable)
		return
	}

	number.SetString(numberstring, 10)

	// Obvious prime testing things
	// 1. It ends in an even number.
	// 2. It ends in a 5.

	// Try a lookup first.
	result, _ = lookupPrime(c, number.String())

	// Nothing in the database. Better test it ourselves.
	if nil == result {
		result = &Result{
			Count:  0,
			Number: &number,
			Prime:  number.ProbablyPrime(10),
			Happy:  happy(&number),
		}
	}

	// Keep track of how many times specific prime numbers have been looked for.
	result.Count++

	// If prime, we store.
	if true == result.Prime {
		if err := storePrime(c, result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	output, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(output))
}

// Because the datastore can't store big.Int types, we need to marshall.
type PrimeRecord struct {
	Data []byte
}

func lookupPrime(c appengine.Context, num string) (*Result, error) {
	key := datastore.NewKey(c, "Prime", num, 0, nil)
	result := new(PrimeRecord)
	if err := datastore.Get(c, key, result); err != nil {
		return nil, err
	}

	primeResult := new(Result)
	json.Unmarshal(result.Data, &primeResult)

	return primeResult, nil
}

func storePrime(c appengine.Context, prime *Result) error {
	output, err := json.Marshal(prime)
	if err != nil {
		return err
	}

	record := &PrimeRecord{
		Data: output,
	}

	key := datastore.NewKey(c, "Prime", prime.Number.String(), 0, nil)
	if _, err := datastore.Put(c, key, record); err != nil {
		return err
	}

	return nil
}
