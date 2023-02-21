package wallet

import (
	"crypto/sha256"
	"github.com/mr-tron/base58"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ripemd160"
	"log"
)

type Address = string

const keyBitSizes = 256

type Wallet struct {
	//助记词
	mnemonic string
	//私钥为 *bip32.Key 类型
	privateKey *bip32.Key
	//公钥由私钥计算推导，使用下面的调用
	//publicKey = privateKey.PublicKey()
	Address Address
}

// NewWallet 构造函数
func NewWallet(pass string) *Wallet {
	w := &Wallet{}
	//生成密钥
	w.GenKey(pass)
	//生成地址
	w.GenAddress()

	return w
}

// 生成key
func (w *Wallet) GenKey(pass string) *Wallet {
	////elliptic.P256() 生成椭圆
	////rand.Reader 生成随机数
	//privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//if err != nil {
	//	//panic(err)
	//	log.Fatal(err)
	//}
	//w.privateKey = privateKey
	//return w

	//使用 bip39
	//墒（随机）
	entropy, err := bip39.NewEntropy(keyBitSizes)
	if err != nil {
		log.Fatal(err)
	}
	//助记词
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatal(err)
	}
	//key的种子 seed
	seed := bip39.NewSeed(mnemonic, pass)
	//生成秘钥
	privateKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Fatal(err)
	}
	//生成公钥

	w.mnemonic = mnemonic
	w.privateKey = privateKey
	return w
}

// 生成Address
func (w *Wallet) GenAddress() *Wallet {
	//利用私钥生成公钥
	//pubKey := w.genPubKey() // []byte
	pubKey := w.privateKey.PublicKey().String()

	// ripemd160(sha256(pubkey))
	//shaHash := sha256.Sum256(pubKey)
	//rpmd := ripemd160.New()
	//rpmd.Write(shaHash[:])
	//pubHash := rpmd.Sum(nil)
	pubHash := HashPubKey([]byte(pubKey))

	//计算checkSum 校验值
	h1 := sha256.Sum256(pubHash)
	checkSum := sha256.Sum256(h1[:])

	//组合，继续base64
	data := append(append([]byte{0}, pubHash...), checkSum[:4]...)
	w.Address = base58.Encode(data)

	return w
}

// 利用私钥生成公钥
//func (w *Wallet) genPubKey() []byte {
//	pubKey := append(
//		w.privateKey.PublicKey.X.Bytes(),
//		w.privateKey.PublicKey.Y.Bytes()...,
//	)
//	return pubKey
//}

// 生成公钥hash值
func HashPubKey(pubKey []byte) []byte {
	shaHash := sha256.Sum256(pubKey)
	rpmd := ripemd160.New()
	rpmd.Write(shaHash[:])
	pubHash := rpmd.Sum(nil)
	return pubHash
}

// getMnemonic
func (w *Wallet) GetMnemonic() string {
	return w.mnemonic
}
