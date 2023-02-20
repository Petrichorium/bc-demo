package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strconv"
)

func main() {
	bits := 8               //256位前8位为0
	target := big.NewInt(1) //00000000 ... 00000001
	//采用左移方案，构建比较数
	//00000001 LSH 1 = 00000010
	//00000001 LSH 2 = 00000100
	target.Lsh(target, uint(256-bits+1))
	fmt.Println(target.String())
	fmt.Println("------Minting------")

	nonce := 0
	//服务字符串
	serviceStr := "block data"
	var hashInt big.Int
	for {
		//服务字符串 连接 随机数
		data := serviceStr + strconv.Itoa(nonce)
		hash := sha256.Sum256([]byte(data))
		hashInt.SetBytes(hash[:])
		fmt.Println(hashInt.String(), nonce)
		//break
		if hashInt.Cmp(target) == -1 { //compare 比较  ==-1时：hashInt < target
			fmt.Println("本机挖矿成功")
			return
		}
		nonce++
	}
}
