module github.com/okex/solidity-sample/okexchain-web3-sample/golang

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.39.2
	github.com/ethereum/go-ethereum v1.9.25
	github.com/okex/exchain v0.18.4
	github.com/tendermint/tendermint v0.33.9
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/okex/cosmos-sdk v0.39.2-exchain4
    github.com/tendermint/iavl => github.com/okex/iavl v0.14.3-exchain
    github.com/tendermint/tendermint => github.com/okex/tendermint v0.33.9-exchain3
)
