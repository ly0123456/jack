package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

type Wallet struct {
	privateKey *ecdsa.PrivateKey
	Pubkey     []byte
}

//生成私钥对
func NewWallet() *Wallet {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubkey := privateKey.PublicKey

	pubkey1 := append(pubkey.X.Bytes(), pubkey.Y.Bytes()...)
	return &Wallet{privateKey: privateKey, Pubkey: pubkey1}
}

//生成钱包地址
func (w *Wallet) GetAddress() string {
	pubkey := HashPubkey(w.Pubkey)
	Vision := []byte{byte(00)}
	pubkey = append(Vision, pubkey...)
	hash4 := TwoSha256Hash(pubkey)
	pubkey = append(pubkey, hash4...)
	address := base58.Encode(pubkey)
	return address
}
func HashPubkey(pubkey []byte) []byte {
	hash := sha256.Sum256(pubkey)
	ripemder := ripemd160.New()
	ripemder.Write(hash[:])
	pubkeyhash := ripemder.Sum(nil)
	return pubkeyhash
}

//做两次sha256
func TwoSha256Hash(pubkey []byte) []byte {
	sum256 := sha256.Sum256(pubkey)
	bytes := sha256.Sum256(sum256[:])
	return bytes[:4]
}

func Isaddress(address string)bool{
	decodeInfo := base58.Decode(address)
	pubkeyhash:=decodeInfo[1:len(decodeInfo)-4]
	checksum1 := decodeInfo[len(decodeInfo)-4:]
	hash := TwoSha256Hash(pubkeyhash)
	checksum2:=hash[0:4]
	if !bytes.Equal(checksum1,checksum2){
		return false
	}

	return true
}
