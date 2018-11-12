package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

const walletFileName = "wallet.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

//实例化钱包
func NewWallets() *Wallets {
	var wallets *Wallets
	wallets.Wallets = make(map[string]*Wallet)
	wallets.LoadFromFile()
	return wallets
}
func (ws *Wallets) CreateWallets() {
	//new一个私钥对
	wallet := NewWallet()
	//生成钱包地址
	address := wallet.GetAddress()

	ws.Wallets[address] = wallet
	ws.SaveToFile()

}
//本地化钱包
func (ws *Wallets) SaveToFile() bool {
	var buffer bytes.Buffer
	//************如果有interface 一定要注册一下*************
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&ws)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	//存入本地
	err = ioutil.WriteFile(walletFileName, buffer.Bytes(), 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
		return false
	}
	return true
}

//加载
func (ws *Wallets) LoadFromFile() bool {
	if !IsExist(walletFileName) {
		return false
	}
	dataInfo, e := ioutil.ReadFile(walletFileName)
	if e != nil {
		fmt.Println(e)
		os.Exit(-1)
	}
	var wallets *Wallets
	//************如果有interface 一定要注册一下*************
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(dataInfo))
	e = decoder.Decode(&wallets)
	if e != nil {
		fmt.Println(e)
		os.Exit(-1)
	}
	ws.Wallets = wallets.Wallets
	return true
}

//获取钱包地址
func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}
