package blockchain

import (
	"bytes"
	"encoding/gob"
	"github.com/cn-org-Pretichor/bc-demo/tx"
	"github.com/cn-org-Pretichor/bc-demo/wallet"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
)

// UTXO （Unspend Transaction Outputs）未消费的交易输出
// 带有缓存的UTXO设计

// UTXO 缓存操作对象
type UTXOCache struct {
	db *leveldb.DB
}

func NewUTXOCache(db *leveldb.DB) *UTXOCache {
	return &UTXOCache{
		db: db,
	}
}

// UTXO结构
type UTXO struct {
	Output         *tx.Output // output本身存储
	HashSrcTx      string     // 所属的 交易 的 hash
	IndexSrcOutput int        // 位于交易输出列表的索引
}

// 单个交易的UTXO集合
type UTXOSet = []*UTXO

// Update
//
//	@Description: 更新UTXO缓存
//	@Author Petrichor 2023-02-26 11:48:47
//	@receiver uc
//	@param tx 新的交易
//	@return *UTXOCache
func (uc *UTXOCache) Update(tx *tx.TX) *UTXOCache {
	// 添加UTXO缓存
	// 遍历tx交易全部的输出，每个输出作为一个新的UTXO缓存来使用即可
	us := UTXOSet{}
	for i, outputt := range tx.Outputs {
		us = append(us, &UTXO{
			Output:         outputt,
			HashSrcTx:      tx.Hash,
			IndexSrcOutput: i,
		})
	}
	ser, err := SerializeUTXOSet(us)
	if err != nil {
		log.Fatal("UTXOSet Serialize failed.")
	}
	// 更新对应的缓存即可
	// []byte("b_"+b.GetHashCurr()
	uc.db.Put([]byte("t_"+tx.Hash), ser, nil)

	// 清理已消费的UTXO缓存
	// 遍历全部新交易的输入
	for _, in := range tx.Inputs {
		// 基于输入的 来源交易的hash 找到 UTXOCache 中的 UTXO
		key := []byte("t_" + in.HashSrcTx)
		data, err := uc.db.Get(key, nil)
		if err != nil {
			log.Println(err)
		}
		// 反序列化
		utxoSet, err := UnSerializeUTXOSet(data)
		if err != nil {
			log.Println(err)
		}
		// 更新该 utxo集合
		newUtxoSet := UTXOSet{}
		for _, utxo := range utxoSet {
			if utxo.IndexSrcOutput != in.IndexSrcOutput {
				// 未消费的 utxo 再存储起来
				newUtxoSet = append(newUtxoSet, utxo)
			}
			// 已消费的 utxo 不做处理
		}

		// 去掉已消费后的，为消费的 utxo
		if len(newUtxoSet) == 0 {
			// 没有剩余了 可以删除了
			uc.db.Delete(key, nil)
		} else {
			// 有剩余 继续存起来
			serializeUTXOSet, err := SerializeUTXOSet(newUtxoSet)
			if err != nil {
				log.Println("UTXOSet Serialize failed.")
				return uc
			}
			uc.db.Put(key, serializeUTXOSet, nil)
		}
	}

	return uc
}

// 序列化单交易的UTXO集合
func SerializeUTXOSet(us UTXOSet) ([]byte, error) {
	buffer := bytes.Buffer{}
	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(us); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func UnSerializeUTXOSet(data []byte) (UTXOSet, error) {
	buffer := bytes.Buffer{}
	buffer.Write(data)
	dec := gob.NewDecoder(&buffer)
	var us UTXOSet
	err := dec.Decode(&us)

	return us, err
}

// 找到哪些 属于 address 的未被消费的输出 Unspend
func (uc *UTXOCache) FindUTXO(address wallet.Address) []*UTXO {
	utxo := []*UTXO{}

	// 遍历utxo缓存，得到utxo即可
	iter := uc.db.NewIterator(util.BytesPrefix([]byte("t_")), nil)
	for iter.Next() {
		value, err := uc.db.Get(iter.Key(), nil)
		if err != nil {
			continue
		}
		// value 是序列化的 utxoSet
		// 反序列化的工作
		us, err := UnSerializeUTXOSet(value)
		if err != nil {
			log.Println(err)
			continue
		}
		// us 中属于 address 的 加入整体的UTXO集合中
		for _, u := range us {
			if u.Output.To == address {
				utxo = append(utxo, u)
			}
		}
	}
	iter.Release()
	return utxo
}
