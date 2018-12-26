package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/zuhrulumam/golang-blockchain/blockchain"
)

// CommandLine model
type CommandLine struct {
	blockchain *blockchain.BlockChain
}

// printUsage print cli usages
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage : ")
	fmt.Println("add -block BLOCK DATA - add block with data ")
	fmt.Println("print - print blockchain ")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		fmt.Println("Argument less then 2")
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Block Added")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Data %s \n", block.Data)
		fmt.Printf("prevHash %x \n", block.PrevHash)
		fmt.Printf("Hash %x \n", block.Hash)
		fmt.Printf("Nonce %d \n", block.Nonce)

		if len(block.PrevHash) == 0 {
			break
		}

	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block Data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "print":
		err := printCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printCmd.Parsed() {
		cli.printChain()
	}

}

func main() {
	defer os.Exit(0)
	chain := blockchain.InitChain()
	defer chain.Database.Close()

	cli := CommandLine{chain}

	cli.run()

	// chain := blockchain.InitChain()

	// chain.AddBlock("Test")

	// for _, block := range chain.Blocks {
	// 	fmt.Printf("Data %s \n", block.Data)
	// 	fmt.Printf("prevHash %x \n", block.PrevHash)
	// 	fmt.Printf("Hash %x \n", block.Hash)
	// 	fmt.Printf("Nonce %d \n", block.Nonce)
	// }
}
