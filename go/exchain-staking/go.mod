module github.com/okex/solidity-sample/go/exchain-staking

go 1.16

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/btcsuite/btcd v0.21.0-beta // indirect
	github.com/cespare/cp v1.1.1 // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/ethereum/go-ethereum v1.10.8
	github.com/go-kit/kit v0.10.0 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/okex/exchain-ethereum-compatible v1.0.2
	github.com/prometheus/tsdb v0.9.1 // indirect
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/okex/cosmos-sdk v0.39.2-exchain16
	github.com/tendermint/iavl => github.com/okex/iavl v0.14.3-exchain1
	github.com/tendermint/tendermint => github.com/okex/tendermint v0.33.9-exchain11
	github.com/tendermint/tm-db => github.com/okex/tm-db v0.5.2-exchain1
)
