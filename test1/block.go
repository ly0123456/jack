package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)
//系统设置难度
const Diff  =20
//定义一个区块的类
type Block struct {
	Vision        uint64 //版本号
	Height        uint64 //区块高度
	PrevBlockHash []byte //前区块hash
	MorkHash      []byte //默克尔根
	Timestamp     uint64 //时间戳
	Data          []byte //交易数据
	Hash          []byte //本区块hash
	Diff          uint64 //当前难度
	Nonce         uint64 //随机值
	Txs           []*Transaction
}
//实例化Block
func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	block:=Block{
		Vision:        0,
		Height:        0,
		PrevBlockHash: prevHash,
		MorkHash:      nil,
		Timestamp:     uint64(time.Now().Unix()),
		Txs:           txs,
		Diff:          Diff,
	}
	//运用pow算出当前hash nonce
	//创建新的工作量证明
	block.setMorkRoot()
	pow:=NewPow(&block)
	//使用工作量证明的方法，返回hash nonce
	hash,nonce:=pow.Run()
	block.Hash=hash
	block.Nonce=nonce

	return  &block
}
//实现一个block的序列化的方法
func (b *Block)Encode()[]byte  {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer) //先定义一个编码器
	//然后用编码器编码
	encoder.Encode(&b)
	return  buffer.Bytes()
}
//实现一个block的反序列化的方法
func Decode(data []byte)*Block  {
	var block Block
	//定义一个解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))
	decoder.Decode(&block)
	return &block
}
func (b *Block) setMorkRoot() {
	info := []byte{}
	for _, tx := range b.Txs {
		info = append(info, tx.TXHash...)
	}
	hash := sha256.Sum256(info)
	b.MorkHash = hash[:]
}