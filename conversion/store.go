package conversion

import (
	"fmt"
	sync "sync"
	"time"
)

// MemoryGraph this is a graph store that effectively maps pairs of currencies for faster lookup
type MemoryGraph struct {
	nodes []*Currency
	edges map[Currency][]*Rate
	lock  sync.RWMutex
}

func NewMemoryGraph() *MemoryGraph {
	return &MemoryGraph{
		lock: sync.RWMutex{},
	}
}

// addNode adds a node to the graph
func (mem *MemoryGraph) addNode(n Currency) {
	mem.lock.Lock()
	if !mem.hasNode(n) {
		mem.nodes = append(mem.nodes, &n)
	}
	mem.lock.Unlock()
}

// addEdge adds an edge to the graph, connects currency pair by their edges,
// one currency can pair with many other currencies
// and vice-versa
func (mem *MemoryGraph) addEdge(provider Provider, n1, n2 Currency) (map[Currency][]*Rate, error) {
	mem.lock.Lock()
	if mem.edges == nil {
		mem.edges = make(map[Currency][]*Rate)
	}

	now := time.Now()

	var edge1 *Rate
	var timestamp time.Time
	var err error
	// if edge already exist,
	ok, rate := mem.hasEdge(n1, n2)
	fmt.Println(ok)
	if ok {
		// check if it's due for another conversion call
		elapsed := rate.timestamp.Add(time.Duration(now.Second())).Second()
		if elapsed < provider.NextTick() {
			mem.lock.Unlock()
			return mem.edges, nil
		}
		// update edge
		newRate, _ := provider.Convert(n1, n2)
		rate = newRate
	} else {
		edge1, err = provider.Convert(n1, n2)
		if err != nil {
			mem.lock.Unlock()
			fmt.Println(err) // just log the error
			return nil, err
		}
		timestamp = time.Now()
		edge1.timestamp = &timestamp
		// store the n1->n2 currency conversion as a connection
		mem.edges[n1] = append(mem.edges[n1], edge1)
	}

	edge2 := &Rate{
		timestamp: edge1.timestamp,
	}
	edge2.tgt_currency = n1
	edge2.src_currency = n2
	nval := float64(1) / (*edge1.value)
	edge2.value = &nval // inverse conversion

	// store the n1<-n2 reverse currency conversion as a connection
	mem.edges[n2] = append(mem.edges[n2], edge2)
	mem.lock.Unlock()

	return mem.edges, nil
}

func (mem *MemoryGraph) hasNode(c Currency) bool {
	var result bool = false
	for _, n := range mem.nodes {
		if n.uid == c.uid {
			result = true
			break
		}
	}

	return result
}

func (mem *MemoryGraph) HasNode(c Currency) bool {
	return mem.hasNode(c)
}

func (mem *MemoryGraph) hasEdge(c1, c2 Currency) (bool, *Rate) {
	var result bool = false
	var rate *Rate = nil
	for _, e := range mem.edges[c1] {
		if e.tgt_currency.uid == c2.uid {
			result = true
			rate = e
			break
		}
	}

	return result, rate
}

func (mem *MemoryGraph) SetPair(provider Provider, p *Pair) (map[Currency][]*Rate, error) {
	mem.addNode(p.src_currency)
	mem.addNode(p.tgt_currency)

	return mem.addEdge(provider, p.src_currency, p.tgt_currency)
}
