package block

import (
	"bytes"
	"encoding/gob"
	"time"
)

// 中间的数据结构，区块数据
type BlockData struct {
	Version         int
	HashPrevBlock   Hash      //前一个区块的Hash
	HashMerkleBlock Hash      //默克尔树的哈希节点
	Time            time.Time //区块的创建时间
	Bits            int       //难度相关
	Nonce           int       //挖矿相关

	Txs        string //交易列表
	TxsCounter int    //交易计数器
	HashCurr   Hash   //当前区块哈希值缓存
}

// 区块序列化
func BlockSerialize(b Block) ([]byte, error) {
	//由于区块的字段都是 unexported field 非导出字段
	//使用中间的数据结构作为桥梁，完成序列化
	//将区块数据复制到BlockData上
	bd := BlockData{
		Version:         b.header.version,
		HashPrevBlock:   b.header.hashPrevBlock,
		HashMerkleBlock: b.header.hashMerkleBlock,
		Time:            b.header.time,
		Bits:            b.header.bits,
		Nonce:           b.header.nonce,
		Txs:             b.txs,
		TxsCounter:      b.txsCounter,
		HashCurr:        b.hashCurr,
	}
	//执行gob序列化即可
	buffer := bytes.Buffer{}
	//编码器
	enc := gob.NewEncoder(&buffer)
	//编码，序列化
	if err := enc.Encode(bd); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func BlockUnSerialize(data []byte) (Block, error) {
	//得到装有内容的缓存
	buffer := bytes.Buffer{}
	buffer.Write(data)
	//解码器
	dec := gob.NewDecoder(&buffer)
	//解码，反序列化
	bd := BlockData{}
	if err := dec.Decode(&bd); err != nil {
		return Block{}, err
	}
	//反序列化成功
	return Block{
		header: BlockHeader{
			version:         bd.Version,
			hashPrevBlock:   bd.HashPrevBlock,
			hashMerkleBlock: bd.HashMerkleBlock,
			time:            bd.Time,
			bits:            bd.Bits,
			nonce:           bd.Nonce,
		},
		txs:        bd.Txs,
		txsCounter: bd.TxsCounter,
		hashCurr:   bd.HashCurr,
	}, nil
}
