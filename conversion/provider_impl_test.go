package conversion

import (
	"reflect"
	"testing"
)

func TestProviderImpl_Convert(t *testing.T) {
	c1 := Currency{"ASTRONAUT", "astronaut"}
	c2 := Currency{"ETH", "eth"}
	type args struct {
		c1 Currency
		c2 Currency
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"Test convert two currencies",
			args{c1, c2},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prov := NewProviderImpl()
			got, err := prov.Convert(tt.args.c1, tt.args.c2)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProviderImpl.Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("ProviderImpl.Convert(), want %v, got %v", tt.want, false)
			}

			if got.value == nil {
				t.Errorf("ProviderImpl.Convert(), want %v, got %v", tt.want, false)
			}
		})
	}
}

func TestProviderImpl_setPairs(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			"Tes setPairs",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prov := NewProviderImpl()
			got := prov.setPairs()
			if got == nil {
				t.Errorf("ProviderImpl.setPairs(), want %v, got %v", tt.want, false)
				return
			}

			if len(got) != len(currencyPairs) {
				t.Errorf("want %v, got %v", len(currencyPairs), len(got))
				return
			}

			if len(got[astro]) < 2 {
				t.Errorf("Should have at least 2 conversion pairs with %v", astro)
				return
			}

			if len(got[astronaut]) < 2 {
				t.Errorf("Should have at least 2 conversion pairs with %v", astronaut)
				return
			}
		})
	}
}

func TestProviderImpl_LiveRates(t *testing.T) {
	type args struct {
		page  *int
		limit *int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"Test Liverates",
			args{},
			8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prov := NewProviderImpl()
			if got := prov.LiveRates(tt.args.page, tt.args.limit); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("len(ProviderImpl.LiveRates()) = %v, want %v", len(got), tt.want)
			}

		})
	}
}

func TestProviderImpl_GetRates(t *testing.T) {
	testProvider := NewProviderImpl()
	_ = testProvider.LiveRates(nil, nil)
	type args struct {
		c Currency
	}
	tests := []struct {
		name string
		args args
		want []*Rate
	}{
		{
			"Test GetRates of single currency",
			args{astro},
			testProvider.store.edges[astro],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prov := testProvider
			if got := prov.GetRates(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProviderImpl.GetRates() = %v, want %v", got, tt.want)
			}
		})
	}
}
