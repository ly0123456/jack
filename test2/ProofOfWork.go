package main

import (
	"bytes"
	. "crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	Block *Block
	//目标值
	target *big.Int
}

func NewPow(block Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, 256-Diff)
	return &ProofOfWork{Block: &block, target: target}
}
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var nonce uint64
	var hash []byte
	for {
		//获取block的当前hash
		hash = pow.preBlockHash(nonce)
		//定义一个big。int容器来接hash用于与target比较
		var hashBig big.Int
		hashBig.SetBytes(hash)
		if pow.target.Cmp(&hashBig) == 1 {
			fmt.Println("挖矿成功。。。。。")
			break
		}
		nonce++
	}
	return hash, nonce
}

//拼装数据
func (pow *ProofOfWork) preBlockHash(nonce uint64) []byte {
	block := pow.Block
	data := bytes.Join([][]byte{
		NumToBytes(block.Version),
		block.PrevBlockHash,
		block.MerkelRoot,
		NumToBytes(block.Difficulty),
		NumToBytes(block.TimeStamp),
		NumToBytes(nonce),
	}, []byte{})
	hash := Sum256(data)
	return hash[:]
}

//校验数据
func (pow *ProofOfWork) IsVaild(nonce uint64) bool {
	//获取block的当前hash
	hash := pow.preBlockHash(nonce)
	//定义一个big。int容器来接hash用于与target比较
	var hashBig big.Int
	hashBig.SetBytes(hash)
	if pow.target.Cmp(&hashBig) == -1 {
		return false
	}
	return true
}
