package main

import (
	bytes2 "bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/app"
	"github.com/okex/exchain/app/codec"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/bytes"
)

var (
	host    = "https://exchaintestrpc.okex.org"
	privKey = "89c81c304704e9890025a5a91898802294658d6e4034a11c6116f4b129ea12d3"
	oecTestnetChainId int64 = 65
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
	client, err := ethclient.Dial(host)
	if err != nil {
		log.Fatalf("failed to initialize client: %+v", err)
	}
	// 0.2 get the chain-id from network
	chainID := big.NewInt(oecTestnetChainId)
	if err != nil {
		log.Fatalf("failed to fetch the chain-id from network: %+v", err)
	}
	// 0.3 unencrypted private key -> secp256k1 private key
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		log.Fatalf("failed to switch unencrypted private key -> secp256k1 private key: %+v", err)
	}
	// 0.4 secp256k1 private key -> pubkey -> address
	pubkey := privateKey.Public()
	pubkeyECDSA, ok := pubkey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalln("failed to switch secp256k1 private key -> pubkey")
	}
	fromAddress := crypto.PubkeyToAddress(*pubkeyECDSA)

	// 0.5 get the gasPrice
	gasPrice := big.NewInt(10000000000)
	//
	// 1. deploy contract
	//
	//contractAddr := deployContract(client, fromAddress, gasPrice, chainID, privateKey)
	contractAddr := common.HexToAddress("0x79BE5cc37B7e17594028BbF5d43875FDbed417db")
	//
	// 2. call contract(write)
	//
	readContract(client, contractAddr)

	writeContract(client, fromAddress, gasPrice, chainID, privateKey, contractAddr)
	time.Sleep(time.Second * 5)
	//
	// 3. call contract(read)
	//
	readContract(client, contractAddr)
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
	hash := getTxHash(signedTx)
	time.Sleep(time.Second * 5)

	receipt, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		log.Fatal(err)
	}

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

func writeContract(client *ethclient.Client,
	fromAddress common.Address,
	gasPrice *big.Int,
	chainID *big.Int,
	privateKey *ecdsa.PrivateKey,
	contractAddr common.Address) {
	// 0. get the value of nonce, based on address
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("failed to fetch the value of nonce from network: %+v", err)
	}

	fmt.Printf(
		"writeContract: \n" +
			"	sender Address<0x%s>, \n" +
		"	gasPrice<%s>, \n" +
		"	contractAddr<%s>\n",
		hex.EncodeToString(fromAddress.Bytes()),
		gasPrice.String(),
		contractAddr.String())

	unsignedTx := writeContractTx(nonce, contractAddr, gasPrice)
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
}

func writeContractTx(nonce uint64, contractAddr common.Address, gasPrice *big.Int) *types.Transaction {
	value := big.NewInt(0)
	gasLimit := uint64(3000000)

	num := big.NewInt(100)
	data, err := sampleContractABI.Pack("add", num)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("	abi <add>: %d\n", num)
	fmt.Printf("	nonce: %d\n", nonce)

	return types.NewTransaction(nonce, contractAddr, value, gasLimit, gasPrice, data)
}

func readContract(client *ethclient.Client, contractAddr common.Address) {
	data, err := sampleContractABI.Pack("getCounter")
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

	ret, err := sampleContractABI.Unpack("getCounter", output)
	if err != nil {
		panic(err)
	}
	fmt.Printf("readContract <getCounter>: %d\n", ret)
}

func getTxHash(signedTx *types.Transaction) common.Hash {
	//ts := types.Transactions{signedTx}
	//rawTx := hex.EncodeToString(ts.)

	rawTxBytes, err := hex.DecodeString("rawTx")
	if err != nil {
		log.Fatal(err)
	}

	tx := new(evmtypes.MsgEthereumTx)
	// RLP decode raw transaction bytes
	if err := rlp.DecodeBytes(rawTxBytes, tx); err != nil {
		log.Fatal(err)
	}

	cdc := codec.MakeCodec(app.ModuleBasics)
	txEncoder := authclient.GetTxEncoder(cdc)
	txBytes, err := txEncoder(tx)
	if err != nil {
		log.Fatal(err)
	}

	var hexBytes bytes.HexBytes
	hexBytes = tmhash.Sum(txBytes)
	hash := common.HexToHash(hexBytes.String())
	return hash
}