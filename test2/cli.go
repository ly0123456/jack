package main

import (
	"fmt"
	"strconv"
)

const Usage = `
Usage:
	createBlockChain ADDRESS "创建区块链"
	printChain "打印区块链"
	getBalance ADDRESS "获取指定地址余额"
	send FROM TO AMOUNT "转账"
	createWallet "创建钱包地址"
	listAllAddress "打印钱包中的所有地址"
	mine ADDRESS DATA "挖矿"
	status "查看交易状态"
`

//定义一个CLI，里面包含BlockChain，所有细节工作交给bc，命令的解析工作交给CLI
type CLI struct {
	//bc *BlockChain
}

func checkArgs(cmds []string, count int) bool {
	if len(cmds) != count {
		fmt.Println("参数无效!")
		return false
		//os.Exit(1)
	}
	return true
}

//定义一个run函数，负责接收命令行的数据，然后根据命令进行解析，并完成最终的调用
func (cli *CLI) Run(cmds []string) {
	//args := os.Args

	cmd := cmds[0]

	switch cmd {

	case "createBlockChain":
		fmt.Printf("创建区块链命令被调用!\n")
		checkArgs(cmds, 2)
		address := cmds[1]
		cli.CreateBlockChain(address)

	case "printChain":
		fmt.Printf("打印区块命令被调用\n")
		checkArgs(cmds, 1)
		cli.PrintChain()

	case "getBalance":
		fmt.Printf("获取余额命令被调用\n")
		checkArgs(cmds, 2)
		address := cmds[1]
		cli.GetBalance(address)

	case "send":
		fmt.Printf("转账send命令被调用\n")
		checkArgs(cmds, 4)

		from := cmds[1]
		to := cmds[2]
		amount, _ := strconv.ParseFloat(cmds[3], 64) //string
		cli.Send(from, to, amount)

	case "createWallet":
		fmt.Printf("createWallet命令被调用\n")
		checkArgs(cmds, 1)

		cli.CreateWallet()

	case "listAllAddress":
		fmt.Printf("listAllAddress命令被调用\n")
		checkArgs(cmds, 1)
		cli.ListAllAddress()

	case "mine":
		fmt.Printf("mine命令被调用\n")

		var miner string
		var data string

		if len(cmds) != 3 {
			miner = "1NVwrN4yZVV3hW1PkXCg38sGcsXMKcYaw7"
			data = "helloworld"
		} else {
			miner = cmds[1]
			data = cmds[2]
		}

		cli.Mine(miner, data)
	case "status":
		cli.Status()

	default:
		fmt.Printf("无效的命令，请检查!\n")
		cli.Help()
	}
}
