package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

type Output struct {
	//交易金额
	Value float64
	//收款人的公钥hash
	PubkeyHash []byte
}
type Input struct {
	//引用的交易id
	Txid []byte
	//索引
	Index int64
	//私钥签名
	Sig []byte
	//自己的公钥
	Pubkey []byte
}
type Transaction struct {
	//交易Id
	TxId []byte
	//所有的输出
	TxOutPuts []*Output
	//所有的输入
	TxIntPuts []*Input
	//交易时间
	TimeStamp int64
}

func NewCoinbaseTx(adderss, data string) *Transaction {
	//输入
	input := Input{nil, -1, nil, []byte(data)}
	//输出
	output := Output{reward, PubkeyHash(adderss)}
	tx := Transaction{nil, []*Output{&output}, []*Input{&input}, time.Now().Unix()}
	tx.TxId = tx.SetHash()
	return &tx
}

//用交易信息生成交易hash作为交易id
func (tx *Transaction) SetHash() []byte {
	var bufer bytes.Buffer
	encoder := gob.NewEncoder(&bufer)
	encoder.Encode(&tx)
	hash := sha256.Sum256(bufer.Bytes())
	return hash[:]
}
