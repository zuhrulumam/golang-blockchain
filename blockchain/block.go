package blockchain

// InitChain init block chain with genesis block
func InitChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
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

// BlockChain model
type BlockChain struct {
	Blocks []*Block
}

// Block model
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

// CreateHash create hash on particular block
// func (b *Block) CreateHash() {
// 	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
// 	hash := sha256.Sum256(info)
// 	b.Hash = hash[:]
// }

// CreateBlock create block with hash, data, prevhash
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{
		Hash:     []byte{},
		Data:     []byte(data),
		PrevHash: prevHash,
		Nonce:    0,
	}

	pow := NewProof(block)

	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block

}
