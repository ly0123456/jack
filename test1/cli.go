package main

import (
	"os"
	"fmt"
)

//定义一个命令行类型
type ClI struct {

}
const Usage = `
	./blockchain createBlockChain "创建区块链"
	./blockchain addBlock DATA   "添加数据到区块链"
	./blockchain printChain "打印区块链"
`
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
		fmt.Println("创建区块链的命令被调用...")
		blockChain := CreateBlockChain()
		defer blockChain.Db.Close()
	case "addBlock":
		fmt.Println("添加区块被调用....")
		if len(args)==3 {
			blockChain := NewBlockChain()
			blockChain.AddBlock(args[2])
			defer blockChain.Db.Close()
		}
	case "printChain":
		fmt.Println("打印区块被调用.....")
		blockChain := NewBlockChain()
		blockChain.PrintBlock()
		defer blockChain.Db.Close()
	default:
		fmt.Println("输入参数有误，请参照")
		fmt.Println(Usage)

	}
}