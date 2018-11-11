package main

import (
	"./bolt"
	"fmt"
	"log"
	"os"
)

//定义一个数据库的名字
const BlockChainDB = "BlockChainDB.db"

//定义一个桶的名字
const BlcokBucket = "BlcokBucket"

//定义一个最后区块的hashKey
const LastHashKey = "LastHashKey"

//定义一个区块链的类型
type BlockChain struct {
	Db       *bolt.DB //数据库的句柄
	LastHash []byte   //最后一个区块的hash用于索引所有的
}

//创建一个新的区块链
func CreateBlockChain(address string) *BlockChain {
	//定义一个当前区块的hash
	var LastHash []byte
	if IsExist() {
		fmt.Println("数据库已经存在")
		os.Exit(-1)
	}
	//打开一个数据库，如果没有就创建
	db, e := bolt.Open(BlockChainDB, 0600, nil)
	if e != nil {
		fmt.Println("数据库打开失败")
		os.Exit(-1)
	}
	//向数据中添加数据
	db.Update(func(tx *bolt.Tx) error {
		buckt := tx.Bucket([]byte(BlcokBucket))
		//如果桶不存在就需要创建一个新的桶
		if buckt == nil {
			buckt, e = tx.CreateBucket([]byte(BlcokBucket))
			if e != nil {
				log.Panic(e)
			}
			tx := NewCoinbaseTx(address, "Gensis Block..... ")
			newBlock := NewBlock([]*Transaction{tx}, []byte{})
			//将新区块存入数据库，key就是当前区块的hash
			e = buckt.Put(newBlock.Hash, newBlock.Encode())
			if e != nil {
				log.Println(e)
			}
			//把最后区块的hash存入数据库
			e = buckt.Put([]byte(LastHashKey), newBlock.Hash)
			if e != nil {
				log.Println(e)
			}
			//满足blc的结构
			LastHash = newBlock.Hash
		}
		return nil
	})
	return &BlockChain{db, LastHash}
}

//实例化已有区块链
func NewBlockChain() *BlockChain {
	//定义一个当前区块的hash
	var LastHash []byte
	if !IsExist() {
		fmt.Println("数据库不存在，请检查")
		os.Exit(-1)
	}
	//打开一个数据库，如果没有就创建
	db, e := bolt.Open(BlockChainDB, 0600, nil)
	if e != nil {
		fmt.Println("数据库打开失败")
		os.Exit(-1)
	}
	//向数据中查看数据
	db.View(func(tx *bolt.Tx) error {
		buckt := tx.Bucket([]byte(BlcokBucket))
		//如果桶不存在就需要创建一个新的桶
		if buckt == nil {
			fmt.Println("有问题请检查")
			os.Exit(-1)
		}
		LastHash = buckt.Get([]byte(LastHashKey))

		//满足blc的结构

		return nil
	})

	return &BlockChain{db, LastHash}
}
func (blc *BlockChain) AddBlock(txs []*Transaction) {
	//找到前区块的hash
	lastHash := blc.LastHash
	//创建新的区块链
	newBlock := NewBlock(txs, lastHash)
	newBlock.Height++
	//将新区快存入数据库
	blc.Db.Update(func(tx *bolt.Tx) error {
		buckt := tx.Bucket([]byte(BlcokBucket))
		//如果桶不存在就需要创建一个新的桶
		if buckt == nil {
			fmt.Println("有问题请检查")
			os.Exit(-1)
		}
		//将新区块存入数据库，key就是当前区块的hash
		e := buckt.Put(newBlock.Hash, newBlock.Encode())
		if e != nil {
			log.Println(e)
		}
		//把最后区块的hash存入数据库
		e = buckt.Put([]byte(LastHashKey), newBlock.Hash)
		if e != nil {
			log.Println(e)
		}
		return nil
	})
}

//定义一个迭代器
type Iterator struct {
	Db        *bolt.DB
	HashPoint []byte // 去开指针总是指向最后一个区块
}

func (blc *BlockChain) NewIterator() *Iterator {
	db := blc.Db
	HashPoint := blc.LastHash
	return &Iterator{db, HashPoint}
}

//迭代器的方法
func (it *Iterator) Next() *Block {
	var block *Block
	it.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlcokBucket))
		if bucket == nil {
			fmt.Println("有问题请检查")
			os.Exit(-1)
		}
		//根据hash指针找到最后一个区块
		databytes := bucket.Get(it.HashPoint)
		//反序列化得到最后一个区块
		block = Decode(databytes)
		//找到最后一个区块后把hash指针前移
		it.HashPoint = block.PrevBlockHash
		return nil
	})

	return block
}

//实现区块链的打印方法
func (blc *BlockChain) PrintBlock() {
	//先new一个迭代器
	iterator := blc.NewIterator()
	for {
		//使用迭代器的方法取出区块
		block := iterator.Next()
		//打印区块

		fmt.Printf("================当前区块高度%d=====================\n", block.Height)
		fmt.Printf("当前版本号:%d\n", block.Vision)
		fmt.Printf("前区块hash:%x\n", block.PrevBlockHash)
		fmt.Printf("当前区块hash:%x\n", block.Hash)
		fmt.Printf("当前区块的data：%s\n", block.Data)
		fmt.Printf("当前难度：%d\n", block.Diff)
		fmt.Printf("时间戳：%d\n", block.Timestamp)
		fmt.Printf("当前随机值：%d\n\n", block.Nonce)

		if len(block.PrevBlockHash) == 0 {
			fmt.Println("打印完毕")
			break
		}
	}
}
func IsExist() bool {
	_, err := os.Stat(BlockChainDB)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
func (blc *BlockChain) GetUTXOs(address string) []*TxOutput {
	var UTXOs []*TxOutput
	spendUtxos := make(map[string][]uint64)

	it := blc.NewIterator()

	for {
		block := it.Next()
		for _, tx := range block.Txs {
		OUT:
			for i, UTXO := range tx.Outputs {

				if UTXO.ScriptPubKey == address {
					fmt.Println("zhaodao l ")
					if len(spendUtxos[string(tx.TXHash)]) != 0 {
						for _, index := range spendUtxos[string(tx.TXHash)] {
							if index == uint64(i) {
								continue OUT
							}
						}
					}

					UTXOs = append(UTXOs, UTXO)
				}
			}
			for _, input := range tx.Inputs {
				if input.Sig == address {
					spendUtxos[string(input.TXId)] = append(spendUtxos[string(input.TXId)], input.Index)
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXOs
}

func (blc *BlockChain) GetBalance(address string) {
	var amount float64
	utxOs := blc.GetUTXOs(address)
	for _, utxo := range utxOs {
		amount += utxo.Value
	}
	fmt.Println(amount)
}
func (blc *BlockChain) FindNeedUtxos(from string, amount uint64) (map[string][]uint64, uint64) {
	//TODO

	return nil, 10
}
