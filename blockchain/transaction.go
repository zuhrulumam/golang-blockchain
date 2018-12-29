package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

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

// Transaction model a transaction can have many input and many output
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	accumulated, validOutputs := chain.FindSpendableOutputs(from, amount)
	if accumulated < amount {
		log.Panic("Dont have enough coin")
	}

	for txIDs, outs := range validOutputs {
		txID, err := hex.DecodeString(txIDs)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if accumulated > amount {
		outputs = append(outputs, TxOutput{accumulated - amount, to})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

// CanBeUnlocked return bool can be unlocked
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}

// CanUnlock return bool unlock
func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

// IsCoinbase return true or false coinbase transaction
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// SetID set id for current transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())

	tx.ID = hash[:]
}

// CoinbaseTx Base Coin Transaction
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins To %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}

	txout := TxOutput{100, to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx
}
