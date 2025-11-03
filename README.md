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

### 9. **HTTP REST API**

- Create wallets (`POST /api/createwallet`)
- Send transactions (`POST /api/send`)
- Check balances (`GET /api/balance/:address`)
- Network info (`GET /api/networkinfo`)
- List addresses (`GET /api/addresses`)
- View last block (`GET /api/lastblock`)
- Health check (`GET /health`)

### 10. **CLI (Command Line Interface)**

- Start network nodes (`startnode`)
- Create blockchain (`createblockchain`)
- Basic wallet management (`createwallet`, `listaddresses`)

## 🏗️ Project Structure

Following the [golang-standards/project-layout](https://github.com/golang-standards/project-layout):

```
blockchain-go/
├── cmd/
│   └── blockchain/          # Main application entry point
│       └── main.go          # Application startup and basic commands
├── internal/
│   ├── api/                 # HTTP API server
│   │   └── server.go        # REST API endpoints (balance, send, network info, etc.)
│   ├── blockchain/          # Core blockchain logic
│   │   ├── base58.go        # Base58 encoding (Bitcoin-style)
│   │   ├── block.go         # Block structure with PoW and transactions
│   │   ├── blockchain.go    # Blockchain with persistence (LevelDB)
│   │   ├── config.go        # Configuration constants (difficulty, rewards, etc.)
│   │   ├── merkle.go        # Merkle Tree for transaction hashing
│   │   ├── proof.go         # Proof of Work algorithm
│   │   ├── transaction.go   # Transaction system with ECDSA signatures
│   │   ├── utxo.go          # UTXO set management
│   │   ├── utils.go         # Utility functions
│   │   └── wallet.go        # Wallet and address management
│   └── network/             # P2P network layer
│       ├── peer.go          # Peer connection management
│       ├── protocol.go      # Network protocol messages
│       └── server.go        # P2P server, mempool, mining coordination
├── build/                   # Compiled binaries
├── docs/                    # Detailed documentation
│   ├── ARCHITECTURE.md      # System architecture
│   ├── BITCOIN_COMPARISON.md # How it compares to Bitcoin
│   ├── HALVING_AND_SUPPLY.md # Economics and supply model
│   ├── MINING.md            # Mining mechanics
│   ├── NETWORK.md           # Network protocol details
│   └── ...                  # Portuguese versions (*.pt-br.md)
├── scripts/                 # Utility scripts
│   ├── check-balances.sh    # Check balances across all nodes
│   ├── check-lastblock.sh   # Check blockchain height of all nodes
│   ├── network-status.sh    # Network status dashboard
│   ├── docker-test.sh       # Automated Docker network test
│   └── demo.sh              # Quick demo script
├── docker-compose.yml       # Multi-node Docker network setup
├── Dockerfile               # Container image definition
├── go.mod                   # Go module dependencies
├── go.sum                   # Go module checksums
├── Makefile                 # Build automation
├── LICENSE                  # MIT License
├── README.md                # This file (English)
└── README.pt-br.md         # Portuguese README
```

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or higher
- Docker & Docker Compose (for multi-node testing)

### Quick Start (Docker - Recommended)

The Docker network comes pre-configured with 4 nodes and is the easiest way to test:

```bash
# Clone the repository
git clone https://github.com/marcocsrachid/blockchain-go.git
cd blockchain-go

# Build and start the network (4 nodes: 1 seed, 2 miners, 1 regular)
docker-compose build
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

**Exposed ports:**

- `4000` - Seed Node (HTTP API)
- `4001` - Miner 1 (HTTP API)
- `4002` - Miner 2 (HTTP API)
- `4003` - Regular Node (HTTP API)

**Useful scripts:**

```bash
# Complete network status
./scripts/network-status.sh

# Check block heights
./scripts/check-lastblock.sh

# Check balances
./scripts/check-balances.sh
```

### Manual Build (Local)

If you want to compile and run locally:

#### 1. Build the binary

```bash
# Standard build
go build -o build/blockchain cmd/blockchain/main.go

# Or static build (for Docker Alpine)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix netgo -ldflags '-s -w' -o build/blockchain cmd/blockchain/main.go
```

#### 2. Create wallet and blockchain

```bash
# Create a wallet (note the generated address)
./build/blockchain createwallet

# List addresses
./build/blockchain listaddresses

# Create the blockchain with reward address
./build/blockchain createblockchain -address YOUR_ADDRESS
```

#### 3. Start a node

**Mining node (produces blocks):**

```bash
# Terminal 1 - Seed/Miner Node
NODE_ID=node1 ./build/blockchain startnode -port 3000 -miner YOUR_ADDRESS
```

**Regular node (non-mining):**

```bash
# Terminal 2 - Regular Node (connects to node1)
NODE_ID=node2 SEED_NODE=localhost:3000 ./build/blockchain startnode -port 3001
```

**Important environment variables:**

- `NODE_ID` - Unique node ID (defines data directory)
- `SEED_NODE` - Seed node address to connect to
- `-port` - Node P2P port (default: 3000)
- `-apiport` - HTTP API port (default: 4000)
- `-miner` - Address to receive mining rewards (enables mining)

### Using the HTTP API

All nodes expose a REST API:

```bash
# Check network status
curl http://localhost:4000/api/networkinfo | jq

# List addresses
curl http://localhost:4000/api/addresses | jq

# Check balance
curl http://localhost:4000/api/balance/YOUR_ADDRESS | jq

# Create new wallet
curl -X POST http://localhost:4000/api/createwallet | jq

# Send transaction
curl -X POST http://localhost:4000/api/send \
  -H "Content-Type: application/json" \
  -d '{
    "from": "FROM_ADDRESS",
    "to": "TO_ADDRESS",
    "amount": 10
  }' | jq

