package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallte struct {
	//私钥
	PrivateKey *ecdsa.PrivateKey
	//公钥的字节流
	Pubkey []byte
}

//生成一个私钥对
func NewWallre() *Wallte {
	//使用椭圆曲线加密函数
	privateKey, e := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if e != nil {
		fmt.Println("钱包创建失败")
		return nil
	}
	pubkey := privateKey.PublicKey
	//将公钥用[]byte转化用于传输
	pubkey1 := append(pubkey.X.Bytes(), pubkey.Y.Bytes()...)
	return &Wallte{PrivateKey: privateKey, Pubkey: pubkey1}
}

//用公钥生成一个地址
func (w *Wallte) CreateAddress() string {
	pubkey1 := HashPubkey(w.Pubkey)
	Vosdon := byte(00)
	pubkey := append([]byte{Vosdon}, pubkey1...)
	hash4 := TwoHash(pubkey)
	lasthash := append(pubkey, hash4...)
	address := base58.Encode(lasthash)
	return address
}
func HashPubkey(pubkey []byte) []byte {
	ripemder := ripemd160.New()
	ripemder.Write(pubkey)
	hash := ripemder.Sum(nil)
	return hash[:]
}
func TwoHash(pubkey []byte) []byte {
	data := sha256.Sum256(pubkey)
	hash := sha256.Sum256(data[:])
	return hash[:4]
}
