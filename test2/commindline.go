package main

import (
	"fmt"
	"os"
)

//创建新的blockchain
func (cli *CLI) CreateBlockChain(adderss string) {
	blocChain := CreateBlocChain(adderss)
	defer blocChain.Db.Close()
}

//打印区块链
func (cli *CLI) PrintChain() {
	chain := NewBlockChain()
	chain.PrintBlockChain()
	chain.Db.Close()
}

//查询余额
func (cli *CLI) GetBalance(address string) {}

//转账
func (cli *CLI) Send(from, to string, amount float64) {}

//创建新的钱包地址
func (cli *CLI) CreateWallet() {
	wallets := NewWallets()
	address := wallets.CreateAddress()
	if address == "" {
		fmt.Printf("创建地址失败\n")
		os.Exit(1)
	}
	fmt.Println("你的新地址是：", address)
}

//查看钱包
func (cli *CLI) ListAllAddress() {
	wallets := NewWallets()
	addersses := wallets.GetAdderss()
	for i, k := range addersses {
		fmt.Printf("第%d的地址：%v\n", i+1, k)
	}
}

//挖矿
func (cli CLI) Mine(miner, data string) {}

//查看交易状态
func (cli *CLI) Status() {}

//帮助信息
func (cli *CLI) Help() {
	fmt.Printf(Usage)
}
