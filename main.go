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
}

// createBlockChain create a blockchain with corresponding address
func (cli *CommandLine) createBlockChain(address string) {
	chain := blockchain.InitChain(address)
	chain.Database.Close()
	fmt.Println("Blockchain Created")
}

// getBalance get balance on particular blockchain address
func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	UTxO := chain.FindUTxO(address)

	for _, out := range UTxO {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

// send token from address to pointed address
func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success send token")
}

// printUsage print cli usages
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage : ")
	fmt.Println("getbalance -address ADDRESS - print blockchain ")
	fmt.Println("createblockchain -address ADDRESS - print blockchain ")
	fmt.Println("send -from FROM -to TO -amount AMOUNT - print blockchain ")
	fmt.Println("printchain - print blockchain ")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		fmt.Println("Argument less then 2")
		cli.printUsage()
		runtime.Goexit()
	}
}

// func (cli *CommandLine) addBlock(data string) {
// 	cli.blockchain.AddBlock(data)
// 	fmt.Println("Block Added")
// }

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()

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

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printCmd := flag.NewFlagSet("print", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "Address of Balance")
	createBlockChainAddress := createBlockChainCmd.String("address", "", "Create This Chain Address")
	sendFrom := sendCmd.String("from", "", "Source Wallet Address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "print":
		err := printCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainAddress == "" {
			createBlockChainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockChainAddress)
	}

	if sendCmd.Parsed() {

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printCmd.Parsed() {
		cli.printChain()
	}

}

func main() {
	defer os.Exit(0)

	cli := CommandLine{}

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
