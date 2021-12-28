package main

import (
	bytes2 "bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/okex/exchain-ethereum-compatible/utils"
)

const (
	RpcUrl          = "https://exchaintestrpc.okex.org"
	ChainId  int64  = 65 //  oec testnet
	PrivKey         = "89c81c304704e9890025a5a91898802294658d6e4034a11c6116f4b129ea12d3"
	GasPrice int64  = 100000000 // 0.1 gwei
	GasLimit uint64 = 3000000
)

var (
	sampleContractByteCode []byte
	sampleContractABI      abi.ABI
)

func init() {
	bin, err := ioutil.ReadFile("../../contracts/counter.bin")
	if err != nil {
		log.Fatal(err)
	}
	sampleContractByteCode = common.Hex2Bytes(string(bin))

	abiByte, err := ioutil.ReadFile("../../contracts/counter.abi")
	if err != nil {
		log.Fatal(err)
	}
	sampleContractABI, err = abi.JSON(bytes2.NewReader(abiByte))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//
	// 0. init
	//
	// 0.1 init client
	client, err := ethclient.Dial(RpcUrl)
	if err != nil {
		log.Fatalf("failed to initialize client: %+v", err)
	}
	// 0.2 get the chain-id from network
	if err != nil {
		log.Fatalf("failed to fetch the chain-id from network: %+v", err)
	}
	// 0.3 unencrypted private key -> secp256k1 private key
	privateKey, err := crypto.HexToECDSA(PrivKey)
	if err != nil {
		log.Fatalf("failed to switch unencrypted private key -> secp256k1 private key: %+v", err)
	}
	// 0.4 secp256k1 private key -> pubkey -> address
	pubkey := privateKey.Public()
	pubkeyECDSA, ok := pubkey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalln("failed to switch secp256k1 private key -> pubkey")
	}
	senderAddress := crypto.PubkeyToAddress(*pubkeyECDSA)

	// 1. deploy contract
	contractAddr := deployContract(client, senderAddress, big.NewInt(GasPrice), big.NewInt(ChainId), privateKey)

	//contractAddr := common.HexToAddress("0x79BE5cc37B7e17594028BbF5d43875FDbed417db")

	// 2. call contract(write)
	ReadContract(client, contractAddr, "getCounter")
	WriteContract(client, contractAddr, senderAddress, privateKey, "add", big.NewInt(100))
	time.Sleep(time.Second * 5)

	ReadContract(client, contractAddr, "getCounter")
}

func WriteContract(client *ethclient.Client,
	contractAddr common.Address,
	fromAddress common.Address,
	privateKey *ecdsa.PrivateKey,
	name string,
	args ...interface{}) {
	// 0. get the value of nonce, based on address
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("failed to fetch the value of nonce from network: %+v", err)
	}

	// 0.5 get the gasPrice
	gasPrice := big.NewInt(GasPrice)

	fmt.Printf(
		"writeContract: \n"+
			"	sender address: <%s>\n"+
			"	contract address: <%s>\n"+
			"	gas price: <%.1f> gwei\n"+
			"	nonce: %d\n"+
			"	ABI: <%s %s>\n",
		fromAddress.Hex(),
		contractAddr.String(),
		float64(gasPrice.Uint64())/1e9,
		nonce, name, args)

	data, err := sampleContractABI.Pack(name, args...)
	if err != nil {
		log.Fatal(err)
	}

	unsignedTx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), GasLimit, gasPrice, data)

	// 2. sign unsignedTx -> rawTx
	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(big.NewInt(ChainId)), privateKey)
	if err != nil {
		log.Fatalf("failed to sign the unsignedTx offline: %+v", err)
	}

	// 3. send rawTx
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadContract(client *ethclient.Client, contractAddr common.Address, name string, args ...interface{}) {
	data, err := sampleContractABI.Pack(name, args...)
	if err != nil {
		log.Fatal(err)
	}

	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	output, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		panic(err)
	}

	ret, err := sampleContractABI.Unpack(name, output)
	if err != nil {
		panic(err)
	}
	fmt.Printf("readContract <%s>: %d\n", name, ret)
}

func deployContract(client *ethclient.Client,
	fromAddress common.Address,
	gasPrice *big.Int,
	chainID *big.Int,
	privateKey *ecdsa.PrivateKey) (contractAddr common.Address) {
	// 0. get the value of nonce, based on address
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("failed to fetch the value of nonce from network: %+v", err)
	}

	//1. simulate unsignedTx as you want, fill out the parameters into a unsignedTx
	unsignedTx := deployContractTx(nonce, gasPrice)

	// 2. sign unsignedTx -> rawTx
	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("failed to sign the unsignedTx offline: %+v", err)
	}

	// 3. send rawTx
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	// 4. get the contract address based on tx hash
	hash, err := utils.Hash(signedTx)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 5)

	receipt, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("new Contract Address: %s\n", receipt.ContractAddress.String())

	return receipt.ContractAddress
}

func deployContractTx(nonce uint64, gasPrice *big.Int) *types.Transaction {
	value := big.NewInt(0)
	gasLimit := uint64(3000000)

	// Constructor
	input, err := sampleContractABI.Pack("")
	if err != nil {
		log.Fatal(err)
	}
	data := append(sampleContractByteCode, input...)
	return types.NewContractCreation(nonce, value, gasLimit, gasPrice, data)
}
