package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

const walletFileName = "wallet.dat"

type Wallets struct {
	WallteMap map[string]*Wallte
}

//实例化一个钱包
func NewWallets() *Wallets {
	var wallets Wallets
	wallets.WallteMap = make(map[string]*Wallte)
	//将本地的钱包加载到内存中
	wallets.LoadFromFile()
	return &wallets
}

//创建一个钱包地址
func (w *Wallets) CreateAddress() string {
	wallre := NewWallre()
	address := wallre.CreateAddress()
	w.WallteMap[address] = wallre
	if !w.SaveToFile() {
		return ""
	}
	return address
}
func (w *Wallets) SaveToFile() bool {
	//把wallets的数据用转码用于存入数据库
	var buffer bytes.Buffer
	//gob: type not registered for interface: elliptic.p256Curve
	//1. gob编码的结构里面如果涉及到了interface类型的数据，要对gob进行注册
	//interface注册函数
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&w)
	if err != nil {
		fmt.Println(err)
		return false
	}
	//用ioutl将数据存到本地
	err = ioutil.WriteFile(walletFileName, buffer.Bytes(), 0600)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//找到钱包中的地址
func (w *Wallets) GetAdderss() []string {
	var addresses []string
	for address := range w.WallteMap {
		addresses = append(addresses, address)
	}
	return addresses
}

//把本地的钱包信息加载到内存中
func (w *Wallets) LoadFromFile() bool {
	file, e := ioutil.ReadFile(walletFileName)
	if e != nil {
		fmt.Println(e)
		return false
	}
	var wallets Wallets
	//使用gob解码需要注册
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(file))
	e = decoder.Decode(&wallets)
	if e != nil {
		fmt.Println(e)
		return false
	}
	w.WallteMap = wallets.WallteMap
	return true
}
