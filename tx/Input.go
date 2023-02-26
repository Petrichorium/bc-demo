package tx

type Input struct {
	HashSrcTx      string //输入来源的 交易 的 hash
	IndexSrcOutput int    //输入来源的 交易 的 输出 的 索引
	Signature      []byte //交易输入签名
	PublicKey      string //公钥
}
