package contract

import (
	"github.com/ethereum/go-ethereum/common"
)

const (
	FUNCTION_BALANCE_OF = "balanceOf"
)

type ERC20 struct {
	BaseContract
}

func NewERC20 (port string, contractABI string ,contractAddress common.Address) *ERC20 {
	token := new(ERC20)
	token.Connect(port)
	token.Address = contractAddress
	token.LoadABI(contractABI)
	return token
}

