package blockchain

import (
	"github.com/dgraph-io/badger"
)

// // BlockChain model
// type BlockChain struct {
// 	Blocks []*Block
// }

// // InitChain init block chain with genesis block
// func InitChain() *BlockChain {
// 	return &BlockChain{[]*Block{Genesis()}}
// }

const (
	dbpath = "./tmp/badger"
)

// BlockChain model
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

// InitChain init block chain with genesis block
func InitChain() *BlockChain {

	var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = dbpath
	opts.ValueDir = dbpath

	db, err := badger.Open(opts)
	Handle(err)

	// check if there is a blockchain via lasthash key

	// if not
	// create genesis block
	// save hash to database with lasthash key --> to disk
	// save serialize block to database with it's hash key --> to disk
	// add genesis hash to lasthash so we get lasthash in memory --> in memory

	// if there is a block chain
	// create block with lasthash key as previous hash
	// add to chain
	// save hash to database with lasthash key
	// save serialize block to database with it's hash key --> to disk

}

// Genesis iniial block for the chain
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

// AddBlock add block to a blockchain
func (chain *BlockChain) AddBlock(data string) {
	// get previous block
	prevBlock := chain.Blocks[len(chain.Blocks)-1]

	// get previous hash from previous chain
	newBlock := CreateBlock(data, prevBlock.Hash)

	chain.Blocks = append(chain.Blocks, newBlock)
}
