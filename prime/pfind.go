package prime

// Search for prime numbers, write them to databases, that sort of thing.

import (
	"appengine"
	"math/big"
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
