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

//var walletLock sync.Mutex

type Wallet struct {
	Address common.Address  `json:"address"`
	PrivateKey *ecdsa.PrivateKey `json:"wallet"`
}

type KeyStore struct {
	Address common.Address `json:"address"`
	PrivateKey string `json:"privateKey"`
}

func NewWallet() *Wallet {
	for {
		privateKeyECDSA,err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
		if err!=nil {
			log.Println(err.Error())
			panic(err)
		}
		if(len(privateKeyECDSA.D.Text(16))!=64) {
			continue;
		}
		address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
		return &Wallet{
			Address:address,
			PrivateKey:privateKeyECDSA,
		}
	}
}

func NewWalletFromFile(filepath string) (*Wallet,error){
	var keystore KeyStore
	data,err:=ioutil.ReadFile(filepath)
	if err!=nil {
		log.Println(err.Error())
		return nil,err
	}
	err =json.Unmarshal(data,&keystore)
	if err!=nil {
		log.Println(err.Error())
		return nil,err
	}

	pk,err:= derivePrivateKey(keystore.PrivateKey)
	if err!=nil {
		log.Println(keystore.PrivateKey,err.Error())
		return nil,err
	}

	return &Wallet{
		Address:keystore.Address,
		PrivateKey:pk,
	},nil
}

func derivePrivateKey(rawString string) (*ecdsa.PrivateKey,error) {
	pk,err:= crypto.HexToECDSA(rawString)
	if err!=nil {
		return nil,err
	}
	return pk,nil
}
