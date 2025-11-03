package blockchain

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	utxoPrefix   = []byte("utxo-")
	prefixLength = len(utxoPrefix)
)

// UTXOSet represents the set of UTXOs (Unspent Transaction Outputs)
// Similar to Bitcoin, maintains a cache of unspent outputs
type UTXOSet struct {
	Blockchain *Blockchain
}

// FindSpendableOutputs finds and returns unspent outputs that can be used
func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	accumulated := 0

	db := u.Blockchain.Database

	iter := db.NewIterator(util.BytesPrefix(utxoPrefix), nil)
	defer iter.Release()

	for iter.Next() {
		k := iter.Key()
		v := iter.Value()

		k = bytes.TrimPrefix(k, utxoPrefix)
		txID := hex.EncodeToString(k)
		outs := DeserializeOutputs(v)

		for outIdx, out := range outs.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)
			}
		}
	}

	if err := iter.Error(); err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOuts
}

// FindUTXO finds all UTXOs for a public key
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput

	db := u.Blockchain.Database

	iter := db.NewIterator(util.BytesPrefix(utxoPrefix), nil)
	defer iter.Release()

	for iter.Next() {
		v := iter.Value()
		outs := DeserializeOutputs(v)

		for _, out := range outs.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	if err := iter.Error(); err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// CountTransactions returns the number of transactions in the UTXO set
func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.Database
	counter := 0

	iter := db.NewIterator(util.BytesPrefix(utxoPrefix), nil)
	defer iter.Release()

	for iter.Next() {
		counter++
	}

	if err := iter.Error(); err != nil {
		log.Panic(err)
	}

	return counter
}

// Reindex rebuilds the UTXO set
func (u UTXOSet) Reindex() {
	db := u.Blockchain.Database

	u.DeleteByPrefix(utxoPrefix)

	UTXO := u.Blockchain.FindAllUTXO()

	for txId, outs := range UTXO {
		key, err := hex.DecodeString(txId)
		if err != nil {
			log.Panic(err)
		}
		key = append(utxoPrefix, key...)

		err = db.Put(key, outs.Serialize(), nil)
		if err != nil {
			log.Panic(err)
		}
	}
}

// Update updates the UTXO set with the block's transactions
func (u *UTXOSet) Update(block *Block) {
	db := u.Blockchain.Database

	for _, tx := range block.Transactions {
		if tx.IsCoinbase() == false {
			for _, in := range tx.Inputs {
				updatedOuts := TXOutputs{}
				inID := append(utxoPrefix, in.ID...)
				
				v, err := db.Get(inID, nil)
				if err != nil {
					log.Panic(err)
				}

				outs := DeserializeOutputs(v)

				for outIdx, out := range outs.Outputs {
					if outIdx != in.Out {
						updatedOuts.Outputs = append(updatedOuts.Outputs, out)
					}
				}

				if len(updatedOuts.Outputs) == 0 {
					if err := db.Delete(inID, nil); err != nil {
						log.Panic(err)
					}
				} else {
					if err := db.Put(inID, updatedOuts.Serialize(), nil); err != nil {
						log.Panic(err)
					}
				}
			}
		}

		newOutputs := TXOutputs{}
		for _, out := range tx.Outputs {
			newOutputs.Outputs = append(newOutputs.Outputs, out)
		}

		txID := append(utxoPrefix, tx.ID...)
		if err := db.Put(txID, newOutputs.Serialize(), nil); err != nil {
			log.Panic(err)
		}
	}
}

// DeleteByPrefix deletes all items with a specific prefix
func (u *UTXOSet) DeleteByPrefix(prefix []byte) {
	db := u.Blockchain.Database

	iter := db.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()

	keysToDelete := make([][]byte, 0)
	for iter.Next() {
		key := make([]byte, len(iter.Key()))
		copy(key, iter.Key())
		keysToDelete = append(keysToDelete, key)
	}

	if err := iter.Error(); err != nil {
		log.Panic(err)
	}

	for _, key := range keysToDelete {
		if err := db.Delete(key, nil); err != nil {
			log.Panic(err)
		}
	}
}

