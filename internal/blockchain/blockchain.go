package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/syndtr/goleveldb/leveldb"
)

// Database path configuration (uses constant from config.go)
var dbPath = getDBPath()

// getDBPath returns the database path, checking environment variable first
func getDBPath() string {
	if path := os.Getenv("BLOCKCHAIN_DATA_DIR"); path != "" {
		return path + "/blocks"
	}
	return DBPath // Use constant from config.go
}

type Blockchain struct {
	LastHash []byte
	Database *leveldb.DB
}

// BlockchainIterator iterates over blockchain blocks
type BlockchainIterator struct {
	CurrentHash []byte
	Database    *leveldb.DB
}

// InitBlockchain initializes a new blockchain with genesis block
func InitBlockchain(address string) *Blockchain {
	var lastHash []byte

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dbPath, os.ModePerm); err != nil {
		Handle(err)
	}

	db, err := leveldb.OpenFile(dbPath, nil)
	Handle(err)

	// Check if blockchain already exists
	data, err := db.Get([]byte("lh"), nil)
	if err != nil && err != leveldb.ErrNotFound {
		Handle(err)
	}

	if data == nil {
		// No existing blockchain, create genesis
		fmt.Println("No existing blockchain found")
		cbtx := CoinbaseTX(address, GenesisData, 0) // Genesis block is height 0
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created")
		
		err = db.Put(genesis.Hash, genesis.Serialize(), nil)
		Handle(err)
		err = db.Put([]byte("lh"), genesis.Hash, nil)
		Handle(err)

		lastHash = genesis.Hash
	} else {
		// Blockchain exists, load last hash
		lastHash = data
	}

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

// ContinueBlockchain continues an existing blockchain
func ContinueBlockchain(address string) *Blockchain {
	if DBexists() == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dbPath, os.ModePerm); err != nil {
		Handle(err)
	}

	db, err := leveldb.OpenFile(dbPath, nil)
	Handle(err)

	// Load last hash
	data, err := db.Get([]byte("lh"), nil)
	Handle(err)
	lastHash = data

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

// MineBlock mines a new block with the provided transactions
func (chain *Blockchain) MineBlock(transactions []*Transaction) *Block {
	return chain.MineBlockWithInterrupt(transactions, nil)
}

func (chain *Blockchain) MineBlockWithInterrupt(transactions []*Transaction, interrupt <-chan bool) *Block {
	var lastHash []byte
	var lastHeight int

	// Verify all transactions
	for _, tx := range transactions {
		if chain.VerifyTransaction(tx) == false {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	// Get last block info
	data, err := chain.Database.Get([]byte("lh"), nil)
	Handle(err)
	lastHash = data

	// Get last block to retrieve height
	blockData, err := chain.Database.Get(lastHash, nil)
	Handle(err)
	lastBlock := Deserialize(blockData)
	lastHeight = lastBlock.Height

	// Create new block with interrupt support
	newBlock := CreateBlockWithInterrupt(transactions, lastHash, lastHeight+1, interrupt)
	
	// If block is nil, mining was interrupted
	if newBlock == nil {
		return nil
	}

	// Save to database
	err = chain.Database.Put(newBlock.Hash, newBlock.Serialize(), nil)
	Handle(err)
	err = chain.Database.Put([]byte("lh"), newBlock.Hash, nil)
	Handle(err)

	chain.LastHash = newBlock.Hash

	return newBlock
}

// AddBlock adds a block to the blockchain (used when receiving blocks from network)
func (chain *Blockchain) AddBlock(block *Block) {
	// Check if block already exists
	_, err := chain.Database.Get(block.Hash, nil)
	if err == nil {
		return // Block already exists
	}

	// Validate block data
	blockData := block.Serialize()
	
	// Save block
	err = chain.Database.Put(block.Hash, blockData, nil)
	Handle(err)

	// Get current last block
	lastData, err := chain.Database.Get([]byte("lh"), nil)
	Handle(err)
	lastBlockData, err := chain.Database.Get(lastData, nil)
	Handle(err)
	lastBlock := Deserialize(lastBlockData)

	// Update last hash if new block has greater height
	if block.Height > lastBlock.Height {
		err = chain.Database.Put([]byte("lh"), block.Hash, nil)
		Handle(err)
		chain.LastHash = block.Hash
	}
}

// GetBlock retrieves a block by its hash
func (chain *Blockchain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	data, err := chain.Database.Get(blockHash, nil)
	if err != nil {
		return block, err
	}

	block = *Deserialize(data)

	return block, nil
}

// GetBestHeight returns the height of the latest block in the chain
func (chain *Blockchain) GetBestHeight() int {
	var lastBlock Block

	data, err := chain.Database.Get(chain.LastHash, nil)
	Handle(err)
	lastBlock = *Deserialize(data)

	return lastBlock.Height
}

// GetLastBlock returns the last block in the blockchain
func (chain *Blockchain) GetLastBlock() *Block {
	var lastBlock Block

	data, err := chain.Database.Get(chain.LastHash, nil)
	Handle(err)
	lastBlock = *Deserialize(data)

	return &lastBlock
}

// GetBlockHashes returns a list of block hashes in the blockchain
func (chain *Blockchain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	iter := chain.Iterator()

	for {
		block := iter.Next()
		blocks = append(blocks, block.Hash)

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return blocks
}

// FindTransaction finds a transaction by its ID
func (chain *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction not found")
}

// SignTransaction signs inputs of a transaction
func (chain *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := chain.FindTransaction(in.ID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

// VerifyTransaction verifies transaction inputs signatures
func (chain *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := chain.FindTransaction(in.ID)
		if err != nil {
			return false
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

// FindUnspentTransactions returns a list of transactions containing unspent outputs
func (chain *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTxs []Transaction
	spentTXOs := make(map[string][]int)
	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
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

// FindUTXO finds and returns all unspent transaction outputs for a specific public key
func (chain *Blockchain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := chain.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// FindAllUTXO finds all unspent transaction outputs and returns them indexed by transaction ID
func (chain *Blockchain) FindAllUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					inTxID := hex.EncodeToString(in.ID)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return UTXO
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (chain *Blockchain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}

// Iterator returns a BlockchainIterator
func (chain *Blockchain) Iterator() *BlockchainIterator {
	iter := &BlockchainIterator{chain.LastHash, chain.Database}
	return iter
}

// Next returns the next block in the iteration
func (iter *BlockchainIterator) Next() *Block {
	var block *Block

	data, err := iter.Database.Get(iter.CurrentHash, nil)
	Handle(err)

	block = Deserialize(data)
	iter.CurrentHash = block.PrevHash

	return block
}

// Close closes the blockchain database
func (chain *Blockchain) Close() {
	chain.Database.Close()
}
