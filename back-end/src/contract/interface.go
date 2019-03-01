package contract

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Contract interface {
	GetAddress() common.Address
	Connect(socket string)
	Close() error
	SendTransaction(tx *types.Transaction) (*types.Transaction, error)
	GetReceiptStatus (txHash common.Hash) (uint64,error)
	Call(funcName string, args ...interface{}) ( []byte, error)
	PackFunction(funcName string, args ...interface{}) ([]byte, error)
	GetNonce(address common.Address) (uint64, error)
	GetChainId() (*big.Int,error)
	EventWatcher() (chan types.Log, <-chan error)
	Unpack(v interface{}, name string, output []byte) error
}
