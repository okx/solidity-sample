package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

var PrivKey666 = "01375084a83cad214b3c53210efa7f4e67cfc409af3245bf8886ff446bf36328"

func main() {
	client, err := ethclient.Dial(RpcUrl)
	if err != nil {
		log.Fatalf("failed to initialize client: %+v", err)
	}
	send(client, "0xa6F064a245B8EF99dB02690b649690f3FB731CeB")

}

func send(client *ethclient.Client, to string) {
	privateKey, senderAddress := initKey(PrivKey666)
	toAddress := common.HexToAddress(to)

	// send 0.001okt
	TransferOKT(client, senderAddress, toAddress, big.NewInt(1*1e15), privateKey, 0)
}
