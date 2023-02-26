package pow

import (
	"crypto/sha256"
	"fmt"
	"github.com/cn-org-Pretichor/bc-demo/block"
	"math"
	"math/big"
	"strconv"
)

type ProofOfWork struct {
	// 需要pow工作量区块的区块
	block *block.Block
	// 证明参数目标
	target *big.Int
}

func NewPow(b *block.Block) *ProofOfWork {
	p := &ProofOfWork{
		block:  b,
		target: big.NewInt(1),
	}
	// 计算target
	p.target.Lsh(p.target, uint(block.HashLen-b.GetBits()+1))
	return p
}

// Proof
//
//	@Description: hashcash 证明
//	@Author petrichor 2023-02-26 18:12:10
//	@receiver p --
//	@return int --使用的nonce 挖矿计数
//	@return block.Hash --形成的区块hash
func (p *ProofOfWork) Proof() (int, block.Hash) {
	var hashInt big.Int
	// 基于block 准备 serviceStr
	serviceStr := p.block.GenServiceStr()
	// nonce 计数器
	nonce := 0
	// 迭代计算hash，设置防止nonce溢出的条件
	fmt.Printf("Target:%d\n", p.target)
	for nonce <= math.MaxInt64 {
		// 生成hash
		data := serviceStr + strconv.Itoa(nonce)
		hash := sha256.Sum256([]byte(data))
		// 得到hash的big.Int
		hashInt.SetBytes(hash[:])
		fmt.Printf("Hash:%s\t%d\n", hashInt.String(), nonce)
		// 判断是否满足难度（数学难题）
		if hashInt.Cmp(p.target) == -1 {
			// 解决问题
			return nonce, block.Hash(fmt.Sprintf("%x", hash))
		}
		nonce++
	}
	return 0, ""
}

// Validate
//
//	@Description: 验证pow 验证区块hash是否正确
//	@Author petrichor 2023-02-26 18:08:39
//	@receiver p --POW（工作量证明）
//	@return bool --布尔值
func (p *ProofOfWork) Validate() bool {
	// 再次生成hash
	serviceStr := p.block.GenServiceStr()
	data := serviceStr + strconv.Itoa(p.block.GetNonce())
	hash := sha256.Sum256([]byte(data))
	// 比较是否相等
	if p.block.GetHashCurr() != fmt.Sprintf("%x", hash) {
		return false
	}
	// 比较是否满足难题
	target := big.NewInt(1)
	target.Lsh(target, uint(block.HashLen-p.block.GetBits()+1))
	hashInt := new(big.Int)
	hashInt.SetBytes(hash[:])
	// 不小于
	if hashInt.Cmp(target) != -1 {
		return false
	}
	return true
}
