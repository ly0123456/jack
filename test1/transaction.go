package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

//交易输入
type TxInput struct {
	TXId  []byte //引用交易ID
	Index uint64 // 所在ID的索引
	Sig   string //解密脚本
}

//交易输出
type TxOutput struct {
	Value        float64 //转账金额
	ScriptPubKey string  //转账地址
}

//交易
type Transaction struct {
	TXHash  []byte //交易ID
	Inputs  []*TxInput
	Outputs []*TxOutput
}

func (t *Transaction) SetId() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(t)
	hash := sha256.Sum256(buffer.Bytes())
	t.TXHash = hash[:]
}

//挖矿交易
func NewCoinbaseTx(address, data string) *Transaction {
	input := TxInput{nil, 0, data}
	output := TxOutput{12.5, address}
	tx := Transaction{nil, []*TxInput{&input}, []*TxOutput{&output}}
	tx.SetId()
	return &tx
}

//创建一个普通交易
func NewTransaction(from, to string, amount uint64, blc *BlockChain) *Transaction {
	var inputs []*TxInput
	var outputs []*TxOutput
	//通过我的名字找到我的utxo
	Needutxos, calcMoney := blc.FindNeedUtxos(from, amount)
	if calcMoney < amount {
		fmt.Println("余额不足，交易失败")
		return nil
	}
	//便利我的utxo
	for i, indexds := range Needutxos {
		for _, index := range indexds {
			input := TxInput{TXId: []byte(i), Index: index, Sig: from}
			inputs = append(inputs, &input)
		}
	}
	output := TxOutput{float64(calcMoney), to}
	outputs = append(outputs, &output)
	//如果找到的钱比支付的钱多久需要找零
	if calcMoney > amount {
		outputs = append(outputs, &TxOutput{float64(calcMoney - amount), from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetId()
	return &tx
}
