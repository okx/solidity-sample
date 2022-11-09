package main

import (
	bytes2 "bytes"
	"crypto/ecdsa"
	_ "embed"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/okex/solidity-sample/go/exchain-gasback/contracts"
	"log"
	"math/big"
)

const (
	RpcUrl                 = "https://exchaintestrpc.okex.org"
	ChainId         int64  = 65 //  okc testnet
	PrivKey                = ""
	GasPrice        int64  = 100000000 // 0.1 gwei
	GasLimit        uint64 = 3000000
	contractAddress        = "0x7C9eBC51F32ED4760AEADeF365875ef5727442Bc"
)

const (
	systemContarctAddress  = "0x7e5E6AF6424BE7A835313777Fc6E0d1912e52Fc8"
	gasBackContractAddress = "0x56CeE7c20F4C83996e72b65a7b410fcE4C3b52B0"
)

func main() {

	//load contracts
	contractAbi, err := abi.JSON(bytes2.NewReader([]byte(contracts.SystemContractABI)))
	if err != nil {
		panic(err)
	}
	systemContract := contracts.SystemContract{ByteCode: common.Hex2Bytes(contracts.SystemContractBIN), Abi: contractAbi, Address: common.HexToAddress(systemContarctAddress)}

	contractAbi, err = abi.JSON(bytes2.NewReader([]byte(contracts.GasBackMSGHelperABI)))
	if err != nil {
		panic(err)
	}
	gasBackMSGHelper := &contracts.ContractWrap{ByteCode: common.Hex2Bytes(contracts.GasBackMSGHelperBIN), Abi: contractAbi, Address: common.HexToAddress(gasBackContractAddress), System: &systemContract}

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

	// Register gasback
	withdrawerAddress := senderAddress
	contractAddr := common.HexToAddress(contractAddress)
	nonces := []*big.Int{big.NewInt(394658)}
	res, err := gasBackMSGHelper.Invoke(client, senderAddress, privateKey, GasPrice, GasLimit, ChainId, "genRegisterMsg", &contractAddr, &withdrawerAddress, nonces)
	if err != nil {
		log.Fatalf("genRegisterMsg failed: +v", err)
	}
	fmt.Println(res)

	// Update withdrawerAddress
	withdrawerAddress = common.HexToAddress("0x04A987fa1Bd4b2B908e9A3Ca058cc8BD43035991")
	res, err = gasBackMSGHelper.Invoke(client, senderAddress, privateKey, GasPrice, GasLimit, ChainId, "genUpdateMsg", &contractAddr, &withdrawerAddress)
	if err != nil {
		log.Printf("genUpdateMsg failed: +v", err)
	}
	fmt.Println(res)

	// Cancel gasback
	res, err = gasBackMSGHelper.Invoke(client, senderAddress, privateKey, GasPrice, GasLimit, ChainId, "genCancelMsg", &contractAddr)
	if err != nil {
		log.Printf("genCancelMsg failed: +v", err)
	}
	fmt.Println(res)

}
