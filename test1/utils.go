package main

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"os"
)
//判断文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
//将pubkey转化为pubkeyhash
func HashPubkey(pubkey []byte)[]byte{
	hash := sha256.Sum256(pubkey)
	ripemder := ripemd160.New()
	ripemder.Write(hash[:])
	pubkeyhash := ripemder.Sum(nil)
	return  pubkeyhash
}
//做两次sha256
func TwoSha256Hash(pubkey []byte)[]byte{
	sum256 := sha256.Sum256(pubkey)
	bytes := sha256.Sum256(sum256[:])
	return bytes[:4]
}