package prime

// Search for prime numbers, write them to databases, that sort of thing.

import (
	"appengine"
	"math/big"
	"time"
)

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

func IsPrime(c appengine.Context, numberstring string) (*Result, error) {
	var number = new(big.Int)

	number.SetString(numberstring, 10)

	return IsPrimeInt(c, number)
}

func IsPrimeInt(c appengine.Context, number *big.Int) (*Result, error) {

	// Obvious prime testing things
	// 1. It ends in an even number.
	// 2. It ends in a 5.

	// Try a lookup first.
	result, _ := LookupPrime(c, number.String())

	// Nothing in the database. Better test it ourselves.
	if nil == result {
		result = &Result{
			Count:  0,
			Number: number,
			Prime:  number.ProbablyPrime(10),
			Happy:  happy(number),
		}
	}

	// Keep track of how many times specific prime numbers have been looked for.
	result.Count++

	// If prime, we store.
	if true == result.Prime {
		if err := StorePrime(c, result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func GetNextPrime(c appengine.Context) (*Result, error) {
	var err error
	result := new(Result)
	lastPrime, err := LookupLastPrime(c)
	two := big.NewInt(2)

	// Have we never lookup up the last prime before? Dang! Start over at 11.
	if nil == lastPrime {
		lastPrime = &LastPrime{
			Number:      "11",
			RequestTime: time.Now(),
		}
	}

	oddNumber := new(big.Int)
	oddNumber.SetString(lastPrime.Number, 10)

	for {
		oddNumber.Add(oddNumber, two)
		result, err = IsPrimeInt(c, oddNumber)
		if err != nil {
			return nil, err
		}

		if true == result.Prime {
			break
		}
	}

	// We have our next prime! Let's save it now.
	lastPrime.Number = oddNumber.String()
	lastPrime.RequestTime = time.Now()
	if err = StoreLastPrime(c, lastPrime); err != nil {
		return nil, err
	}

	return result, nil
}
