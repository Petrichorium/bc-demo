package blockchain

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"time"
)

type Blockchain struct {
	lasthash Hash //最后一个区块的hash
	//blocks   map[Hash]*Block //全部区块信息，由区块hash作为 key 来检索
	db *leveldb.DB //leveldb的链接
}

// NewBlockchain 构造方法
func NewBlockchain(db *leveldb.DB) *Blockchain {
	//初始化区块链
	bc := &Blockchain{
		//blocks: map[Hash]*Block{},
		db: db,
	}
	//读取最新的区块hash
	data, err := bc.db.Get([]byte("lasthash"), nil)
	if err == nil { //读取到lasthash
		bc.lasthash = Hash(data)
	}
	return bc
}

// AddBlock 添加区块
// 提供区块的数据，目前事字符串
func (bc *Blockchain) AddBlock(txs string) *Blockchain {
	//实例化区块
	b := NewBlock(bc.lasthash, txs)

	//将区块链最hash设置为当前区块的hash
	bc.lasthash = b.hashCurr

	//将区块加入到链的存储结构中
	//bc.blocks[b.hashCurr] = b
	if bs, err := BlockSerialize(*b); err != nil {
		log.Fatal("block can not be serialized.")
	} else if err = bc.db.Put([]byte("b_"+b.hashCurr), bs, nil); err != nil {
		log.Fatal("block can not be saved.")
	}

	//将最新的区块hash存储到数据库中
	if err := bc.db.Put([]byte("lasthash"), []byte(b.hashCurr), nil); err != nil {
		log.Fatal("lastHash can not be saved.")
	}

	return bc
}

// AddGensisBlock 添加创世区块
func (bc *Blockchain) AddGensisBlock() *Blockchain {
	//校验是否可以添加创世区块
	if bc.lasthash != "" {
		return bc
	}

	//将区块加入到区块链中
	//创世区块只有 txs 是特殊的
	bc.AddBlock("The Gensis Block.")

	return bc
}

// GetBlock 通过hash获取区块
func (bc *Blockchain) GetBlock(hash Hash) (*Block, error) {
	//从数据库中读取对应的区块
	data, err := bc.db.Get([]byte("b_"+hash), nil)
	if err != nil {
		return nil, err
	}
	//反序列化
	b, err := BlockUnSerialize(data)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// Iterate 迭代展示区块的方法
func (bc *Blockchain) Iterate() {
	//从最后一个区块的hash开始迭代
	//for hash := bc.lasthash; hash != ""; {
	//	b := bc.blocks[hash]
	//	fmt.Println("HashCurr", b.hashCurr)
	//	fmt.Println("Txs", b.txs)
	//	fmt.Println("Time", b.header.time.Format(time.UnixDate))
	//	fmt.Println("HashPrevBlock", b.header.hashPrevBlock)
	//	fmt.Println("")
	//	hash = b.header.hashPrevBlock
	//}

	//从最后一个区块的hash开始迭代
	for hash := bc.lasthash; hash != ""; {
		//得到区块
		b, err := bc.GetBlock(hash)
		if err != nil {
			log.Fatalf("Block <%s> is not exsits.", hash)
		}
		fmt.Println("HashCurr:", b.hashCurr)
		fmt.Println("Txs:", b.txs)
		fmt.Println("Time:", b.header.time.Format(time.UnixDate))
		fmt.Println("HashPrevBlock:", b.header.hashPrevBlock)
		fmt.Println("")
		hash = b.header.hashPrevBlock
	}
}
func (bc *Blockchain) Clear() {
	// 数据库中全部区块链的key全部删除
	bc.db.Delete([]byte("lastHash"), nil)
	// 迭代删除，全部的b_的key
	iter := bc.db.NewIterator(util.BytesPrefix([]byte("b_")), nil)
	for iter.Next() {
		bc.db.Delete(iter.Key(), nil)
	}
	iter.Release()
	//清空bc对象
	bc.lasthash = ""
}
