package blockchain

import (
	"github.com/cn-org-Pretichor/bc-demo/wallet"
)

func (bc *Blockchain) GetBalance(address wallet.Address) int {
	// 获取 address 的对应 UTXO
	// 统计余额
	balance := 0
	for _, utxo := range bc.UTXOCache.FindUTXO(address) {
		balance += utxo.Output.Value
	}
	return balance
}
