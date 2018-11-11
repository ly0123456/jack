package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
)

//定义一个工作量证明的类
type ProofOfWork struct {
	Block *Block
	//pow目标值
	target *big.Int
}

//实例化工作量证明
func NewPow(block *Block) *ProofOfWork {
	//初始化一个大数0000000000000000000000000000000000000000000000000000001
	target := big.NewInt(1)
	//实现我们要求的目标值
	target.Lsh(target, 256-Diff) //0000100000000000000000000000000
	return &ProofOfWork{block, target}
}

//实现pow找满足target的函数，返回hash和nonce
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var hash [32]byte
	var nonce uint64
	for {
		//将区块的数据拼接成[]byte
		hashData := pow.preHashData(nonce)
		//算出当前的hash
		hash = sha256.Sum256(hashData)
		//定义一个大数，用于存储算出来的hash用于与target比较
		var hashInt big.Int
		hashInt.SetBytes(hash[:])
		//如果找的hash比target小就说明找了我们规定的hash成功然后退出
		if pow.target.Cmp(&hashInt) == 1 {
			fmt.Println("找到了")
			break
		}
		//如果没有找到就nonce++继续找
		nonce++
	}

	return hash[:], nonce
}

//Vision uint64//版本号
//Height uint64//区块高度
//PrevBlockHash []byte//前区块hash
//MorkHash []byte //默克尔根
//Timestamp uint64 //时间戳
//Data []byte //交易数据
//Hash  []byte //本区块hash
//Diff  uint64 //当前难度
//Nonce uint64 //随机值
func (pow *ProofOfWork) preHashData(nonce uint64) []byte {
	data := bytes.Join([][]byte{
		Num2Byte(pow.Block.Vision),
		Num2Byte(pow.Block.Height),
		pow.Block.PrevBlockHash,
		pow.Block.MorkHash,
		Num2Byte(pow.Block.Timestamp),
		Num2Byte(pow.Block.Diff),
		Num2Byte(nonce),
	}, []byte{})
	return data
}

//用于处理数字转[]byte
func Num2Byte(num uint64) []byte {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, &num) //binary.BigEndian大端存储
	return buffer.Bytes()
}
//实现一个pow的校验方法Isvaild
func (pow *ProofOfWork)Isvaild()bool {

	hashData := pow.preHashData(pow.Block.Nonce)
	hash := sha256.Sum256(hashData)
	var hashInt big.Int
	hashInt.SetBytes(hash[:])
	if pow.target.Cmp(&hashInt)==-1{
		return  false
	}
	return true
}