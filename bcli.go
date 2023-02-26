package main

import (
	"flag"
	"fmt"
	"github.com/cn-org-Pretichor/bc-demo/blockchain"
	"github.com/cn-org-Pretichor/bc-demo/tx"
	"github.com/cn-org-Pretichor/bc-demo/wallet"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"strconv"
	"strings"
)

// 命令行工具
/* 参数接收 os.Args
用于获取命令行的全部参数，为[]string结构，其中0元素固定，为当前执行的脚本名
*/
func main() {
	// 初始化数据库
	// 数据库链接
	dbpath := "data"
	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 释放数据库链接
	defer db.Close()

	// 初始化区块链
	bc := blockchain.NewBlockchain(db)
	// 添加创世区块
	// bc.AddGensisBlock()

	// 初始化第一个命令参数
	arg1 := ""
	if len(os.Args) >= 2 {
		arg1 = os.Args[1]
	}
	switch strings.ToLower(arg1) {
	case "create:block":
		// 为 create:block 命令增加一个 flag集合，标志集合
		// flag.ExitOnError 的错误处理为，一旦解析失败，则 exit
		fs := flag.NewFlagSet("create:block", flag.ExitOnError)
		// 在集合中，添加需要解析的flag标志
		// txs := fs.String("txs", "", "")
		address := fs.String("address", "", "")
		// 解析命令行参数
		fs.Parse(os.Args[2:])
		bc.AddBlock(*address)

	// 判断是否解析成功
	// if !fs.Parsed() {
	//	log.Fatal("createblock args parsed error.")
	// }
	// fmt.Println(txs, *txs)
	case "balance":
		fs := flag.NewFlagSet("balance", flag.ExitOnError)
		address := fs.String("address", "", "")
		fs.Parse(os.Args[2:])

		fmt.Printf("Address:%s\nBalance:%s\n", *address, strconv.Itoa(bc.GetBalance(*address)))
	case "send":
		newFlagSet := flag.NewFlagSet("send", flag.ExitOnError)
		from := newFlagSet.String("from", "", "user from")
		to := newFlagSet.String("to", "", "user to")
		value := newFlagSet.Int("value", 0, "value")
		newFlagSet.Parse(os.Args[2:])

		if newFlagSet.Parsed() { // parsed() 解析成功
			err := bc.Transfer(*from, *to, *value*tx.BTC)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println("Send Success! waiting for block packed.")
			}

		}

	case "show":
		// 展示全部区块
		bc.Iterate()
	case "init":
		// 为 init 命令增加一个 flag集合，标志集合
		// flag.ExitOnError 的错误处理为，一旦解析失败，则 exit
		fs := flag.NewFlagSet("init", flag.ExitOnError)
		// 在集合中，添加需要解析的flag标志
		address := fs.String("address", "", "")
		// 解析命令行参数
		fs.Parse(os.Args[2:])
		if fs.Parsed() {
			if *address == "" {
				log.Println("Please set address by flag -address.")
			} else {
				// 删除已有的全部区块，增加重新增加一个区块
				// 清空，真实情况不应该有Clear操作
				bc.Clear()

				// 增加创世区块
				bc.AddGensisBlock(*address)
			}
		}
	case "create:wallet":
		// 为 create:wallet 命令增加一个 flag集合，标志集合
		// flag.ExitOnError 的错误处理为，一旦解析失败，则 exit
		fs := flag.NewFlagSet("create:wallet", flag.ExitOnError)
		// 在集合中添加flag标志
		pass := fs.String("pass", "", "")
		// 解析命令行参数
		fs.Parse(os.Args[2:])
		w := wallet.NewWallet(*pass)
		fmt.Printf("your Mnemonic: %s \n", w.GetMnemonic())
		fmt.Printf("your Address : %s \n", w.Address)
	case "help":
		fallthrough
	default:
		Usage()
	}
}

// Usage 输出 bcli 的帮助信息
func Usage() {
	fmt.Println("Bcli is a tool for Blockchnia.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("\t%s\t\t%s\n", "bcli help", "help info for bcli.")
	// @Modified By petrichor
	// @Modify description init 命令增加了 address
	fmt.Printf("\t%s\t\t%s\n", "bcli init -adddress=<adddress>", "initial blockchain.")
	// fmt.Printf("\t%s\t%s\n", "bcli create:block -txs=<txs>", "create block on blockchain.")
	fmt.Printf("\t%s\t%s\n", "bcli create:block -adddress=<adddress>", "create block on blockchain.")
	fmt.Printf("\t%s\t%s\n", "bcli send -from=A -to=B -value=N", "")
	fmt.Printf("\t%s\t\t%s\n", "bcli show", "show blocks in chain.")
	fmt.Printf("\t%s\t\t%s\n", "bcli create:wallet -pass=<pass>", "create wallet base on pass")
}
