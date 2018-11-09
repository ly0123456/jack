package main

import (
	"bytes"
	. "crypto/sha256"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
)


type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewPow(block *Block)*ProofOfWork  {
	target:=big.NewInt(1)
	 target.Lsh(target, 256-Diff)

	return &ProofOfWork{block,target}
}
func (pow *ProofOfWork)Run() ([]byte,uint64){
	var hashint big.Int
	var hash [32]byte
	nonce:=0
	for  {
		dataHash := pow.preDataHash(uint64(nonce))
		hash = Sum256(dataHash)
		hashint.SetBytes(hash[:])
		if pow.target.Cmp(&hashint)==1 {
			fmt.Println("找到了")
			break
		}
		nonce++
	}

	return hash[:],uint64(nonce)
}
func (pow *ProofOfWork)preDataHash(nonce uint64)[]byte{

	block:=pow.block

	joindata := bytes.Join([][]byte{
		num2bytes(block.Height),
		block.PrevBlockHash,
		block.Data,
		block.MerkelRoot,
		num2bytes(block.TimeSTamp),
		num2bytes(block.Diff),
		num2bytes(nonce),
	}, []byte{})
return joindata
}
func num2bytes(num uint64)[]byte{
	var buffer bytes.Buffer
	binary.Write(&buffer,binary.BigEndian,&num)
	return buffer.Bytes()
}



func (pow *ProofOfWork)IsVild()bool{

	dataHash := pow.preDataHash(pow.block.Nonce)
	hash := sha256.Sum256(dataHash)
	var hashint big.Int
	hashint.SetBytes(hash[:])
	if pow.target.Cmp(&hashint)==-1 {
		return false
	}
	return true
}