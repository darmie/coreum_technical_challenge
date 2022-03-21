package conversion

import (
	"fmt"
	"net/http"
	sync "sync"
	"time"

	coingecko "github.com/superoo7/go-gecko/v3"
	"github.com/thoas/go-funk"
)

var (
	btc = Currency{
		name: "Bitcoin",
		uid:  "btc",
	}

	eth = Currency{
		name: "Ether",
		uid:  "eth",
	}

	astro = Currency{
		name: "Astro",
		uid:  "astro",
	}

	astronaut = Currency{
		name: "Astronaut",
		uid:  "astronaut",
	}

	currencyPairs = []*Pair{
		{astro, eth},
		{astro, btc},
		{astronaut, eth},
		{astronaut, btc},
	}
)

type ProviderImpl struct {
	store *MemoryGraph
	proc  sync.WaitGroup
}

func NewProviderImpl() *ProviderImpl {
	return &ProviderImpl{store: new(MemoryGraph), proc: sync.WaitGroup{}}
}

// Convert currency pair
func (prov ProviderImpl) Convert(c1, c2 Currency) (*Rate, error) {

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	CG := coingecko.NewClient(httpClient)
	res, err := CG.SimplePrice([]string{c1.uid.(string)}, []string{c2.uid.(string)})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ret := &Rate{}

	ret.tgt_currency = c2
	ret.src_currency = c1

	if res != nil {
		result := *res
		uid1 := c1.uid.(string)
		uid2 := c2.uid.(string)
		value := result[uid1]
		if funk.Contains(value, uid2) {
			val := float64(value[uid2])
			ret.value = &val
			stamp := time.Now()
			ret.timestamp = &stamp
		}
	}

	return ret, nil
}

// NextTick get the time duration before the next conversion call
func (prov ProviderImpl) NextTick() int {
	return 30
}

func (prov ProviderImpl) setPairs() map[Currency][]*Rate {
	mem := prov.store

	for _, pair := range currencyPairs {
		// Queue this conversion process to the wait group
		prov.proc.Add(1)
		// Run this process concurrently
		go func(pair *Pair) {

			defer prov.proc.Done() // Release this process from queue

			_, err := mem.SetPair(&prov, pair) // calls prov.Convert() internally
			if err != nil {
				prov.proc.Done()
			}

		}(pair)
	}

	// wait for concurrently running processes to finish
	prov.proc.Wait()
	return mem.edges
}

func (prov ProviderImpl) LiveRates(page *int, limit *int) []*Rate {
	_page := 1
	_limit := 10
	if page != nil {
		_page = *page
	}

	if limit != nil {
		_limit = *limit
	}

	rates := prov.setPairs()
	var flat_rates []*Rate

	// flatten all rates
	for _, c_rates := range rates {
		flat_rates = append(flat_rates, c_rates...)
	}

	size := len(flat_rates)
	offset := _page - 1
	length := _limit

	// adjust the length for the last page
	if size < (length + offset) {
		length = size - offset
		fmt.Println(length)
	}

	return flat_rates[offset:length]
}

func (prov ProviderImpl) GetRates(c Currency) []*Rate {
	rates := prov.setPairs()
	fmt.Println(rates)
	return rates[c]
}
