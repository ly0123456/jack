package main

import (
	"os"
	"fmt"
)
//判断文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
func Welcome() {
	fmt.Printf("\n====================================================================\n")
	fmt.Printf("                   欢迎来到传智播客!(英雄联盟)\n")
	fmt.Printf("====================================================================\n")
}
