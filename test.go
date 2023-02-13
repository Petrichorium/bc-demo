package main

import (
	"fmt"
	"github.com/cn-org-Pretichor/bc-demo/blockchain"
)

func main() {
	b := blockchain.NewBlock("", "Gensis Block")
	fmt.Println(b)
}
