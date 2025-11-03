package blockchain

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp    int64
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
	Height       int
	Difficulty   int    // Mining difficulty used for this block
	MerkleRoot   []byte // Merkle root of transactions (calculated once, stored for validation)
}

// HashTransactions returns the hash of all transactions using Merkle Tree
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}
	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data
}

// CreateBlock creates a new block with transactions
func CreateBlock(txs []*Transaction, prevHash []byte, height int) *Block {
	return CreateBlockWithInterrupt(txs, prevHash, height, nil)
}

func CreateBlockWithInterrupt(txs []*Transaction, prevHash []byte, height int, interrupt <-chan bool) *Block {
	// Use UTC timestamp to ensure consistency across different timezones
	block := &Block{
		Timestamp:    time.Now().UTC().Unix(),
		Hash:         []byte{},
		Transactions: txs,
		PrevHash:     prevHash,
		Nonce:        0,
		Height:       height,
		Difficulty:   Difficulty,
		MerkleRoot:   []byte{}, // Will be calculated by HashTransactions
	}

	// Calculate and store Merkle Root ONCE
	block.MerkleRoot = block.HashTransactions()

	pow := NewProof(block)
	nonce, hash := pow.RunWithInterrupt(interrupt)

	// If hash is nil, mining was interrupted
	if hash == nil {
		return nil
	}

	block.Hash = hash
	block.Nonce = nonce
	return block
}

func CreateBlockWithDifficulty(txs []*Transaction, prevHash []byte, height int, difficulty int) *Block {
	// Use UTC timestamp to ensure consistency across different timezones
	block := &Block{
		Timestamp:    time.Now().UTC().Unix(),
		Hash:         []byte{},
		Transactions: txs,
		PrevHash:     prevHash,
		Nonce:        0,
		Height:       height,
		Difficulty:   difficulty,
		MerkleRoot:   []byte{}, // Will be calculated by HashTransactions
	}

	// Calculate and store Merkle Root ONCE
	block.MerkleRoot = block.HashTransactions()

	pow := NewProofWithDifficulty(block, difficulty)
	nonce, hash := pow.RunWithInterrupt(nil)

	// If hash is nil, mining was interrupted (shouldn't happen for genesis)
	if hash == nil {
		return nil
	}

	block.Hash = hash
	block.Nonce = nonce
	return block
}

// Genesis creates the genesis block with a coinbase transaction
// Uses lower difficulty (GenesisDifficulty) for faster initialization
func Genesis(coinbase *Transaction) *Block {
	return CreateBlockWithDifficulty([]*Transaction{coinbase}, []byte{}, 0, GenesisDifficulty)
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	Handle(encoder.Encode(b))

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	Handle(decoder.Decode(&block))

	return &block
}
