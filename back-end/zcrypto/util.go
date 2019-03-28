package zcrypto

import (
	"crypto/rand"
	"math/big"
)

func RandomNumberInGroup(N *big.Int) (*big.Int,error) {
	n,err:= rand.Int(rand.Reader, new(big.Int).Sub(N,big.NewInt(1)))
	return new(big.Int).Add(n,big.NewInt(1)),err
}