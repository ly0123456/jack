package main

import (
	"os"
	"fmt"
	"strconv"
)

//定义一个命令行类型
type ClI struct {

}
const Usage = `
Usage:
		./blockchain createBlockChain ADDRESS "创建区块链"
	./blockchain addBlock DATA   "添加数据到区块链"
	./blockchain printChain "打印区块链"
	./blockchian getBalance ADDRESS "获取指定地址余额"
	./blockchian send FROM TO AMOUNT MINER DATA  "转账"
	./blockchian createWallet "创建钱包地址"
	./blockchian listAllAddress "打印钱包中的所有地址"
`
func checkArgs(count int) {
	if len(os.Args) != count {
		fmt.Println("参数无效!")
		os.Exit(1)
	}
}
//实现命令行的方法
func (C *ClI)Run()  {
	args := os.Args
	if len(args)<2{
		fmt.Println("输入的参数有误，请参照")
		fmt.Println(Usage)
		os.Exit(-1)
	}
	switch args[1] {
	case "createBlockChain":
		checkArgs(3)
			fmt.Println("创建区块链的命令被调用...")
			address:=args[2]
			C.CreateBlockChain(address)

	case "send":
		checkArgs(7)
		from:=args[2]
		to:=args[3]
		amount,_:=strconv.ParseFloat(args[4],64)
		miner:=args[5]
		data:=args[6]
		C.Send(from,to,amount,miner,data)
	case "printChain":
		fmt.Println("打印区块被调用.....")
		C.PrintChain()
	case "getBalance":
		checkArgs(3)
		address:=args[2]
		C.GetBalance(address)
	case "createWallet":
		fmt.Printf("createWallet命令被调用\n")
		checkArgs(2)

		C.CreateWallet()
	case "listAllAddress":
		fmt.Printf("listAllAddress命令被调用\n")
		checkArgs(2)
		C.ListAllAddress()
	case "help":
		C.Help()
	default:
		fmt.Println("输入参数有误，请参照")

	}
}