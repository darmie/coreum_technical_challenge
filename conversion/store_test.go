package conversion

import (
	"testing"
)

func TestMemoryGraph_addNode(t *testing.T) {
	c := Currency{"BTC", "btc"}
	type args struct {
		n Currency
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Test add node",
			args{c},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewMemoryGraph()
			mem.addNode(tt.args.n)

			if len(mem.nodes) != 1 {
				t.Error("There should be at least 1 memory graph node")
				return
			}

			if mem.nodes[0].uid != c.uid {
				t.Errorf("Wants %s, got %s", c.uid, mem.nodes[0].uid)
				return
			}
		})

	}
}

func TestMemoryGraph_addEdge(t *testing.T) {
	c1 := Currency{"ASTRONAUT", "astronaut"}
	c2 := Currency{"ETH", "eth"}
	type args struct {
		provider Provider
		n1       Currency
		n2       Currency
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Test add edge",
			args{
				NewProviderImpl(),
				c1,
				c2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewMemoryGraph()
			mem.addEdge(tt.args.provider, tt.args.n1, tt.args.n2)

			if len(mem.edges) < 2 {
				t.Error("There should be at least 2 memory graph edges")
				return
			}

			if val, ok := mem.edges[c1]; ok {
				if len(val) == 0 {
					t.Errorf("There should be at least 1 currency pair in edges[%s]", c1.uid)
					return
				}
			} else {
				t.Errorf("Memory graph edge should have a key of currency %v", c1)
				return
			}

			if val, ok := mem.edges[c2]; ok {
				if len(val) == 0 {
					t.Errorf("There should be at least 1 currency pair in edges[%s]", c2.uid)
					return
				}
			} else {
				t.Errorf("Memory graph edge should have a key of currency %v", c1)
				return
			}

		})
	}
}

func TestMemoryGraph_SetPair(t *testing.T) {
	c1 := Currency{"ASTRONAUT", "astronaut"}
	c2 := Currency{"ETH", "eth"}
	type args struct {
		provider Provider
		p        *Pair
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Test set pair",
			args{
				NewProviderImpl(),
				&Pair{
					c1,
					c2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewMemoryGraph()
			mem.SetPair(tt.args.provider, tt.args.p)

			if len(mem.nodes) < 1 {
				t.Error("There should be at least 2 memory graph nodes")
				return
			}

			if len(mem.edges) < 2 {
				t.Error("There should be at least 2 memory graph edges")
				return
			}

			if val, ok := mem.edges[c1]; ok {
				if len(val) == 0 {
					t.Errorf("There should be at least 1 currency pair in edges[%s]", c1.uid)
					return
				}
			} else {
				t.Errorf("Memory graph edge should have a key of currency %v", c1)
				return
			}

			if val, ok := mem.edges[c2]; ok {
				if len(val) == 0 {
					t.Errorf("There should be at least 1 currency pair in edges[%s]", c2.uid)
					return
				}
			} else {
				t.Errorf("Memory graph edge should have a key of currency %v", c1)
				return
			}
		})
	}
}