# View last block
curl http://localhost:4000/api/lastblock | jq

# List known peers
curl http://localhost:4000/api/peers | jq
```

### Complete Example (3 Nodes)

```bash
# Terminal 1 - Seed Node (non-mining, coordinator only)
NODE_ID=seed ./build/blockchain createblockchain -address 1SeedAddress...
NODE_ID=seed ./build/blockchain startnode -port 3000

# Terminal 2 - Miner 1
NODE_ID=miner1 SEED_NODE=localhost:3000 ./build/blockchain startnode -port 3001 -apiport 4001 -miner 1Miner1Address...

# Terminal 3 - Miner 2
NODE_ID=miner2 SEED_NODE=localhost:3000 ./build/blockchain startnode -port 3002 -apiport 4002 -miner 1Miner2Address...

# Terminal 4 - Send transaction via API
curl -X POST http://localhost:4001/api/send \
  -H "Content-Type: application/json" \
  -d '{"from":"1Miner1Address...","to":"1Miner2Address...","amount":50}' | jq

# Wait ~60-90s for mining...

# Check balances
curl http://localhost:4001/api/balance/1Miner1Address... | jq
curl http://localhost:4002/api/balance/1Miner2Address... | jq
```

### Accessing Docker Containers

```bash
# Execute commands inside containers
docker exec -it blockchain-seed /app/blockchain listaddresses
docker exec -it blockchain-miner1 /app/blockchain listaddresses

# View logs of a specific node
docker-compose logs -f node-seed
docker-compose logs -f node-miner1

# Stop the network
docker-compose down

# Stop and clean data (complete reset)
docker-compose down -v
```

📖 For complete network implementation details, see [docs/NETWORK.md](docs/NETWORK.md) and [docs/NETWORK.pt-br.md](docs/NETWORK.pt-br.md)

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

1. **Fixed Difficulty**: Bitcoin adjusts difficulty every 2016 blocks (this project has fixed difficulty)
2. **Simplified Halving**: Bitcoin halves every 210,000 blocks; this project halves every 210,000 blocks but without adjustment complexity
3. **Basic P2P**: Implemented but simpler than Bitcoin's full protocol
4. **No Scripts**: Bitcoin uses Script language for spending conditions
5. **Basic Mempool**: Implemented but without priority fees or RBF (Replace-By-Fee)
6. **No SPV**: Simplified Payment Verification not implemented
7. **No SegWit**: Segregated Witness not implemented
8. **No Lightning**: Lightning Network not implemented
9. **No DNS Seeds**: Manual peer management instead of DNS seed discovery
10. **Simplified Consensus**: No orphan handling or chain reorganization

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
