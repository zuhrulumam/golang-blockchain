package blockchain

// TxInput Transaction Input
type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

// TxOutput Transaction Output
type TxOutput struct {
	Value  int
	PubKey string
}

// Transaction model
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}
