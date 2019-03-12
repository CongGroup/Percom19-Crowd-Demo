package user

import (
	"contract"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"sync/atomic"
)

const (
	MASTER_ADDRESS = "0x3c62aa7913bc303ee4b9c07df87b556b6770e3fc"
	MASTER_KEY = "e27cb51d1eb94ad42b8f196e341e082042639677df43fd7d1440c07b40e2a065"
	TRANSFER_VALUE = "20000000000000000"   // 0.02 ether
	THRESHOLD_VALUE = "10000000000000000"  // 0.01 ether only send ether if less than THRESHOLD_VALUE
)

type User struct {
	Address common.Address
	privateKey *ecdsa.PrivateKey
	c contract.Contract
	Nonce uint64  // for cache
}

func NewUser(address common.Address, privateKey *ecdsa.PrivateKey, c contract.Contract) *User {
	return &User{
		Address:address,
		privateKey: privateKey,
		c:c,
	}
}

func (this *User) Send(funcName string, args ...interface{}) error {
	return this.SendWithValue(big.NewInt(0),funcName,args...)
}

func (this *User) SendWithValue(value *big.Int,funcName string, args ...interface{}) error {
	input,err:= this.c.PackFunction(funcName,args...)
	if err!=nil {
		return err
	}

	return this.sendHelper(value,input)
}

func (this *User) sendHelper(value *big.Int, input []byte) error {
	nonce, err:= this.c.GetNonce(this.Address)
	if err!=nil {
		return err
	}
	tx:= types.NewTransaction(nonce,this.c.GetAddress(),value,contract.GasLimit,big.NewInt(contract.GasPrice),input);
	chainID, err := this.c.GetChainId()
	if err!=nil {
		return err
	}

	signedTx, err:= types.SignTx(tx,types.NewEIP155Signer(chainID),this.privateKey)

	_,err = this.c.SendTransaction(signedTx)
	if err!=nil {
		return err
	}
	_,err = this.c.GetReceiptStatus(signedTx.Hash())
	return err
}


// only for owner user
func (this *User) Transfer(to common.Address, value *big.Int) error {
	nonce:=atomic.AddUint64(&this.Nonce,1)

	tx:= types.NewTransaction(nonce-1,to,value,contract.GasLimit,big.NewInt(contract.GasPrice),[]byte("0x"));
	chainID, err := this.c.GetChainId()
	if err!=nil {
		return err
	}
	signedTx, err:= types.SignTx(tx,types.NewEIP155Signer(chainID),this.privateKey)

	_,err = this.c.SendTransaction(signedTx)
	if err!=nil {
		return err
	}
	_,err = this.c.GetReceiptStatus(signedTx.Hash())
	return err
}

