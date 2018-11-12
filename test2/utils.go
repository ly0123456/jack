package main

import (
	"bytes"
	"encoding/binary"
	"github.com/base58"
	"os"
)

//数字转[]byte
func NumToBytes(num uint64) []byte {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, num)
	return buffer.Bytes()
}

//将地址转成公钥hash
func PubkeyHash(address string) []byte {
	decode := base58.Decode(address)
	pubkeyhash := decode[1 : len(decode)-4]
	return pubkeyhash
}

//查询当前文件是否存在
func Isexitis(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
