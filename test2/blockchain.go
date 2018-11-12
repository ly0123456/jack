package main

import (
	"fmt"
	"github.com/bolt-master"
	"log"
	"os"
	"time"
)

const reward = 12.5
const blockChainName = "blockChain.db"
const blockBucket = "blockBucket"
const lastHashKey = "lastHashKey"

type BlockChain struct {
	//数据库句柄
	Db *bolt.DB
	//最后一个区块的hash用于查询
	LastHash []byte
}

func CreateBlocChain(address string) *BlockChain {
	var lashHash []byte
	//判断数据库文件是否存在
	if Isexitis(blockChainName) {
		fmt.Println("数据库已存在，请检查")
		return nil
	}
	db, e := bolt.Open(blockChainName, 0600, nil)
	if e != nil {
		fmt.Println("打开数据库失败", e)
		os.Exit(-1)
	}
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			bucket, e = tx.CreateBucket([]byte(blockBucket))
			if e != nil {
				fmt.Println("创建桶失败，请检查", e)
				os.Exit(-1)
			}
			//创建建挖矿交易
			coinbases := NewCoinbaseTx(address, genesisInfo)
			block := NewBlock([]*Transaction{coinbases}, nil)
			//用当前hash作为key 当前区块的字节流作为value存入数据库
			e = bucket.Put(block.Hash, block.Encode())
			if e != nil {
				log.Panic(e)
			}
			//吧当前hash存入数据库用于查询
			e = bucket.Put([]byte(lastHashKey), block.Hash)
			if e != nil {
				log.Panic(e)
			}
			lashHash = block.Hash
		}

		return nil
	})

	return &BlockChain{db, lashHash}
}

//生成一个区块链实例
func NewBlockChain() *BlockChain {
	var lashHash []byte
	//判断数据库文件是否存在
	if Isexitis(blockChainName) == false {
		fmt.Println("数据库不存在，请检查")
		os.Exit(-1)
	}
	db, e := bolt.Open(blockChainName, 0600, nil)
	if e != nil {
		fmt.Println("打开数据库失败", e)
		os.Exit(-1)
	}
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			fmt.Println("桶不存在，请检查....")
			os.Exit(-1)
		}
		//lasthash就是数据库中最后存的hash
		lashHash = bucket.Get([]byte(lastHashKey))
		return nil
	})

	return &BlockChain{db, lashHash}
}
func (blc *BlockChain) AddBlock(txs []*Transaction) {
	blc.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			fmt.Println("桶不存在，请检查....")
			os.Exit(-1)
		}
		lasthash := blc.LastHash
		newblcok := NewBlock(txs, lasthash)
		e := bucket.Put(newblcok.Hash, newblcok.Encode())
		if e != nil {
			log.Panic(e)
		}
		e = bucket.Put([]byte(lastHashKey), newblcok.Hash)
		if e != nil {
			log.Panic(e)
		}
		blc.LastHash = newblcok.Hash
		return nil
	})
}

//定义一个迭代器
type Iterator struct {
	Db          *bolt.DB //来自于区块链
	currentHash []byte   //随着移动改变
}

//实例化一个迭代器
func (blc *BlockChain) NewItertor() *Iterator {
	return &Iterator{blc.Db, blc.LastHash}
}

//迭代器的迭代方法
func (it *Iterator) Next() *Block {
	var block *Block
	it.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			fmt.Println("桶不存在，请检查....")
			os.Exit(-1)
		}
		tempdata := bucket.Get(it.currentHash)
		block = Decode(tempdata)
		//更新区块指针，让他总是指向最后一个
		it.currentHash = block.PrevBlockHash
		return nil
	})
	return block
}

//实现打印区块链的方法
func (blc *BlockChain) PrintBlockChain() {
	itertor := blc.NewItertor()
	for {
		block := itertor.Next()

		fmt.Printf("===============================\n")
		fmt.Printf("Version :%d\n", block.Version)
		fmt.Printf("PrevBlockHash :%x\n", block.PrevBlockHash)
		fmt.Printf("MerkeRoot :%x\n", block.MerkelRoot)
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp: %s\n", timeFormat)
		//fmt.Printf("TimeStamp :%d\n", block.TimeStamp)
		fmt.Printf("Difficulty :%d\n", block.Difficulty)
		fmt.Printf("Nonce :%d\n", block.Nonce)
		fmt.Printf("Hash :%x\n", block.Hash)
		fmt.Printf("Data :%s\n", block.Transactions[0].TxIntPuts[0].Pubkey)
		pow := NewPow(*block)
		fmt.Printf("IsValid : %v\n\n", pow.IsVaild(block.Nonce))
		if len(block.PrevBlockHash) == 0 {
			fmt.Println("打印完毕")
			break
		}
	}
}
