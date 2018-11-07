package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/base58"
	"log"
	"math/big"
	"strings"
)

//交易输入
type TxInput struct {
	TXId   []byte //引用交易ID
	Index  uint64 // 所在ID的索引
	Sig    []byte //解密脚本
	Pubkey []byte
}

//交易输出
type TxOutput struct {
	Value      float64 //转账金额
	PubkeyHash []byte  //转账地址
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
func NewOutput(value float64, address string) *TxOutput {
	var Output TxOutput
	Output.Value = value
	Output.LockWithHash(address)
	return &Output
}
func (output *TxOutput) LockWithHash(address string) {
	decode := base58.Decode(address)
	pubkeyHash := decode[1 : len(decode)-4]
	output.PubkeyHash = pubkeyHash
}

const reward = 12.5

//挖矿交易
func NewCoinbaseTx(address, data string) *Transaction {
	input := TxInput{nil, 0, nil, []byte(data)}
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
	privateKey := wallet.PrivateKey
	pubkey := wallet.Pubkey
	pubkeyHash := HashPubkey(pubkey)
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
			input := TxInput{TXId: []byte(i), Index: uint64(index), Sig: nil, Pubkey: pubkey}
			inputs = append(inputs, &input)
		}
	}
	output := NewOutput(amount, to)
	outputs = append(outputs, output)
	//如果找到的钱比支付的钱多久需要找零
	if calcMoney > amount {
		output1 := NewOutput(calcMoney-amount, from)
		outputs = append(outputs, output1)
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetId()
	blc.SignTransaction(&tx, privateKey)
	return &tx
}
func (tx *Transaction) Sign(privayekey *ecdsa.PrivateKey, prevTxs map[string]*Transaction) bool {
	fmt.Printf("对交易进行签名...\n")
	txCopy := tx.TrimmedCopy()
	//遍历copy的input给每个input添加pubkeyhash
	for i, input := range txCopy.Inputs {
		//查看
		prevTx := prevTxs[string(input.TXId)]
		//将每个input的pubkey字段改为PubkeyHash
		txCopy.Inputs[i].Pubkey = prevTx.Outputs[input.Index].PubkeyHash
		txCopy.SetId()
		//统一清零
		txCopy.Inputs[i].Pubkey = nil
		//需要签名的数据源
		signdatahash := txCopy.TXHash
		//私钥签名
		r, s, err := ecdsa.Sign(rand.Reader, privayekey, signdatahash)
		if err != nil {

			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Inputs[i].Sig = signature
	}

	return true
}

func (tx *Transaction) TrimmedCopy() *Transaction {
	var inputs []*TxInput
	for _, input := range tx.Inputs {

		inputs = append(inputs, &TxInput{input.TXId, input.Index, nil, nil})
	}

	return &Transaction{tx.TXHash, inputs, tx.Outputs}
}
func (tx *Transaction) Verify(prevTxs map[string]*Transaction) bool {
	txCopy := tx.TrimmedCopy()
	//对当前交易的inputs遍历
	for i, input := range tx.Inputs {
		prevTx := prevTxs[string(input.TXId)]
		//与前面的数据源准备一致
		txCopy.Inputs[i].Pubkey = prevTx.Outputs[input.Index].PubkeyHash
		txCopy.SetId()
		txData := txCopy.TXHash
		fmt.Println(txData)
		txCopy.Inputs[i].Pubkey = nil
		pubKey := input.Pubkey
		signature := input.Sig
		fmt.Println("椒盐的", input.Sig)

		//根据signature 切出来r1, s1, 一分为二
		r1 := big.Int{}
		s1 := big.Int{}

		r1Data := signature[:len(signature)/2]
		s1Data := signature[len(signature)/2:]

		r1.SetBytes(r1Data)
		s1.SetBytes(s1Data)

		//切pubkey字节流
		x1 := big.Int{}
		y1 := big.Int{}

		x1Data := pubKey[:len(pubKey)/2]
		y1Data := pubKey[len(pubKey)/2:]

		x1.SetBytes(x1Data)
		y1.SetBytes(y1Data)

		curve := elliptic.P256()
		pubKeyOrigin := ecdsa.PublicKey{curve, &x1, &y1}

		if !ecdsa.Verify(&pubKeyOrigin, txData, &r1, &s1) {
			fmt.Println("校验失败")
			return false
		}
	}
	fmt.Printf("恭喜，校验成功！\n")
	return true
}
func (tx *Transaction) String() string {
	//fmt.Sprintf("打印交易细节...\n")
	//return string("hello")
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.TXHash))

	for i, input := range tx.Inputs {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.TXId))
		lines = append(lines, fmt.Sprintf("       Index:       %d", input.Index))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Sig))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.Pubkey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %f", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubkeyHash))
	}

	return strings.Join(lines, "\n")
}
