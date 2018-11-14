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
