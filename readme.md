# Coreum Technical Challenge

## A. Currency Conversion
### (1)
Prepare .proto file with grpc service & message type definitions (no implementation
required). Describe service methods to perform conversion, batch conversion, listing of rates
with pagination.

#### Solution
see [generated grpc code](conversion/conversion_grpc.pb.go)

```proto
syntax = "proto3";
package conversion;

option go_package = "/conversion";

message ConversionRequest {
    repeated string sources = 1;
    repeated string targets = 2;

    double value = 3; // The value to be converted, proto3 doesn't support defauly values; client may be forced to provide one.
}

// This is a general result structure, that can be used for single or many conversions. 
message ConversionResult {
    message ConversionData {
        // source currency
        string src_currency = 1;
        
        // A map of target currency and the conversion value
        message DataMap {
            string tgt_currency = 1;
            double rate = 2;
        }

        // supports many mapped rates in case we want for single to single, single to many or many to many conversions
        repeated DataMap rate = 2;
    }

    repeated ConversionData data = 2;
}

message BatchConversionRequest {
    ConversionRequest conversion = 1;
    int32 batch_size = 2;
}



message RatesRequest {
    ConversionRequest conversion = 1;
    int32 limit = 2;
    int32 page = 3;
   
}

// This represents result data that will hold the paginated list of rates 
message RatesResult {
    ConversionResult result = 1;
    int32 curr_page = 2;
    bool has_next = 3; 
}

service ConversionService {
    // convert currencies
    rpc Convert(ConversionRequest) returns (ConversionResult){}
    // perform batch conversion, client sends a (stream) sequence of batched conversion requests
    rpc ConvertBatch(stream BatchConversionRequest) returns (ConversionResult){}
    // get list of rates, the ones already persisted on db/storage
    rpc Rates(RatesRequest) returns (RatesResult){}
}
```
### (2) 
To perform conversion fast we want to store live exchange rates in local memory and sync
them with some external source via API e.g (Coingecko API or CurrencyLayer). Implement
part of exchange rates sync between external API and local memory (concurrency should be
handled). Also define an interface which provider should implement to fetch live rates
suitable for different providers (Coingecko/CurrencyLayer etc)

see implementations [here](conversion/store.go) & [here](conversion/provider_impl.go)

#### Problem (i) 
Implement part of exchange rates sync between external API and local memory (concurrency should be
handled)
#### Problem (ii)
Define an interface which provider should implement to fetch live rates
suitable for different providers (Coingecko/CurrencyLayer etc)

### Solution (i)
#### _Local Memory_
Implementing local memory to store converted currency pairs could have been done via a one-to-one `map`, which leads to new problem like; 
> How do we map a single source/primary currency to mutiple currency rates? 

I used a map of arrays `map[Currency][]*Rate`. This also helps us also to perform two-way conversion, e.g `BTC -> ETH`  and `ETH -> BTC`, then I can quickly look up `ETH -> BTC` without another conversion call. We can also store conversions like this `USD -> [BTC, ETH, BNB, SOL]` and have `BTC -> USD`, `ETH -> USD`, `SOL -> USD` already stored in memory.

> How do we prevent overwriting rates that are already stored to a source/primary currency? 

I append new rates to the location of the currency if it already exists.

> How do we avoid calling conversion API too many times ?

I check if a rate pair already exist in memory. Memory is implemented as a graph to facilitate faster look up. Rate pairs are stored as edges, to a currency node. I look up the edge which is stored as an array and mapped to a source currency as a key. E.G `edges[BTC] -> [ETH, BNB, SOL]`, then loop through `[ETH, BNB, SOL]` for the right-hand pair. 

If the pair exist in memory, we can then check if the timestamp from the last conversion call as elapsed the provided duration before proceeding to call the API. 

#### _Concurrency_
To synchronize data fetched from API concurrently as I write the latest data to the memory, I implememnted a Mutex lock while performing API call and writing data to memory, then unlock it afterwards.  While data is being fetched concurrently, it is being written to memory one at a time.

### Solution (ii)
see [provider_impl.go](conversion/provider_impl.go) for an example implementation of the interface
```go
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
```


## Testing
Run the following in terminal `go install ./conversion && go test ./conversion`






