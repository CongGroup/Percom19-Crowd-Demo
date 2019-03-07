package app

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"log"
)

type Wallet struct {
	Address common.Address  `json:"address"`
	PrivateKey *ecdsa.PrivateKey `json:"wallet"`
}

type KeyStore struct {
	Address common.Address `json:"address"`
	PrivateKey string `json:"privateKey"`
}

func NewWallet() *Wallet {
	privateKeyECDSA,err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err!=nil {
		log.Println(err.Error())
		panic(err)
	}
	address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
	return &Wallet{
		Address:address,
		PrivateKey:privateKeyECDSA,
	}
}

func NewWalletFromFile(filepath string) *Wallet{
	var keystore KeyStore
	data,err:=ioutil.ReadFile(filepath)
	if err!=nil {
		log.Println(err.Error())
		panic(err)
	}
	err =json.Unmarshal(data,&keystore)
	if err!=nil {
		log.Println(err.Error())
		panic(err)
	}

	pk,err:= derivePrivateKey(keystore.PrivateKey)
	if err!=nil {
		log.Println(err.Error())
		panic(err)
	}

	return &Wallet{
		Address:keystore.Address,
		PrivateKey:pk,
	}
}

func derivePrivateKey(rawString string) (*ecdsa.PrivateKey,error) {
	pk,err:= crypto.HexToECDSA(rawString)
	if err!=nil {
		return nil,err
	}
	return pk,nil
}
