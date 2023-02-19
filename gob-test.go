package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Block struct {
	CurrHash string
	Txs      string
}

func main() {
	//Block数据，数据不需要指针类型
	b := &Block{
		CurrHash: "123415666",
		Txs:      "first transaction",
	}

	//gob编码
	//先得到编码器
	//定义可以写入内容的容器，通常使用byte型的缓存
	//提供的缓存应该具备可写功能
	var network bytes.Buffer // Stand-in for a network connection
	//编码器需要该缓存，将编码的结果写入该缓存
	enc := gob.NewEncoder(&network) // Will write to network.
	//编码数据，编码器的结果写入了编码器的缓存中
	enc.Encode(b)
	fmt.Println(network.Bytes(), network.String())
	result := network.Bytes()

	//解码
	//解码时，解码的数据从缓存中获取
	//提供的缓存应该具备可读功能
	var bbr bytes.Buffer
	//将之前编码的数据，放入缓存中
	bbr.Write(result)
	dec := gob.NewDecoder(&bbr)
	//解码时，需要提供解码的数据类型
	b1 := Block{}
	dec.Decode(&b1)
	fmt.Println(b1)
}
