package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/base58"
	"log"
)

type Wallet struct {
	privateKey *ecdsa.PrivateKey
	Pubkey []byte
}
//生成私钥对
func NewWallet()*Wallet{
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubkey:=privateKey.PublicKey

	pubkey1:= append(pubkey.X.Bytes(), pubkey.Y.Bytes()...)
	return &Wallet{privateKey:privateKey,Pubkey:pubkey1}
}
//生成钱包地址
func (w *Wallet)GetAddress()string{
	pubkey := HashPubkey(w.Pubkey)
	Vision:=[]byte{00}
	pubkey= append(pubkey, Vision...)
	hash4 := TwoSha256Hash(pubkey)
	pubkey=append(pubkey, hash4...)
	address := base58.Encode(pubkey)
	return  address
}
