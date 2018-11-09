package main

import (
	"bytes"
	. "encoding/gob"
	"time"
)
const Diff  =20
type Block struct {
	Height        uint64
	PrevBlockHash []byte
	TimeSTamp     uint64
	Data          []byte
	Hash          []byte
	MerkelRoot    []byte
	Diff          uint64
	Nonce         uint64
}

func NewBlock(data string ,prevBlockHash []byte)* Block {
	block:=Block{
		Height:0,
		PrevBlockHash:prevBlockHash,
		TimeSTamp:uint64(time.Now().Unix()),
		Data:[]byte(data),
		MerkelRoot:[]byte{},
		Diff:Diff,
	}
	//使用pow算出当前hash和 Nonce
	pow:=NewPow(&block)
	hash,nonce:=pow.Run()
	block.Hash=hash
	block.Nonce=nonce
	return &block
}
func (b *Block)EncSec()[]byte{
	var buferr bytes.Buffer
	encoder := NewEncoder(&buferr)
	encoder.Encode(b)
	return buferr.Bytes()
}

func Decode(data []byte)*Block{
	var block  Block
	decoder := NewDecoder(bytes.NewReader(data))
	decoder.Decode(&block)
return  &block
}