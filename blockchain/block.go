package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// Block model
type Block struct {
	Hash []byte
	// Data     []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

// HashTransactions create a hash from transactions
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// Handle error
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// Serialize change block into byte for database purpose
func (block *Block) Serialize() []byte {
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)

	err := encoder.Encode(block)

	Handle(err)

	return buff.Bytes()
}

// Deserialize turn bytes from database into block
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

// CreateHash create hash on particular block
// func (b *Block) CreateHash() {
// 	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
// 	hash := sha256.Sum256(info)
// 	b.Hash = hash[:]
// }

// CreateBlock create block with hash, data, prevhash
func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{
		Hash:         []byte{},
		Transactions: txs,
		PrevHash:     prevHash,
		Nonce:        0,
	}

	pow := NewProof(block)

	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block

}
