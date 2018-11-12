package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/base58"
)

//交易输入
type TxInput struct {
	TXId  []byte //引用交易ID
	Index uint64 // 所在ID的索引
	Sig   []byte //解密脚本
	Pubkey []byte
 }

//交易输出
type TxOutput struct {
	Value        float64 //转账金额
	PubkeyHash  []byte  //转账地址
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
func NewOutput(value float64,address string) *TxOutput {
	var Output *TxOutput
	Output.Value=value
	Output.LockWithHash(address)
	return Output
}
func (output *TxOutput )LockWithHash(address string){
	decode := base58.Decode(address)
	pubkeyHash:=decode[1:len(decode)-4]
	output.PubkeyHash=pubkeyHash
}

const reward = 12.5
//挖矿交易
func NewCoinbaseTx(address,data string) *Transaction {
	input := TxInput{nil, 0, nil ,[]byte(data)}
	txOutput := NewOutput(reward, address)
	tx := Transaction{nil, []*TxInput{&input}, []*TxOutput{txOutput}}
	tx.SetId()
	return &tx
}
func (tx *Transaction) IsCoinbaseTx() bool {
	if tx.Inputs[0].TXId == nil && len(tx.Inputs) == 1 && tx.Inputs[0].Index == 0 {
		return true
	}
	return false
}

//创建一个普通交易
func NewTransaction(from, to string, amount float64, blc *BlockChain) *Transaction {
	wallets := NewWallets()
	if wallets.Wallets[from] == nil {
		fmt.Printf("本地没有 %s 的钱包，无法创建交易\n", from)
		return nil
	}
	wallet := wallets.Wallets[from]
	//privateKey:=wallet.privateKey
	pubkey:=wallet.Pubkey
	pubkeyHash:=HashPubkey(pubkey)
	var inputs []*TxInput
	var outputs []*TxOutput
	//通过我的名字找到我的utxo
	Needutxos, calcMoney := blc.FindNeedUtxos(pubkeyHash, amount)
	if calcMoney < amount {
		fmt.Println("余额不足，交易失败")
		return nil
	}
	//便利我的utxo
	for i, indexds := range Needutxos {
		for _, index := range indexds {
			input := TxInput{TXId: []byte(i), Index: uint64(index), Sig: nil,Pubkey:pubkey}
			inputs = append(inputs, &input)
		}
	}
	output := NewOutput(amount,to)
	outputs = append(outputs, output)
	//如果找到的钱比支付的钱多久需要找零
	if calcMoney > amount {
		output1:=NewOutput(calcMoney-amount,from)
		outputs = append(outputs, output1)
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetId()
	return &tx
}
