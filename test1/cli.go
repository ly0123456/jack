package main

import (
	"fmt"
	"strconv"
)

//定义一个命令行类型
type ClI struct {

}
const Usage = `
Usage:
	createbc ADDRESS "创建区块链"
	print "打印区块链"
	printtx "打印交易"
	balc ADDRESS "获取指定地址余额"
	send <FROM> <TO> <AMOUNT> "转账"
	mine [MINER] [DATA] "挖矿"，默认:  1NVwrN4yZVV3hW1PkXCg38sGcsXMKcYaw7
	createwt "创建钱包地址"
	list "打印钱包中的所有地址"
	status "查看当前待确认交易数量"
`
func checkArgs(cmds []string, count int) bool {
	if len(cmds) != count {
		fmt.Println("参数无效!")
		//os.Exit(1)
		return false
	}

	return true
}
//实现命令行的方法
func (C *ClI)Run(cmds []string)  {
	args:=cmds
	switch args[0] {
	case "createbc":
		if !checkArgs(cmds,2) {return}
			fmt.Println("创建区块链的命令被调用...")
			address:=args[1]
			C.CreateBlockChain(address)

	case "send":
		if !checkArgs(cmds,4) {return}
		from:=args[1]
		to:=args[2]
		amount,_:=strconv.ParseFloat(args[3],64)

		C.Send(from,to,amount)
	case "mine":
		var miner ,data string
		if len(cmds)==3 {
			miner=args[1]
			data=args[2]
		}else {
			miner=""
			data="奖励"
		}
		C.Mine(miner,data)

	case "print":
		if !checkArgs(cmds,1) {return}

		fmt.Println("打印区块被调用.....")
		C.PrintChain()
	case "balc":
		if !checkArgs(cmds,2) {return}
		address:=args[1]
		C.GetBalance(address)
	case "createwt":
		if !checkArgs(cmds,1) {return}
		fmt.Printf("createWallet命令被调用\n")
		C.CreateWallet()
	case "list":
		if !checkArgs(cmds,1) {return}
		fmt.Printf("listAllAddress命令被调用\n")
		C.ListAllAddress()
	case "status":
		if !checkArgs(cmds,1) {return}
		fmt.Printf("status命令被调用\n")
		C.Status()
	case "printtx":
		if !checkArgs(cmds,1) {return}
		C.PrintTx()
	case "help":
		C.Help()
	default:
		fmt.Println("输入参数有误，请参照")

	}
}