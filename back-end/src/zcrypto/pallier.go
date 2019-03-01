package zcrypto

import (
	"crypto/rand"
	"math/big"
)

const (
	N = "10608696505376415216836126482574953390792074036402633331552687092108694586747202641175339209267949553037389623984664304482170753400853800057208033993322043"
	G = "10608696505376415216836126482574953390792074036402633331552687092108694586747202641175339209267949553037389623984664304482170753400853800057208033993322044"
	Lambda = "10608696505376415216836126482574953390792074036402633331552687092108694586746996351441634229506849036970974238216158956473092263641175734091312055291354760"
	Mu = "2187834066285000450347932489094704936460943388243677631913843275610372585229307056365408569652759734622347832192372613093504606689716692392562656706677656"
)

var PubKey *PublicKey
var PriKey *PrivateKey

type PublicKey struct {
	N, G *big.Int
	n2 *big.Int
}

type Cypher struct {
	C *big.Int
}

type PrivateKey struct {
	PublicKey
	Lambda *big.Int
	Mu *big.Int
}

func (this *PublicKey) GetNSquare() *big.Int {
	if this.n2 !=nil {
		return this.n2
	}
	this.n2 = new(big.Int).Mul(this.N,this.N)
	return this.n2
}

func LCM(x, y *big.Int) *big.Int {
	return new(big.Int).Mul(new(big.Int).Div(x, new(big.Int).GCD(nil, nil, x, y)), y)
}

func minusOne(x *big.Int) *big.Int {
	return new(big.Int).Add(x, big.NewInt(-1))
}

func computeMu(g, lambda, n *big.Int) *big.Int {
	n2 := new(big.Int).Mul(n, n)
	u := new(big.Int).Exp(g, lambda, n2)
	return new(big.Int).ModInverse(L(u, n), n)
}

func L(u, n *big.Int) *big.Int {
	t := new(big.Int).Add(u, big.NewInt(-1))
	return new(big.Int).Div(t, n)
}

func NewPallier(bits int) (*PublicKey,*PrivateKey,error) {
	random:= rand.Reader
	var p,q *big.Int
	var errChan = make(chan error,1)
	go func(){
		var err error
		p, err = rand.Prime(random,bits)
		errChan <-err
	}()

	q, err:= rand.Prime(random,bits)
	if err!=nil {
		return nil,nil,err
	}


	n := new(big.Int).Mul(p, q)
	lambda := new(big.Int).Mul(minusOne(p), minusOne(q))
	g := new(big.Int).Add(n, big.NewInt(1))
	mu := new(big.Int).ModInverse(lambda, n)

	pub:= &PublicKey{
		N:n,
		G:g,
	}

	pri:= &PrivateKey{
		PublicKey: PublicKey {
			N: new(big.Int).Set(n),
			G: new(big.Int).Set(g),
		},
		Lambda:lambda,
		Mu:mu,
	}

	return pub,pri,nil
}

func (this *PublicKey)Encrypt(m *big.Int) (*Cypher,error) {
	r,err:= RandomNumberInGroup(this.N)
	if err != nil {
		return nil, err
	}
	nSquare := this.GetNSquare()

	gm := new(big.Int).Exp(this.G, m, nSquare)
	rn := new(big.Int).Exp(r, this.N, nSquare)
	return &Cypher{new(big.Int).Mod(new(big.Int).Mul(rn, gm), nSquare)}, nil
}

func (this *PrivateKey) Decrypt(cypher *Cypher) (*big.Int) {
	tmp := new(big.Int).Exp(cypher.C, this.Lambda, this.GetNSquare())
	msg := new(big.Int).Mod(new(big.Int).Mul(L(tmp, this.N), this.Mu), this.N)
	return msg
}

func init() {
	_n , _ := new(big.Int).SetString(N,10)
	_g , _ := new(big.Int).SetString(N,10)
	_l,_:= new(big.Int).SetString(N,10)
	_m,_:= new(big.Int).SetString(N,10)
	PubKey = &PublicKey{
		N: _n,
		G: _g,
	}
	PriKey = &PrivateKey {
		PublicKey:PublicKey{
			N:_n,
			G:_g,
		},
		Lambda:_l,
		Mu: _m,
	}
}