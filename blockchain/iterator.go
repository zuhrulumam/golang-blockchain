package blockchain

import (
	"github.com/dgraph-io/badger"
)

// BlockChainIterator model
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// Iterator get current block in the blockchain as iterator
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

// Next looping through database to get block in blockchain
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	// get block from deserialize currenthash
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)

		encodedBlock, err := item.Value()

		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	// change iteration currenthash as block prevhash
	iter.CurrentHash = block.PrevHash

	return block
}
