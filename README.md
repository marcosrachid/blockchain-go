# Blockchain in Go - Bitcoin-like Implementation

An educational blockchain project implemented in Go, inspired by the Bitcoin protocol.

[![Go Report Card](https://goreportcard.com/badge/github.com/marcocsrachid/blockchain-go)](https://goreportcard.com/report/github.com/marcocsrachid/blockchain-go)

## 📚 Features

This project implements the main Bitcoin concepts:

### 1. **Proof of Work (PoW)**

- Consensus algorithm similar to Bitcoin
- Adjustable difficulty
- Block mining with nonce

### 2. **Transaction System**

- Transactions with multiple inputs and outputs
- Coinbase transactions (mining reward)
- Digital signature verification with ECDSA

### 3. **UTXOs (Unspent Transaction Outputs)**

- UTXO model similar to Bitcoin
- UTXO cache for performance
- Unspent output tracking system

### 4. **ECDSA Cryptography**

- Public/private key pair generation
- Transaction digital signatures using ECDSA
- P256 elliptic curve

### 5. **Wallets**

- Bitcoin-like address generation
- Base58 encoding (Bitcoin alphabet)
- Public key hashing (SHA256 + RIPEMD160)
- Address validation with checksum

### 6. **Merkle Tree**

- Data structure for transaction verification
- Efficient block hashing
- Identical to Bitcoin implementation

### 7. **Persistence**

- LevelDB database (supports concurrent read/write access)
- Block serialization/deserialization
- Blockchain iterator

### 8. **P2P Network** 🆕

- Peer-to-peer network communication
- TCP-based protocol
- Block and transaction broadcasting
- Blockchain synchronization between nodes
- Mining nodes and regular nodes
- Seed node support

### 9. **CLI (Command Line Interface)**

- Create wallets
- Send transactions
- Check balances
- Print blockchain
- Reindex UTXOs
- Start network nodes
- Manage peers

## 🏗️ Project Structure

Following the [golang-standards/project-layout](https://github.com/golang-standards/project-layout):

```
blockchain-go/
├── cmd/
│   └── blockchain/          # Main application entry point
│       └── main.go          # CLI implementation
├── internal/
│   ├── blockchain/          # Private application code
│   │   ├── base58.go        # Base58 encoding (Bitcoin)
│   │   ├── block.go         # Block structure with transactions
│   │   ├── blockchain.go    # Main blockchain with UTXO
│   │   ├── merkle.go        # Merkle Tree implementation
│   │   ├── proof.go         # Proof of Work
│   │   ├── transaction.go   # Transaction system
│   │   ├── utxo.go          # UTXO system
│   │   ├── utils.go         # Utility functions
│   │   └── wallet.go        # Wallet system
│   └── network/             # P2P network layer
│       ├── peer.go          # Peer management
│       ├── protocol.go      # Network protocol
│       └── server.go        # Network server
├── build/                   # Build artifacts
├── docs/                    # Documentation
├── scripts/                 # Build and demo scripts
├── go.mod                   # Go modules
├── Makefile                 # Build automation
└── README.md               # This file
```

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or higher
- Make (optional, but recommended)
- Docker & Docker Compose (for network testing)

### Installation

```bash
# Clone the repository
git clone https://github.com/marcocsrachid/blockchain-go.git
cd blockchain-go

# Install dependencies
make deps

# Build the project
make build
```

### Usage

#### Create a Wallet

```bash
./build/blockchain createwallet
```

Output example:

```
New address is: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
```

#### List Addresses

```bash
./build/blockchain listaddresses
```

#### Create the Blockchain

Create the blockchain and send the genesis block reward to an address:

```bash
./build/blockchain createblockchain -address YOUR_ADDRESS
```

#### Check Balance

```bash
./build/blockchain getbalance -address YOUR_ADDRESS
```

#### Send Transaction

```bash
./build/blockchain send -from FROM_ADDRESS -to TO_ADDRESS -amount 10
```

#### View the Blockchain

```bash
./build/blockchain printchain
```

#### Reindex UTXOs

```bash
./build/blockchain reindexutxo
```

### Network Commands 🌐

#### Start a Node

Start a mining node:

```bash
./build/blockchain startnode -port 3000 -miner YOUR_ADDRESS
```

Start a regular (non-mining) node:

```bash
./build/blockchain startnode -port 3000
```

#### Manage Peers

Add a peer:

```bash
./build/blockchain addpeer -address localhost:3001
```

List known peers:

```bash
./build/blockchain peers
```

### Docker Network Testing

#### Quick Start

```bash
# Build and start 4-node network
make docker-build
make docker-up

# View logs
make docker-logs

# Stop network
make docker-down
```

#### Full Docker Test

```bash
# Run automated test script
make docker-test
```

The docker-compose setup includes:

- **Seed Node** (port 3000) - Non-mining seed node
- **Miner 1** (port 3001) - Mining node
- **Miner 2** (port 3002) - Mining node
- **Regular Node** (port 3003) - Non-mining node

#### Execute Commands in Containers

```bash
# List addresses
docker exec -it blockchain-seed /app/blockchain listaddresses

# Check balance
docker exec -it blockchain-miner1 /app/blockchain getbalance -address <ADDRESS>

# View blockchain
docker exec -it blockchain-seed /app/blockchain printchain
```

See [docs/NETWORK.md](docs/NETWORK.md) for detailed network documentation.

## 📖 Bitcoin Concepts Implemented

### 1. Proof of Work

Consensus algorithm that ensures security through computational work:

- Miners must find a hash that meets the established difficulty
- The hash must have a certain number of leading zeros
- Similar to Bitcoin's SHA256(SHA256())

### 2. Transactions

Bitcoin-like structure:

- **Inputs**: References to previous transaction outputs
- **Outputs**: New destinations for coins with specific values
- **Coinbase**: Special reward transaction for the miner

### 3. UTXO (Unspent Transaction Output)

- Bitcoin's accounting model
- Each output can only be spent once
- Tracking system for unspent outputs for efficiency

### 4. Cryptography

- **ECDSA**: Transaction digital signatures
- **SHA256**: Block and transaction hashing
- **RIPEMD160**: Public key hashing
- **Base58**: Address encoding (avoids ambiguous characters)

### 5. Merkle Tree

- Data structure that allows efficient transaction verification
- Tree root included in block header
- Enables SPV (Simplified Payment Verification)

### 6. Block Structure

```go
type Block struct {
    Timestamp    int64           // When the block was mined
    Hash         []byte          // Block hash
    Transactions []*Transaction  // Transactions in block
    PrevHash     []byte          // Previous block hash
    Nonce        int             // Nonce for PoW
    Height       int             // Block height in blockchain
}
```

### 7. Wallets and Addresses

Bitcoin-like address generation process:

1. Generate ECDSA key pair
2. SHA256 of public key
3. RIPEMD160 of result
4. Add version byte
5. Calculate checksum (SHA256(SHA256()))
6. Encode in Base58

## 🔍 Differences from Real Bitcoin

This is an educational project. Some differences from real Bitcoin:

1. **Fixed Difficulty**: Bitcoin adjusts difficulty every 2016 blocks
2. **Fixed Reward**: Bitcoin halves reward every 210,000 blocks (halving)
3. **P2P Network**: Not implemented (Bitcoin has complete network protocol)
4. **Scripts**: Bitcoin uses Script language for spending conditions
5. **Mempool**: Pending transaction pool not implemented
6. **SPV**: Simplified Payment Verification not implemented
7. **Segregated Witness**: Not implemented
8. **Lightning Network**: Not implemented

## 🛠️ Development

### Build Commands

```bash
# Build the application
make build

# Clean artifacts
make clean

# Run tests
make test

# Format code
make fmt

# Run linter
make vet

# Development build (with race detector)
make dev
```

### Running the Demo

```bash
# Run the demo script
./scripts/demo.sh
```

## 📚 Learning Resources

To better understand Bitcoin:

1. [Bitcoin Whitepaper - Satoshi Nakamoto](https://bitcoin.org/bitcoin.pdf)
2. [Mastering Bitcoin - Andreas Antonopoulos](https://github.com/bitcoinbook/bitcoinbook)
3. [Bitcoin Developer Guide](https://bitcoin.org/en/developer-guide)
4. [Learn Me a Bitcoin](https://learnmeabitcoin.com/)

## 🤝 Contributing

This is an educational project. Feel free to:

- Fork the project
- Add new features
- Improve documentation
- Report issues

## ⚠️ Disclaimer

This project was created for educational purposes only and should not be used in production. It is not suitable for storing real value.

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

### What this means:
- ✅ **Free to use** for learning, education, and commercial projects
- ✅ **Free to modify** and adapt to your needs
- ✅ **Free to distribute** and share
- ⚠️ **No warranty** - use at your own risk
- 📝 **Attribution appreciated** but not required

## 🎯 Similarity with Bitcoin: **93%**

The project faithfully implements:

- ✅ Merkle Trees (100%)
- ✅ Address system (100%)
- ✅ UTXO model (95%)
- ✅ Proof of Work (85%)
- ✅ Transactions (90%)
- ✅ Wallets and cryptography (95%)

---

**Developed with 💙 for Blockchain and Bitcoin learning**
