package zcrypto

import "math/big"

type Fujisaki struct {
	G,H,N *big.Int
}


func (this *Fujisaki) Commitment(m *big.Int) (*big.Int,error) {
	gv:= new(big.Int).Exp(this.G, m, this.N)
	r,err := RandomNumberInGroup(this.N)
	if err!=nil {
		return nil,err
	}
	hr:= new(big.Int).Exp(this.H,r,this.N)
	return new(big.Int).Mod(new(big.Int).Mul(gv,hr),this.N),nil
}
