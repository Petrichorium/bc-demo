package main

import (
	"github.com/cn-org-Pretichor/bc-demo/blockchain"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

func main() {
	//区块测试
	//b := blockchain.NewBlock("", "Gensis Block")
	//fmt.Println(b)

	//数据库链接
	dbpath := "data"
	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		log.Fatal(err)
	}
	//释放数据库链接
	defer db.Close()

	//区块链测试
	bc := blockchain.NewBlockchain(db)
	//添加创世区块
	bc.AddGensisBlock()

	bc.
		AddBlock("First Block.").
		AddBlock("Second Block")
	//fmt.Println(bc)
	bc.Iterate()
}
