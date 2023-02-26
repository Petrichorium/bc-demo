package blockchain

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/cn-org-Pretichor/bc-demo/tx"
	"github.com/cn-org-Pretichor/bc-demo/wallet"
)

// 交易缓存队列
var TXS = []*tx.TX{}

// 与转账相关的操作
func (bc Blockchain) Transfer(from, to wallet.Address, value int) error {

	// 构建交易数据
	// ## 查询当前用户的全部UTXO
	utxos := bc.UTXOCache.FindUTXO(from)
	// ## 够凑转账金额
	amount, spendableUtxo := FindSpendableUTXO(utxos, value)
	// 是否够支付
	if amount < value {
		// 不够支付
		return errors.New(fmt.Sprintf("Balance of %s not enough.", from))
	}
	// 金额足够，说明可以转账，新建交易。
	// ## 构建交易的输入
	inConuter := len(spendableUtxo)
	ins := make([]*tx.Input, inConuter)
	for i, u := range spendableUtxo {
		ins[i] = &tx.Input{
			HashSrcTx:      u.HashSrcTx,
			IndexSrcOutput: u.IndexSrcOutput,
		}
	}
	// ## 构建交易的输出
	outs := []*tx.Output{
		&tx.Output{
			Value: value,
			To:    to,
		},
	}
	// ### 找零输出
	if amount > value {
		outs = append(outs, &tx.Output{
			Value: amount - value,
			To:    from,
		})
	}

	// 构建交易
	t := tx.NewTransferTX(ins, outs)

	// 存储交易，将交易数据放入集合（内存）
	TXS = append(TXS, t)

	// TODO-petrichor 暂时无法接收别人发来的交易，先存入持久化中以便测试
	bc.TxsSave()

	return nil
}

// 凑够转账的 UTXO 金额
// 共：1，2，4，6，10 的utxo
// 需要 8 的 utxo
// 返回：13, [1, 2, 4, 6]UTX0
func FindSpendableUTXO(utxos []*UTXO, value int) (int, []*UTXO) {
	amount := 0 // 已经凑的金额
	sautxos := []*UTXO{}
	// 遍历全部的utxo
	for _, u := range utxos {
		amount += u.Output.Value
		sautxos = append(sautxos, u)
		// 凑够了
		if amount >= value {
			break
		}
	}
	return amount, sautxos

}

//
// TxsSave
//  @Description: txs 持久化，以便测试
//  @Author petrichor 2023-02-26 14:23:44
//  @receiver bc
//  @return error
//
func (bc *Blockchain) TxsSave() error {
	// 序列化
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(TXS)
	if err != nil {
		return err
	}

	// 存储
	key := []byte("txs")
	bc.db.Put(key, buffer.Bytes(), nil)
	if err != nil {
		return err
	}

	return nil
}

//
// TxsInit
//  @Description: txs反序列化，以便测试时获取txs
//  @Author petrichor 2023-02-26 14:24:19
//  @receiver bc
//  @return error
//
func (bc *Blockchain) TxsInit() error {
	// 获取
	key := []byte("txs")
	data, err := bc.db.Get(key, nil)
	if err != nil {
		return err
	}

	// 解码
	buffer := bytes.Buffer{}
	buffer.Write(data)
	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&TXS)

	return err
}
