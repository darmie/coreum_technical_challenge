package conversion

import "time"

type Currency struct {
	name string
	uid  interface{} // could be an int or string
}

type Pair struct {
	src_currency Currency
	tgt_currency Currency
}

type Rate struct {
	Pair
	value     *float64
	timestamp *time.Time
}

type Provider interface {
	// Covert currency pair
	Convert(c1, c2 Currency) (*Rate, error)
	// NextTick get the time duration  before the next conversion call
	NextTick() int
	// LiveRates fetch all rates of currencies
	LiveRates(page *int, limit *int) []*Rate
	// GetRates of specific currency
	GetRates(c Currency) []*Rate
}
