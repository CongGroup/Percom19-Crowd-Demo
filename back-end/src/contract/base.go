package contract

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"log"
	"math/big"
	"strings"
)


const (
	maxWaitingBlock = 50  // come from web3js, if wait more than 50 blocks, transaction time out
	GasLimit            = 2000000
	GasPrice =  4500000000
)

type ContractConfig struct {
	Port string `json:"port"`
	Address string `json:"address"`
	Abi string `json:"abi"`
}

type SendOpts struct {
	From common.Address
	To common.Address
	GasLimit uint64
	GasPrice *big.Int
	Value *big.Int
}

type BaseContract struct {
	Client   *ethclient.Client
	Address  common.Address
	ABI      *abi.ABI
}

func (c *BaseContract) Connect(socket string) {
	if c.Client != nil {
		return
	}
	var err error
	c.Client, err = ethclient.Dial(socket)
	if err != nil {
		panic(err.Error())
	}
}

func (c* BaseContract) LoadABI(contractAbi string) {
	var err error
	c.ABI= new(abi.ABI)
	*c.ABI, err = abi.JSON(strings.NewReader(contractAbi))
	if err !=nil{
		log.Fatal(err.Error())
	}
}

func (c *BaseContract) Close() error {
	if c.Client == nil {
		return errors.New("has not connected")
	}
	c.Client.Close()
	return nil
}

func (c *BaseContract) SendTransaction(tx *types.Transaction) (*types.Transaction, error) {
	err := c.Client.SendTransaction(context.Background(), tx)
	if err!=nil {
		log.Println(err.Error())
		return nil,err
	}
	_,err =c.GetReceiptStatus(tx.Hash())
	return tx, err
}

func (c *BaseContract) GetEther(address common.Address) (*big.Int, error) {
	balance, err := c.Client.BalanceAt(context.Background(), address, nil)
	return balance, err
}

func (c *BaseContract) GetNonce(address common.Address) (uint64, error) {
	nonce, err := c.Client.PendingNonceAt(context.Background(), address)
	return nonce, err
}

func (c *BaseContract) GetChainId() (*big.Int,error) {
	chainID, err := c.Client.NetworkID(context.Background())
	return chainID,err
}

func (c *BaseContract) GenerateKeyStore(file string, password string) (common.Address, error) {
	ks := keystore.NewKeyStore(file, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	keyJson, err := ks.Export(account, password, password)
	ioutil.WriteFile(file, keyJson, 0777)
	return account.Address, err
}

func (c *BaseContract) GetAddress() common.Address {
	return c.Address
}

func (c *BaseContract) GetReceiptStatus (txHash common.Hash) (uint64,error) {
	count:= 0
	ch:= make(chan *types.Header)
	sub,err:= c.Client.SubscribeNewHead(context.Background(),ch)
	if err!=nil {
		log.Println(err.Error())
		return 0,err
	}

	for {
		select {
		case err:= <- sub.Err():
			log.Println(err.Error())
			return 0,err
		case <- ch:
			count+=1
			if count>=maxWaitingBlock {
				return 0, errors.New("transaction time out")
			} else {
				receipt, err:= c.Client.TransactionReceipt(context.Background(),txHash)
				if err == nil {
					//log.Println(receipt.ContractAddress.String())
					//log.Println(receipt.TxHash.String())
					if receipt.Status ==0 {
						return receipt.Status, errors.New("transaction revert")
					}
					return receipt.Status, nil
				} else {
				}
			}
		}
	}
}

func (c *BaseContract) Call(funcName string, args ...interface{}) ( []byte, error) {
	input,err:= c.ABI.Pack(funcName, args...)
	if err!=nil {
		return nil,err
	}
	rVal, err:= c.Client.CallContract(context.Background(),ethereum.CallMsg{
		To: &c.Address,
		Data:input,
	},nil)
	return rVal,err
}

func (c *BaseContract) PackFunction(funcName string, args ...interface{}) ([]byte, error) {
	input,err:= c.ABI.Pack(funcName, args...)
	if err!=nil {
		return nil,err
	}
	return input,nil
}

func (c *BaseContract) EventWatcher() (chan types.Log, <-chan error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{c.Address},
	}
	logs := make(chan types.Log)
	sub, err := c.Client.SubscribeFilterLogs(context.Background(), query, logs)

	if err != nil {
		log.Fatal(err)
	}

	return logs, sub.Err()
}

func (c *BaseContract) Unpack(v interface{}, name string, output []byte) error {
	err:=c.ABI.Unpack(v, name, output)
	return err
}