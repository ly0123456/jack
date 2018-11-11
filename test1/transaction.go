package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

type TxInput struct {
	TXId  []byte //引用交易ID
	Index uint64 // 所在ID的索引
	Sig   string //解密脚本
}
type TxOutput struct {
	Value        float64 //转账金额
	ScriptPubKey string  //转账地址
}
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

func NewCoinbaseTx(address, data string) *Transaction {
	input := TxInput{nil, -1, data}
	output := TxOutput{12.5, address}
	tx := Transaction{nil, []*TxInput{&input}, []*TxOutput{&output}}
	tx.SetId()
	return &tx
}
func NewTransaction(from ,to string ,amount uint64,blc *BlockChain )*Transaction{
	var inputs []*TxInput
	var outputs []*TxOutput
	Needutxos,calcMoney:=blc.FindNeedUtxos(from,amount)
	if calcMoney<amount {
		fmt.Println("余额不足，交易失败")
		return nil
	}
	for i,indexds:=range Needutxos  {
		for _,index:=range indexds {
			input:=TxInput{TXId:[]byte(i),Index:index,Sig:from}
					inputs=append(inputs,&input)
		}
	}
	output:=TxOutput{float64(calcMoney),to}
	outputs = append(outputs, &output)
	if calcMoney>amount {
		outputs=append(outputs,&TxOutput{float64(calcMoney-amount),from})
	}


	tx := Transaction{nil, inputs, outputs}
	tx.SetId()
	return &tx
}
