package conversion

import (
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
		{btc, eth},
		{eth, btc},
		{btc, astro},
		{eth, astro},
		{astro, eth},
		{astro, btc},
		{astronaut, eth},
		{astronaut, btc},
		{btc, astronaut},
		{eth, astronaut},
	}
)

type ProviderImpl struct {
	store *MemoryGraph
	proc  sync.WaitGroup
}

func NewProviderImpl() *ProviderImpl {
	return &ProviderImpl{store: new(MemoryGraph)}
}

// Convert currency pair
func (prov ProviderImpl) Convert(c1, c2 Currency) (*Rate, error) {
	// Queue this conversion process to the wait group
	prov.proc.Add(1)

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	CG := coingecko.NewClient(httpClient)
	res, err := CG.SimplePrice([]string{c1.uid.(string)}, []string{c2.uid.(string)})
	if err != nil {
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

	// Release this process from queue
	prov.proc.Done()
	return ret, nil
}

// NextTick get the time duration before the next conversion call
func (prov ProviderImpl) NextTick() int {
	return 30
}

func (prov ProviderImpl) setPairs() map[Currency][]*Rate {
	mem := prov.store
	for _, pair := range currencyPairs {
		// Run this process concurrently, calls prov.Convert() internally
		go mem.SetPair(prov, pair)
	}

	// wait for concurrently running processes to finish
	prov.proc.Wait()
	return mem.edges
}

func (prov ProviderImpl) LiveRates(page *int, limit *int) []*Rate {
	if page != nil {
		*page = 1
	}

	if limit != nil {
		*limit = 10
	}

	rates := prov.setPairs()
	var flat_rates []*Rate

	// flatten all rates
	for _, c_rates := range rates {
		flat_rates = append(flat_rates, c_rates...)
	}

	size := len(flat_rates)
	offset := *page - 1
	length := *limit

	// adjust the length for the last page
	if size < (length + offset) {
		length = size - offset
	}

	return flat_rates[offset:length]
}

func (prov ProviderImpl) GetRates(c Currency) []*Rate {
	rates := prov.setPairs()
	return rates[c]
}
