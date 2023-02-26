package tx

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/cn-org-Pretichor/bc-demo/wallet"
	"log"
	"time"
)

// 一个satoshi等于一亿分之一的 BTC(0.00000001 BTC，汶也是比特市单面最小的货市单位（就像是 1分的硬市）
const satoshi = 1 // 1 中本聪
const s = satoshi
const Ks = 1000 * 5  //千
const Ms = 1000 * Ks //百万
const GS = 1000 * Ms //十亿
const BTC = 100000000 * satoshi

// CoinbaseSubsidy 挖矿奖励金
const CoinbaseSubsidy = 12 * BTC

type TX struct {
	Hash string
	// 输入集合
	Inputs []*Input
	// 输出集合
	Outputs []*Output
	//交易时间
	Time time.Time
}

// Tx的构造器
func NewTX(ins []*Input, outs []*Output) *TX {
	tx := &TX{
		Inputs:  ins,
		Outputs: outs,
		Time:    time.Now(),
	}
	tx.SetHash()

	return tx
}

// 新建转账交易
func NewTransferTX(ins []*Input, outs []*Output) *TX {
	return NewTX(ins, outs)
}

// CoinBase交易构造器
func NewCoinBaseTX(to wallet.Address) *TX {
	//输入，没有输入
	ins := []*Input{}
	//输出，仅存在一个输出，给目标为to的用户挖矿奖励
	output := &Output{
		Value: CoinbaseSubsidy, // 挖矿奖励金
		To:    to,
	}
	outs := []*Output{
		output,
	}

	return NewTX(ins, outs)
}

func (tx *TX) SetHash() *TX {
	//先序列化 gob
	ser, err := SerializeTx(*tx)
	if err != nil {
		log.Fatal(err)
	}
	//hash
	//再生成hash sha256
	hash := sha256.Sum256(ser)
	tx.Hash = fmt.Sprintf("%x", hash)
	return tx
}

// 序列化Tx数据
func SerializeTx(tx TX) ([]byte, error) {
	buffer := bytes.Buffer{}
	enc := gob.NewEncoder(&buffer)
	//序列化
	if err := enc.Encode(tx); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// 反序列化Tx数据 (反串行化，解码)
func UnserializeTx(data []byte) (TX, error) {
	buffer := bytes.Buffer{}
	buffer.Write(data)
	dec := gob.NewDecoder(&buffer)
	tx := TX{}
	if err := dec.Decode(&tx); err != nil {
		return tx, err
	}
	return tx, nil
}
