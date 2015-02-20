package prime

// Prime number storage.
import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"time"
)

// Because the datastore can't store big.Int types, we need to marshall.
type PrimeRecord struct {
	Data []byte
}

func LookupPrime(c appengine.Context, num string) (*Result, error) {
	key := datastore.NewKey(c, "Prime", num, 0, nil)
	result := new(PrimeRecord)
	if err := datastore.Get(c, key, result); err != nil {
		return nil, err
	}

	primeResult := new(Result)
	json.Unmarshal(result.Data, &primeResult)

	return primeResult, nil
}

func StorePrime(c appengine.Context, prime *Result) error {
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

type LastPrime struct {
	Number      string
	RequestTime time.Time
}

func LookupLastPrime(c appengine.Context) (*LastPrime, error) {
	key := datastore.NewKey(c, "Prime", "LastPrime", 0, nil)
	result := new(LastPrime)
	if err := datastore.Get(c, key, result); err != nil {
		return nil, err
	}

	return result, nil
}

func StoreLastPrime(c appengine.Context, record *LastPrime) error {
	key := datastore.NewKey(c, "Prime", "LastPrime", 0, nil)
	if _, err := datastore.Put(c, key, record); err != nil {
		return err
	}

	return nil
}
