package contracts

import (
	"context"
	"crypto/ecdsa"
	_ "embed"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

var (

	//go:embed GasBackMSGHelper.json
	GasBackMSGHelperABI string
	//go:embed GasBackMSGHelper.bin
	GasBackMSGHelperBIN string

	//go:embed SystemContract.json
	SystemContractABI string
	//go:embed SystemContract.bin
	SystemContractBIN string
)

type SystemContract struct {
	ByteCode []byte
	Abi      abi.ABI
	Address  common.Address
}

func (s SystemContract) encodeData(data interface{}) ([]byte, error) {
	ret, err := s.Abi.Pack("invoke", data)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type ContractWrap struct {
	ByteCode []byte
	Abi      abi.ABI
	Address  common.Address
	System   *SystemContract
}

func (c ContractWrap) Invoke(client *ethclient.Client,
	fromAddress common.Address,
	privateKey *ecdsa.PrivateKey,
	gasPrice int64,
	gasLimit uint64,
	chainid int64,
	name string,
	args ...interface{}) (*types.Receipt, error) {
	// 0. get the value of nonce, based on address
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the value of nonce from network: %+v", err)
	}

	// 0.5 get the gasPrice
	gasprice := big.NewInt(gasPrice)

	// print transaction information
	fmt.Printf(
		"writeContract: \n"+
			"	sender address: <%s>\n"+
			"	contract address: <%s>\n"+
			"	gas price: <%.1f> gwei\n"+
			"	nonce: %d\n"+
			"	ABI: <%s %s>\n",
		fromAddress.Hex(),
		c.System.Address,
		float64(gasprice.Uint64())/1e9,
		nonce, name, args)

	data, err := c.Abi.Pack(name, args...)
	if err != nil {
		return nil, err
	}

	// gen msg
	processedData, err := client.CallContract(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &c.Address,
		Data: data,
	}, nil)

	if err != nil {
		return nil, fmt.Errorf("process data failed: +v", err)
	}

	// unpack result
	ret, err := c.Abi.Unpack(name, processedData)
	if err != nil {
		return nil, err
	}

	// invoke system contract

	// 1. encode data using system contract
	data, err = c.System.encodeData(ret[0])
	if err != nil {
		return nil, err
	}

	unsignedTx := types.NewTransaction(nonce, c.System.Address, big.NewInt(0), gasLimit, gasprice, data)
	// 2. sign unsignedTx -> rawTx
	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(big.NewInt(chainid)), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign the unsignedTx offline: %+v", err)
	}

	// 3. send rawTx
	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		return nil, err
	}

	return getReceipt(client, signedTx)
}

func getReceipt(client *ethclient.Client, signedTx *types.Transaction) (*types.Receipt, error) {

	hash := signedTx.Hash()
	log.Println(hash.String())
	time.Sleep(time.Second * 15)

	receipt, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
