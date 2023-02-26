package block

import (
	"fmt"
	"github.com/cn-org-Pretichor/bc-demo/tx"
	"time"
)

type Hash = string

const HashLen = 256
const nodeVersion = 0
const blockBits = 12

// Block 区块主体
type Block struct {
	header BlockHeader
	// txs        string //交易列表
	txs        []*tx.TX // 交易列表
	txsCounter int      // 交易计数器
	hashCurr   Hash     // 当前区块哈希值缓存
}

// BlockHeader 区块头
type BlockHeader struct {
	version         int
	hashPrevBlock   Hash      // 前一个区块的Hash
	hashMerkleBlock Hash      // 默克尔树的哈希节点 表示区块所存储的交易的 默克尔树的根
	time            time.Time // 区块的创建时间
	bits            int       // 难度相关
	nonce           int       // 挖矿相关
}

// NewBlock 构造区块
func NewBlock(prevHash Hash) *Block {
	b := &Block{
		header: BlockHeader{
			version:       nodeVersion,
			hashPrevBlock: prevHash, // 设置前面区块的hash
			time:          time.Now(),
			bits:          blockBits,
		},
		txs:        []*tx.TX{},
		txsCounter: 0,
	}
	// 计算设置当前区块的hash
	// b.SetHashCurr()

	return b
}

// Stringify 头信息的字符串化
// func (bh *BlockHeader) Stringify() string {
//	return fmt.Sprintf("%d%s%s%d%d%d",
//		bh.version,
//		bh.hashPrevBlock,
//		bh.hashMerkleBlock,
//		bh.time.UnixNano(), //得到时间戳，nano级别
//		bh.bits,
//		bh.nonce,
//	)
// }

// SetHashCurr 设置当前区块哈希
// func (b *Block) SetHashCurr() *Block {
//	//生成头信息的拼接字符串
//	headerStr := b.header.Stringify()
//	//计算哈希值
//	b.hashCurr = fmt.Sprintf("%x", sha256.Sum256([]byte(headerStr)))
//	return b
// }

// 添加交易
func (b *Block) AddTx(tx *tx.TX) *Block {
	b.txs = append(b.txs, tx)
	b.txsCounter++
	return b
}

// SetNonce 设置Nonce
func (b *Block) SetNonce(nonce int) {
	b.header.nonce = nonce
}

// SetHashCurr 设置当前区块哈希
func (b *Block) SetHashCurr(hash Hash) {
	b.hashCurr = hash
}

// GetHashCurr 获取当前区块哈希
func (b Block) GetHashCurr() Hash {
	return b.hashCurr
}

// GetBits bits属性的getter
func (b Block) GetBits() int {
	return b.header.bits
}

// GenServiceStr
//
//	@Description: 生成用于pow的 服务字符串（生成区块hash的一部分字符串）
//	@Author petrichor 2023-02-26 18:04:58
//	@receiver b 区块
//	@return string blockHeader的属性组成的字符串
func (b *Block) GenServiceStr() string {
	return fmt.Sprintf("%d%s%s%s%d",
		b.header.version,
		b.header.hashPrevBlock,
		b.header.hashMerkleBlock,
		// b.header.time.String(), //得到时间戳，nano级别
		b.header.time.Format("2006-01-02 15:94:05.999999999 -0780 MST"),
		b.header.bits,
	)
}

func (b Block) GetTxs() []*tx.TX {
	return b.txs
}
func (b Block) GetTime() time.Time {
	return b.header.time
}
func (b Block) GetHashPrevBlock() Hash {
	return b.header.hashPrevBlock
}
func (b Block) GetNonce() int {
	return b.header.nonce
}
