package main

import (
	"fmt"
	"github.com/bolt-master"
	"log"
	"os"
	"time"
)

const BlockChainDB = "BlockChian.db"
const BlockChainBuckey = "BlockChainBucket"
const lastHashKey = "lastHashKey"

type BlockChain struct {
	Db       *bolt.DB
	LastHash []byte
}

func CreateBlockChain() *BlockChain {
	var LastHash []byte
	if isDbExist(){
		fmt.Println("数据库已存在")
		os.Exit(-1)
	}
	db, e := bolt.Open(BlockChainDB, 0600, nil)
	if e != nil {
		fmt.Println("数据库打开失败")
		os.Exit(-1)
	}
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlockChainBuckey))
		if bucket == nil {
			bucket, err := tx.CreateBucket([]byte(BlockChainBuckey))
			if err != nil {
				log.Panic(err)
			}
			newBlock := NewBlock("Genesis data...", []byte{})
			bucket.Put(newBlock.Hash, newBlock.EncSec())
			bucket.Put([]byte(lastHashKey), newBlock.Hash)
			LastHash = newBlock.Hash
		}
		return nil
	})
	return &BlockChain{db, LastHash}
}

func NewBlockChain() *BlockChain {
	var lastHash []byte
	if !isDbExist(){
		fmt.Println("数据库不存在")
		os.Exit(-1)
	}
	db, e := bolt.Open(BlockChainDB, 0600, nil)
	if e != nil {
		fmt.Println("数据库打开失败")
		os.Exit(-1)
	}
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlockChainBuckey))
		if bucket == nil {
			fmt.Printf("获取区块链实例时bucket不应为空!")
			os.Exit(1)
		}
		lastHash = bucket.Get([]byte(lastHashKey))
		return nil
	})
	return &BlockChain{db, lastHash}
}
func (blc *BlockChain) AddBlockChain(data string) {
	blc.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlockChainBuckey))
		if bucket == nil {
			fmt.Printf("获取区块链实例时bucket不应为空!")
			os.Exit(1)
		}
		newBlock := NewBlock(data, blc.LastHash)
		bucket.Put(newBlock.Hash, newBlock.EncSec())
		bucket.Put([]byte(lastHashKey), newBlock.Hash)
		blc.LastHash=newBlock.Hash
		return nil
	})
}

type Iterator struct {
	Db          *bolt.DB //来自于区块链
	currentHash []byte   //随着移动改变
}

func (bc *BlockChain) NewIterator() *Iterator {
	return &Iterator{Db: bc.Db, currentHash: bc.LastHash}
}
func(it *Iterator)Next()*Block{
	var block *Block
	 it.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlockChainBuckey))
		if bucket == nil {
			fmt.Printf("获取区块链实例时bucket不应为空!")
			os.Exit(1)
		}
		blockTemp := bucket.Get(it.currentHash)
		block=Decode(blockTemp)
		it.currentHash=block.PrevBlockHash
		return nil
	})
		return block
}

func (blc *BlockChain)PrintBlockChain()  {
	it:=blc.NewIterator()

	for ; ;  {
		block := it.Next()

		fmt.Printf("===============================\n")
		fmt.Printf("PrevBlockHash :%x\n", block.PrevBlockHash)
		fmt.Printf("MerkeRoot :%x\n", block.MerkelRoot)
		timeFormat := time.Unix(int64(block.TimeSTamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp: %s\n", timeFormat)
		//fmt.Printf("TimeStamp :%d\n", block.TimeStamp)
		fmt.Printf("Difficulty :%d\n", block.Diff)
		fmt.Printf("Nonce :%d\n", block.Nonce)
		fmt.Printf("Hash :%x\n", block.Hash)
		fmt.Printf("Data :%s\n", block.Data)
		pow := NewPow(block)
		fmt.Printf("IsValid : %v\n\n", pow.IsVild())



		if len(block.PrevBlockHash)==0{
			fmt.Println("打印完毕")
			break
		}
	}
}

//判断区块链文件是否存在
func isDbExist() bool {
	if _, err := os.Stat(BlockChainDB); os.IsNotExist(err) {
		return false
	}

	return true
}








