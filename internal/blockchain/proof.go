package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"
	"time"
)

// Difficulty is now defined in config.go

type ProofOfWork struct {
	Block      *Block
	Target     *big.Int
	Difficulty int
}

func NewProof(b *Block) *ProofOfWork {
	return NewProofWithDifficulty(b, Difficulty)
}

func NewProofWithDifficulty(b *Block, difficulty int) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	// Ensure block's difficulty field is set
	if b.Difficulty == 0 {
		b.Difficulty = difficulty
	}

	pow := &ProofOfWork{b, target, difficulty}
	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	// Use stored MerkleRoot instead of recalculating to ensure consistency
	// across serialization/deserialization
	nonceBytes := toHex(int64(nonce))
	diffBytes := toHex(int64(pow.Block.Difficulty))
	timeBytes := toHex(pow.Block.Timestamp)

	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.MerkleRoot, // Use stored Merkle Root
			nonceBytes,
			diffBytes,
			timeBytes,
		},
		[]byte{},
	)
	return data
}

// DebugInitData prints each component for debugging
func (pow *ProofOfWork) DebugInitData(nonce int) {
	nonceBytes := toHex(int64(nonce))
	diffBytes := toHex(int64(pow.Block.Difficulty))
	timeBytes := toHex(pow.Block.Timestamp)

	log.Printf("🔍 InitData components:")
	log.Printf("   PrevHash: %x", pow.Block.PrevHash)
	log.Printf("   MerkleRoot (stored): %x", pow.Block.MerkleRoot)
	log.Printf("   Nonce: %d (%x)", nonce, nonceBytes)
	log.Printf("   Difficulty: %d (%x)", pow.Block.Difficulty, diffBytes)
	log.Printf("   Timestamp: %d (%x)", pow.Block.Timestamp, timeBytes)
}

func (pow *ProofOfWork) Run() (int, []byte) {
	return pow.RunWithInterrupt(nil)
}

func (pow *ProofOfWork) RunWithInterrupt(interrupt <-chan bool) (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0
	checkInterval := 10000    // Check for interrupts every 10k iterations
	logInterval := 100000     // Log progress every 100k hashes
	timestampInterval := 1000 // Update timestamp every 1k iterations

	for nonce < math.MaxInt64 {
		// Update timestamp periodically (every ~1k hashes) to keep it current
		// Uses UTC to ensure consistency across different timezones
		if nonce%timestampInterval == 0 {
			pow.Block.Timestamp = time.Now().UTC().Unix()
		}

		// Check for interrupt signal periodically
		if interrupt != nil && nonce%checkInterval == 0 {
			select {
			case <-interrupt:
				// Mining interrupted - return zero values
				log.Printf("⛏️  Mining interrupted at nonce %d", nonce)
				return 0, nil
			default:
				// Continue mining
			}
		}

		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			// Found valid hash - DO NOT update timestamp as it would invalidate the hash!
			log.Printf("✅ Found valid hash: %x at nonce %d", hash, nonce)
			// Debug: show what data was used
			log.Printf("🔍 MINING: Raw InitData (len=%d): %x", len(data), data)
			pow.DebugInitData(nonce)
			break
		}

		// Log progress periodically
		if nonce > 0 && nonce%logInterval == 0 {
			log.Printf("⛏️  Mining... nonce: %d", nonce)
		}

		nonce++
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func toHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
