package blockchain

import (
	"fmt"

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

	// Checking and also updating database with blockchain
	err = db.Update(func(txn *badger.Txn) error {
		// check if there is a blockchain via lasthash key
		// if not
		if _, err := txn.Get([]byte("lasthash")); err == badger.ErrKeyNotFound {
			fmt.Println("No Blockchain Exist. Creating One ...")
			// create genesis block
			genesis := Genesis()

			// save hash to database with lasthash key --> to disk
			// the purpose of this is later when we want to get blocks from block chain we can get from deserialized this block
			err := txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)

			// save serialize block to database with it's hash key --> to disk
			err = txn.Set([]byte("lasthash"), genesis.Hash)

			// set lasthash with genesis hash so we get lasthash in memory --> in memory
			lastHash = genesis.Hash

			return err

		} else {
			// if there is a block chain
			// get lasthash value
			item, err := txn.Get([]byte("lasthash"))
			Handle(err)

			// set lasthash to lasthash in memory --> in memory
			lastHash, err = item.Value()

			return err
		}
	})

	Handle(err)

	blockchain := BlockChain{lastHash, db}

	return &blockchain

}

// Genesis iniial block for the chain
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

// AddBlock add block to a blockchain
func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	// get previous hash via database
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lasthash"))
		Handle(err)

		// set lasthash as the value returned from db
		lastHash, err = item.Value()

		return err
	})

	// create block with data and lasthash
	newBlock := CreateBlock(data, lastHash)

	// save new block to database
	err = chain.Database.Update(func(txn *badger.Txn) error {
		// save serialize block with hash as a key --> to disk
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)

		// save hash to database with lasthash key --> to disk
		err = txn.Set([]byte("lasthash"), newBlock.Hash)

		// set chain lasthash in memory with hash
		chain.LastHash = newBlock.Hash

		return err
	})

	Handle(err)

	// // get previous block
	// prevBlock := chain.Blocks[len(chain.Blocks)-1]

	// // get previous hash from previous chain
	// newBlock := CreateBlock(data, prevBlock.Hash)

	// chain.Blocks = append(chain.Blocks, newBlock)
}
