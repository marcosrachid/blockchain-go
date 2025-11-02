# Architecture

This document describes the architecture and code organization of the blockchain project, following the [golang-standards/project-layout](https://github.com/golang-standards/project-layout).

## Project Structure

```
blockchain-go/
├── cmd/                          # Main applications
│   └── blockchain/               # Blockchain CLI application
│       └── main.go              # Entry point and CLI implementation
├── internal/                     # Private application code
│   └── blockchain/               # Blockchain core implementation
│       ├── base58.go            # Base58 encoding/decoding
│       ├── block.go             # Block structure and operations
│       ├── blockchain.go        # Blockchain management
│       ├── merkle.go            # Merkle tree implementation
│       ├── proof.go             # Proof of Work algorithm
│       ├── transaction.go       # Transaction handling
│       ├── utxo.go              # UTXO set management
│       ├── utils.go             # Utility functions
│       └── wallet.go            # Wallet management
├── build/                        # Build artifacts (generated)
├── docs/                         # Documentation
│   ├── BITCOIN_COMPARISON.md   # Comparison with Bitcoin
│   ├── IMPROVEMENTS.md          # Suggested improvements
│   ├── README.md                # Original README (Portuguese)
│   └── TUTORIAL.md              # Tutorial (Portuguese)
├── scripts/                      # Build, install, and demo scripts
│   └── demo.sh                  # Demo script
├── tmp/                          # Runtime data (generated)
│   ├── blocks/                  # BadgerDB blockchain data
│   └── wallets.dat              # Wallet storage
├── .gitignore                   # Git ignore rules
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── Makefile                     # Build automation
├── ARCHITECTURE.md              # This file
└── README.md                    # Project README (English)
```

## Directory Purposes

### `/cmd`
Main applications for this project. The directory name matches the binary name (blockchain).

- Contains the CLI implementation
- Keeps the main application entry point simple
- All CLI commands are implemented here

### `/internal`
Private application and library code. This layout pattern is enforced by the Go compiler.

- Code here cannot be imported by other projects
- Contains the core blockchain implementation
- All blockchain logic, cryptography, and data structures

### `/build`
Build artifacts directory.

- Contains compiled binaries
- Gitignored (should not be committed)
- Created by `make build`

### `/docs`
Design and user documents.

- Bitcoin comparison analysis
- Improvement suggestions
- Tutorials and guides
- Original Portuguese documentation

### `/scripts`
Scripts for build, install, analysis, and other operations.

- Demo script showing the blockchain in action
- Keeps the root Makefile simple

## Core Components

### Block (`internal/blockchain/block.go`)
```go
type Block struct {
    Timestamp    int64
    Hash         []byte
    Transactions []*Transaction
    PrevHash     []byte
    Nonce        int
    Height       int
}
```
- Represents a single block in the blockchain
- Contains multiple transactions
- Includes Merkle tree hash of transactions
- Implements Proof of Work

### Blockchain (`internal/blockchain/blockchain.go`)
```go
type Blockchain struct {
    LastHash []byte
    Database *badger.DB
}
```
- Manages the entire blockchain
- Provides persistence with BadgerDB
- Handles block mining and validation
- Manages transaction verification

### Transaction (`internal/blockchain/transaction.go`)
```go
type Transaction struct {
    ID      []byte
    Inputs  []TXInput
    Outputs []TXOutput
}
```
- Bitcoin-like transaction model
- Multiple inputs and outputs
- Digital signatures (ECDSA)
- Coinbase transactions for mining rewards

### Wallet (`internal/blockchain/wallet.go`)
```go
type Wallet struct {
    PrivateKey ecdsa.PrivateKey
    PublicKey  []byte
}
```
- ECDSA key pair management
- Bitcoin-like address generation
- Base58 encoding
- Address validation with checksum

### UTXO Set (`internal/blockchain/utxo.go`)
```go
type UTXOSet struct {
    Blockchain *Blockchain
}
```
- Caches unspent transaction outputs
- Improves transaction validation performance
- Similar to Bitcoin's UTXO database

### Merkle Tree (`internal/blockchain/merkle.go`)
```go
type MerkleTree struct {
    RootNode *MerkleNode
}
```
- Efficient transaction verification
- Hash tree of all transactions in a block
- Identical to Bitcoin's implementation

### Proof of Work (`internal/blockchain/proof.go`)
```go
type ProofOfWork struct {
    Block  *Block
    Target *big.Int
}
```
- Mining algorithm
- Adjustable difficulty
- SHA256 hashing
- Nonce finding

## Data Flow

### Creating a Transaction

1. User calls CLI: `./build/blockchain send -from A -to B -amount 10`
2. CLI validates addresses
3. Wallet loaded from file
4. UTXO set queried for spendable outputs
5. Transaction created with inputs/outputs
6. Transaction signed with sender's private key
7. Coinbase transaction created (mining reward)
8. Block mined with both transactions
9. Block added to blockchain
10. UTXO set updated

### Mining a Block

1. Transactions collected
2. Merkle tree created from transactions
3. Block header constructed
4. Proof of Work algorithm runs:
   - Increment nonce
   - Hash block header
   - Check if hash < target
   - Repeat until valid hash found
5. Block stored in database
6. Blockchain updated

### Address Generation

1. Generate ECDSA key pair (P256 curve)
2. Hash public key:
   - SHA256(publicKey)
   - RIPEMD160(hash)
3. Add version byte (0x00)
4. Calculate checksum:
   - SHA256(SHA256(versioned hash))
   - Take first 4 bytes
5. Encode in Base58
6. Result: Bitcoin-like address

## Key Design Decisions

### 1. Internal Package
- All blockchain code in `/internal` prevents external imports
- Enforces encapsulation
- Clear API boundary

### 2. Single Binary
- CLI included in main application
- Simpler deployment
- No separate libraries needed

### 3. BadgerDB for Storage
- Embedded key-value store
- No external database required
- Fast and reliable
- Used by many Go projects

### 4. ECDSA P256
- Standard Go crypto library
- Well-tested and secure
- Bitcoin uses secp256k1, but P256 serves well for education

### 5. English Code, Portuguese Docs
- Code and comments in English (industry standard)
- Documentation available in Portuguese
- Makes code accessible internationally

## Build Process

```bash
# Install dependencies
make deps

# Build binary
make build
# Output: ./build/blockchain

# Run tests
make test

# Clean
make clean
```

## Testing Strategy

The project can be tested at multiple levels:

1. **Unit Tests**: Test individual components
2. **Integration Tests**: Test component interaction
3. **Manual Testing**: Use demo script
4. **CLI Testing**: Test through command line

## Performance Considerations

### UTXO Set Caching
- Keeps unspent outputs in memory/database
- Faster transaction validation
- Must be rebuilt if corrupted

### Database Indexing
- Blocks indexed by hash
- Transactions indexed in UTXO set
- Fast lookups

### Merkle Tree
- O(log n) verification
- Only need branch, not full tree
- Efficient for large block sizes

## Security Considerations

### 1. Private Key Storage
- Stored encrypted in wallet file
- Never transmitted
- Should use proper encryption in production

### 2. Transaction Validation
- Signature verification (ECDSA)
- UTXO validation
- Double-spend prevention

### 3. Proof of Work
- Computational cost prevents spam
- Makes chain immutable
- Higher difficulty = more security

## Future Enhancements

See `docs/IMPROVEMENTS.md` for detailed suggestions:
- Dynamic difficulty adjustment
- Mining reward halving
- Transaction fees
- P2P networking
- Multi-signature support
- Time-locks
- REST API
- Web interface

## References

- [Bitcoin Whitepaper](https://bitcoin.org/bitcoin.pdf)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- [Mastering Bitcoin](https://github.com/bitcoinbook/bitcoinbook)
- [BadgerDB](https://github.com/dgraph-io/badger)

---

**Note**: This is an educational implementation. Do not use in production for real value.

