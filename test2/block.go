package main

import (
	"bytes"
	"crypto/sha256"
	. "encoding/gob"
	"time"
)

//当前系统难度
const Diff = 20
const genesisInfo = "2009年1月3日，财政大臣正处于实施第二轮银行紧急援助的边缘"

type Block struct {
	Version       uint64 //版本号
	PrevBlockHash []byte //前区块哈希值

	MerkelRoot []byte //这是一个哈希值，后面v5用到

	TimeStamp uint64 //时间戳，从1970.1.1到现在的秒数

	Difficulty uint64 //通过这个数字，算出一个哈希值：0x00010000000xxx

	Nonce uint64 // 这是我们要找的随机数，挖矿就找证书

	Hash []byte //当前区块哈希值, 正常的区块不存在，我们为了方便放进来

	//Data []byte //数据本身，区块体，先用字符串表示，v4版本的时候会引用真正的交易结构
	Transactions []*Transaction
}

//创建一个新区块
func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := Block{
		Version:       0,
		PrevBlockHash: prevBlockHash,
		TimeStamp:     uint64(time.Now().Unix()),
		Difficulty:    Diff,
		Transactions:  txs,
	}
	block.MerkelRoot = block.GetMerkelRoot()
	//根据pow生成当前hash和nonce
	pow := NewPow(block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

//根据交易生成merkelroot
func (bl *Block) GetMerkelRoot() []byte {
	var merkroot []byte
	var buffer bytes.Buffer
	for _, tx := range bl.Transactions {
		encoder := NewEncoder(&buffer)
		encoder.Encode(&tx)
		merkroot = append(merkroot, buffer.Bytes()...)
	}
	merkhash := sha256.Sum256(merkroot)
	return merkhash[:]
}

//使用gob包将区块数据转换成[]byte
func (bl *Block) Encode() []byte {
	var bufer bytes.Buffer
	encoder := NewEncoder(&bufer)
	encoder.Encode(&bl)
	return bufer.Bytes()
}

//将[]byte转化成block

func Decode(data []byte) *Block {
	var block Block
	decoder := NewDecoder(bytes.NewReader(data))
	decoder.Decode(&block)
	return &block
}
