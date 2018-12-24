package main

import (
	"fmt"

	"github.com/zuhrulumam/golang-blockchain/blockchain"
)

func main() {

	chain := blockchain.InitChain()

	chain.AddBlock("Test")

	for _, block := range chain.Blocks {
		fmt.Printf("Data %s \n", block.Data)
		fmt.Printf("prevHash %x \n", block.PrevHash)
		fmt.Printf("Hash %x \n", block.Hash)
	}
}
