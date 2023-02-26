package blockchain

import (
	"fmt"
	"github.com/cn-org-Pretichor/bc-demo/block"
	"github.com/cn-org-Pretichor/bc-demo/pow"
	"github.com/cn-org-Pretichor/bc-demo/tx"
	"github.com/cn-org-Pretichor/bc-demo/wallet"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"time"
)

type Blockchain struct {
	lasthash block.Hash // 最后一个区块的hash
	// blocks   map[Hash]*Block //全部区块信息，由区块hash作为 key 来检索
	db *leveldb.DB // leveldb的链接

	UTXOCache *UTXOCache // UTXO 缓存对象，完成UTXO缓存的相关操作
}

// NewBlockchain 构造方法
func NewBlockchain(db *leveldb.DB) *Blockchain {
	// 初始化区块链
	bc := &Blockchain{
		// blocks: map[Hash]*Block{},
		db:        db,
		UTXOCache: NewUTXOCache(db),
	}
	// 读取最新的区块hash
	data, err := bc.db.Get([]byte("lasthash"), nil)
	if err == nil { // 读取到lasthash
		bc.lasthash = block.Hash(data)
	}
	return bc
}

// AddBlock 添加区块
// 提供区块的数据，目前事字符串
// address 添加该区块的地址
func (bc *Blockchain) AddBlock(address wallet.Address) *Blockchain {
	// 实例化区块
	b := block.NewBlock(bc.lasthash)

	// 为区块增加交易，任何区块都有 coinbase 交易
	cbtx := tx.NewCoinBaseTX(address)
	// 将交易加入到区块中
	b.AddTx(cbtx)

	// 更新交易对应的UTXO
	bc.UTXOCache.Update(cbtx)

	// TODO-petrichor 为了获取TXS，需要进行反序列化，以便测试用
	bc.TxsInit()
	// 处理区块常规交易
	// 从交易缓存队列中，获取交易
	// 假设每个区块仅可以存储3个交易，去掉coinbase交易，此处从交易缓存中获取两个交易
	// 实操时，通过交易size在控制，一个区块1M，一笔交易大约为250Byte，装满为止
	for i, l := 0, len(TXS); i < l && i < 2; i++ {
		t := TXS[0]
		// TODO-Petrichor 删除交易缓存列表中被取出的交易
		if len(TXS) > 1 {
			TXS = TXS[1:]
		} else {
			// TODO-Petrichor 清空交易缓存列表
			TXS = []*tx.TX{}
		}
		// 将交易加入到区块中
		b.AddTx(t)
		// 更新交易对应的UTXO
		bc.UTXOCache.Update(t)
	}

	// 对区块做POW，工作量证明
	// pow对象
	p := pow.NewPow(b)
	// 开始证明
	nonce, hash := p.Proof()
	if nonce == 0 || hash == "" {
		log.Fatal("block hashcash proof faild!")
	}
	b.SetNonce(nonce)
	b.SetHashCurr(hash)

	// 将区块链最hash设置为当前区块的hash
	bc.lasthash = b.GetHashCurr()

	// 将区块加入到链的存储结构中
	// bc.blocks[b.hashCurr] = b
	if bs, err := block.BlockSerialize(*b); err != nil {
		log.Fatal("block can not be serialized.")
	} else if err = bc.db.Put([]byte("b_"+b.GetHashCurr()), bs, nil); err != nil {
		log.Fatal("block can not be saved.")
	}

	// 将最新的区块hash存储到数据库中
	if err := bc.db.Put([]byte("lasthash"), []byte(b.GetHashCurr()), nil); err != nil {
		log.Fatal("lastHash can not be saved.")
	}

	return bc
}

// AddGensisBlock 添加创世区块
func (bc *Blockchain) AddGensisBlock(address wallet.Address) *Blockchain {
	// 校验是否可以添加创世区块
	if bc.lasthash != "" {
		return bc
	}

	// 将区块加入到区块链中
	// 创世区块只有 txs 是特殊的
	// bc.AddBlock("The Gensis Block.")
	bc.AddBlock(address)

	return bc
}

// GetBlock 通过hash获取区块
func (bc *Blockchain) GetBlock(hash block.Hash) (*block.Block, error) {
	// 从数据库中读取对应的区块
	data, err := bc.db.Get([]byte("b_"+hash), nil)
	if err != nil {
		return nil, err
	}
	// 反序列化
	b, err := block.BlockUnSerialize(data)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// Iterate 迭代展示区块的方法
func (bc *Blockchain) Iterate() {
	// 从最后一个区块的hash开始迭代
	// for hash := bc.lasthash; hash != ""; {
	//	b := bc.blocks[hash]
	//	fmt.Println("HashCurr", b.hashCurr)
	//	fmt.Println("Txs", b.txs)
	//	fmt.Println("Time", b.header.time.Format(time.UnixDate))
	//	fmt.Println("HashPrevBlock", b.header.hashPrevBlock)
	//	fmt.Println("")
	//	hash = b.header.hashPrevBlock
	// }

	// 从最后一个区块的hash开始迭代
	for hash := bc.lasthash; hash != ""; {
		// 得到区块
		b, err := bc.GetBlock(hash)
		if err != nil {
			log.Fatalf("Block <%s> is not exsits.", hash)
		}
		// 做hashcash验证
		pow := pow.NewPow(b)
		if !pow.Validate() {
			log.Fatalf("Block <%s> is not validate.", hash)
			continue
		}
		fmt.Println("HashCurr:", b.GetHashCurr())
		txs := b.GetTxs()
		fmt.Printf("Txs: %d transactions:\n", len(txs))
		// fmt.Println("Txs:", b.GetTxs())
		for i, t := range txs {
			fmt.Printf("\tindex:%d,hash:%s,Inputs:%d,Outputs:%d\n", i, t.Hash, len(t.Inputs), len(t.Outputs))
		}
		fmt.Println("Time:", b.GetTime().Format(time.UnixDate))
		fmt.Println("HashPrevBlock:", b.GetHashPrevBlock())
		fmt.Println("")
		hash = b.GetHashPrevBlock()
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

	// 清理UTXO缓存
	iteru := bc.db.NewIterator(util.BytesPrefix([]byte("t_")), nil)
	for iteru.Next() {
		bc.db.Delete(iteru.Key(), nil)
	}
	iteru.Release()
	// 清空bc对象
	bc.lasthash = ""

	// TODO-petrichor 清理txs的持久存储
	bc.db.Delete([]byte("txs"), nil)
}
