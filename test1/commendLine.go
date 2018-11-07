package main

import (
	"fmt"
	"os"
)

func (cli *ClI)Send(from ,to string,amount float64 ){
	fmt.Printf("%s 向 %s 转账 %f, 由 %s , data : %s\n", from, to, amount)
	if !IsValidAddress(from){
		fmt.Println("from地址有问题请检查")
		return
	}
	if !IsValidAddress(to){
		fmt.Println("to地址有问题请检查")
		return
	}
	//if !IsValidAddress(miner){
	//	fmt.Println("miner地址有问题请检查")
	//	return
	//}
	blockChain := NewBlockChain()
	//创建挖矿

	//创建普通交易
	transaction := NewTransaction(from, to, amount, blockChain)
	if transaction!=nil {
		gTx=append(gTx, transaction)
	}

}
func (cli *ClI) CreateBlockChain(address string) {
	if !IsValidAddress(address){
		fmt.Println("CreateBlockChain地址有问题请检查")
		return
	}
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
	if !IsValidAddress(address){
		fmt.Println("GetBalance地址有问题请检查")
		return
	}
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

	fmt.Printf("生成的新地址%v\n" ,address)
}
func (cli *ClI)ListAllAddress()  {
	wallets := NewWallets()
	addresses := wallets.GetAddresses()
	for  _,address:=range addresses  {
		fmt.Println(address)
	}
}
func (cli *ClI)Help()  {
	fmt.Println(Usage)
}
func (cli *ClI)Mine(miner ,data string){
	if !IsValidAddress(miner){
		fmt.Println("miner地址有问题请检查")
		return
	}
	blockChain := NewBlockChain()
	//创建挖矿
	coinbase:=NewCoinbaseTx(miner,data)
	gTx =append(gTx, coinbase)

	//添加到区块中ll
	blockChain.AddBlock(gTx)
	gTx=[]*Transaction{}
	fmt.Println("添加区块成功")
}
func (cli *ClI)Status(){
	i:=len(gTx)
fmt.Printf("待确认交易数量%d\n",i)
}
func (cli *ClI)PrintTx()  {
	bc:=NewBlockChain()
	it:=bc.NewIterator()
	for   {
		block:=it.Next()
		fmt.Printf("++++++++++++++++++++++++++++++++++\n")

		for _,tx:=range block.Txs {
			fmt.Println(tx)
		}
		if len(block.PrevBlockHash)==0 {
			break
		}
	}

}