package main

import (
	"fmt"
	"os"
)

func (cli *ClI)Send(from ,to string,amount float64,miner ,data string ){
	fmt.Printf("%s 向 %s 转账 %f, 由 %s , data : %s\n", from, to, amount, miner, data)

	blockChain := NewBlockChain()
	//创建挖矿
	coinbaseTx := NewCoinbaseTx(miner,data)
	txs:= []*Transaction{coinbaseTx}
	//创建普通交易
	transaction := NewTransaction(from, to, amount, blockChain)
	if transaction!=nil {
		txs=append(txs, transaction)
	}
	//添加到区块中ll
	blockChain.AddBlock(txs)
	fmt.Println("添加区块成功")
}
func (cli *ClI) CreateBlockChain(address string) {
	//3. 调用真正的添加区块函数
	bc := CreateBlockChain(address)
	defer bc.Db.Close()
}

func (cli *ClI) PrintChain() {
	bc := NewBlockChain()
	defer bc.Db.Close()
	bc.PrintBlock()
}

func (cli *ClI) GetBalance(address string) {
	bc := NewBlockChain()
	bc.GetBalance(address)
}
func (cli *ClI)CreateWallet() {
	wallets := NewWallets()
	address := wallets.CreateWallets()
	if address == "" {
		fmt.Printf("创建地址失败\n")
		os.Exit(1)
	}

	fmt.Printf("生成的新地址%x\n" ,address)
}
func (cli *ClI)ListAllAddress()  {
	wallets := NewWallets()
	addresses := wallets.GetAddresses()
	for  i,address:=range addresses  {
		fmt.Printf("第%d的地址为%x\n",i+1,address)
	}
}
func (cli *ClI)Help()  {
	fmt.Println(Usage)
}