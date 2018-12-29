package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

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
	dbpath      = "./tmp/badger"
	dbFile      = "./tmp/badger/MANIFEST"
	genesisData = "First Genesis Created"
)

// DBExist return if db created or not
func DBExist() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// BlockChain model
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

// FindSpendableOutputs get spendable output
func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0

	unspentTxs := chain.FindUnspentTransactions(address)

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated > amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// FindUTxO find unspent transaction output
func (chain *BlockChain) FindUTxO(address string) []TxOutput {
	var UTxOs []TxOutput

	unspentTransaction := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransaction {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTxOs = append(UTxOs, out)
			}
		}
	}

	return UTxOs
}

// FindUnspentTransactions find coresponding address unspent transaction
func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spentTxOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTxOs[txID] != nil {
					for _, spentOut := range spentTxOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)

						spentTxOs[inTxID] = append(spentTxOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs
}

// ContinueBlockChain get last block chain created
func ContinueBlockChain(address string) *BlockChain {
	var lastHash []byte

	if DBExist() == false {
		fmt.Println("Blockchain Not Yet Exist")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions
	opts.Dir = dbpath
	opts.ValueDir = dbpath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lasthash"))
		Handle(err)

		lastHash, err = item.Value()

		return err

	})
	Handle(err)

	blockchain := BlockChain{lastHash, db}

	return &blockchain
}

// InitChain init block chain with genesis block
func InitChain(address string) *BlockChain {

	var lastHash []byte

	if DBExist() {
		fmt.Println("Blockchain Already Exist")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions
	opts.Dir = dbpath
	opts.ValueDir = dbpath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		// create transaction with address and genesis data
		cbtx := CoinbaseTx(address, genesisData)

		// create first block with genesis transaction
		genesis := Genesis(cbtx)

		fmt.Println("Genesis Block Created ")

		err := txn.Set([]byte("lasthash"), genesis.Hash)
		Handle(err)

		err = txn.Set(genesis.Hash, genesis.Serialize())

		lastHash = genesis.Hash

		return err
	})

	// // Checking and also updating database with blockchain
	// err = db.Update(func(txn *badger.Txn) error {
	// 	// check if there is a blockchain via lasthash key
	// 	// if not
	// 	if _, err := txn.Get([]byte("lasthash")); err == badger.ErrKeyNotFound {
	// 		fmt.Println("No Blockchain Exist. Creating One ...")
	// 		// create genesis block
	// 		genesis := Genesis()

	// 		// save hash to database with lasthash key --> to disk
	// 		// the purpose of this is later when we want to get blocks from block chain we can get from deserialized this block
	// 		err := txn.Set(genesis.Hash, genesis.Serialize())
	// 		Handle(err)

	// 		// save serialize block to database with it's hash key --> to disk
	// 		err = txn.Set([]byte("lasthash"), genesis.Hash)

	// 		// set lasthash with genesis hash so we get lasthash in memory --> in memory
	// 		lastHash = genesis.Hash

	// 		return err

	// 	} else {
	// 		// else if there is a block chain
	// 		// get lasthash value
	// 		item, err := txn.Get([]byte("lasthash"))
	// 		Handle(err)

	// 		// set lasthash to lasthash in memory --> in memory
	// 		lastHash, err = item.Value()

	// 		return err
	// 	}
	// })

	Handle(err)

	blockchain := BlockChain{lastHash, db}

	return &blockchain

}

// Genesis iniial block for the chain
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

// AddBlock add block to a blockchain
func (chain *BlockChain) AddBlock(txs []*Transaction) {
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
	newBlock := CreateBlock(txs, lastHash)

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
