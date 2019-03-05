package user

import (
	"contract"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
)

type User struct {
	Address common.Address
	privateKey *ecdsa.PrivateKey
	c contract.Contract
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

	nonce, err:= this.c.GetNonce(this.Address)
	if err!=nil {
		log.Println(err.Error())
		return err
	}

	tx:= types.NewTransaction(nonce,this.c.GetAddress(),value,contract.GasLimit,big.NewInt(contract.GasPrice),input);
	chainID, err := this.c.GetChainId()
	if err!=nil {
		return err
	}

	signedTx, err:= types.SignTx(tx,types.NewEIP155Signer(chainID),this.privateKey)
	//log.Println("contract Address:",signedTx.To().String())
	//log.Println("hash:", signedTx.Hash().String())
	_,err = this.c.SendTransaction(signedTx)
	if err!=nil {
		return err
	}
	_,err = this.c.GetReceiptStatus(signedTx.Hash())
	return err
}

