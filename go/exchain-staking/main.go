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
	"github.com/okex/solidity-sample/go/exchain-staking/contracts"
	"log"
)

const (
	RpcUrl          = "https://exchaintestrpc.okex.org"
	ChainId  int64  = 65 //  oec testnet
	PrivKey         = ""
	GasPrice int64  = 100000000 // 0.1 gwei
	GasLimit uint64 = 3000000
)

const (
	// define testnet systemContractAddress. mainnet systemContractAddress is here: https://www.oklink.com/zh-cn/okc/address/0xe9196e65a0b6705777fbe829dfa94ec8b9f2ba48
	systemContarctAddress = "0x7e5E6AF6424BE7A835313777Fc6E0d1912e52Fc8"
	//  // define testnet stakingContractAddress. mainnet stakingContractAddress is here: https://www.oklink.com/zh-cn/okc/address/0x1b29c875bd7ec9a12c29fc6eef8e451207352ef3
	stakingMSGHelperContractAddress = "0xE6947172736e3ed558Ec375cf16C62488627a18e"
)

func main() {

	// load contracts
	contractAbi, err := abi.JSON(bytes2.NewReader([]byte(contracts.SystemContractABI)))
	if err != nil {
		panic(err)
	}
	systemContract := contracts.SystemContract{ByteCode: common.Hex2Bytes(contracts.SystemContractBIN), Abi: contractAbi, Address: common.HexToAddress(systemContarctAddress)}

	contractAbi, err = abi.JSON(bytes2.NewReader([]byte(contracts.StakingMSGHelperABI)))
	if err != nil {
		panic(err)
	}

	stakingMSGHelper := contracts.ContractWrap{ByteCode: common.Hex2Bytes(contracts.StakingMSGHelperABI), Abi: contractAbi, Address: common.HexToAddress(stakingMSGHelperContractAddress), System: &systemContract}

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

	// Deposit
	res, err := stakingMSGHelper.Invoke(client, senderAddress, privateKey, GasPrice, GasLimit, ChainId, "genDepositMsg", "10")
	if err != nil {
		log.Fatalf("genDepositMsg failed: +v", err)
	}
	fmt.Println(res)

	// Withdraw
	res, err = stakingMSGHelper.Invoke(client, senderAddress, privateKey, GasPrice, GasLimit, ChainId, "genWithdrawMsg", "2")
	if err != nil {
		log.Printf("withdraw failed: +v", err)
	}
	fmt.Println(res)

	// AddShares
	vAddress := []string{"exvaloper15w73dl3n7qmrk8ssax5w2k45qd29jlkejyc23q"}
	res, err = stakingMSGHelper.Invoke(client, senderAddress, privateKey, GasPrice, GasLimit, ChainId, "genAddSharesMsg", vAddress)
	if err != nil {
		log.Printf("addshares failed: +v", err)
	}
	fmt.Println(res)

	res, err = stakingMSGHelper.Invoke(client, senderAddress, privateKey, GasPrice, GasLimit, ChainId, "genWithdrawAllRewardsMsg")
	if err != nil {
		log.Printf("WithdrawAllReward failed: +v", err)
	}
	fmt.Println(res)
}
