# Blockchain Network Documentation

## Overview

This blockchain implementation includes a peer-to-peer (P2P) network layer that allows multiple nodes to communicate, synchronize, and maintain a distributed ledger.

## Network Architecture

### Components

1. **Server** (`internal/network/server.go`)
   - TCP-based server for handling incoming connections
   - Manages peer connections and message routing
   - Coordinates blockchain synchronization

2. **Peer Management** (`internal/network/peer.go`)
   - Maintains list of connected peers
   - Handles peer connection lifecycle
   - Thread-safe peer operations

3. **Protocol** (`internal/network/protocol.go`)
   - Defines message types and formats
   - Serialization/deserialization using gob encoding
   - Command routing

### Message Types

| Command | Description |
|---------|-------------|
| `version` | Handshake message exchanged when peers connect |
| `getblocks` | Request blockchain hashes from a peer |
| `inv` | Inventory message listing available blocks/transactions |
| `getdata` | Request specific block or transaction |
| `block` | Send block data |
| `tx` | Send transaction data |
| `addr` | Share peer addresses |
| `ping/pong` | Keep-alive messages |

## Network Flow

### 1. Node Startup

```
Node A starts → Listens on port → Connects to seed nodes → Exchanges version
```

### 2. Blockchain Synchronization

```
Node A ← version (height: 10) ← Node B (seed)
Node A → getblocks → Node B
Node A ← inv (block hashes) ← Node B
Node A → getdata (block hash) → Node B
Node A ← block (data) ← Node B
```

### 3. Transaction Broadcasting

```
Node A creates transaction → broadcasts to all peers
Node B receives transaction → adds to mempool
Node B (miner) mines block → broadcasts new block
Node A receives new block → validates → adds to chain
```

## Running the Network

### Local Testing

#### Method 1: Multiple Terminals

**Terminal 1 - Seed Node + Miner:**
```bash
./build/blockchain startnode -port 3000 -miner <ADDRESS>
```

**Terminal 2 - Mining Node:**
```bash
./build/blockchain startnode -port 3001 -miner <ADDRESS>
```

**Terminal 3 - Regular Node:**
```bash
./build/blockchain startnode -port 3002
```

#### Method 2: Network Demo Script

```bash
make network-demo
# Follow the instructions to start nodes in separate terminals
```

### Docker Testing

#### Quick Start

```bash
make docker-build    # Build images
make docker-up       # Start network (4 nodes)
make docker-logs     # View logs
```

#### Full Test

```bash
make docker-test     # Automated test script
```

#### Manual Docker Commands

```bash
# Start network
docker-compose up -d

# View logs
docker-compose logs -f node-seed
docker-compose logs -f node-miner1

# Execute commands in containers
docker exec -it blockchain-seed /app/blockchain listaddresses
docker exec -it blockchain-miner1 /app/blockchain getbalance -address <ADDRESS>
docker exec -it blockchain-seed /app/blockchain printchain

# Stop network
docker-compose down

# Clean up (removes data)
docker-compose down -v
```

## Network Configuration

### Docker Compose Network

The `docker-compose.yml` creates a network with:

- **Seed Node** (172.20.0.2:3000) - Non-mining, seed node
- **Miner 1** (172.20.0.3:3001) - Mining node
- **Miner 2** (172.20.0.4:3002) - Mining node
- **Regular Node** (172.20.0.5:3003) - Non-mining node

### Port Mapping

| Container | Internal Port | External Port |
|-----------|--------------|---------------|
| node-seed | 3000 | 3000 |
| node-miner1 | 3001 | 3001 |
| node-miner2 | 3002 | 3002 |
| node-regular | 3003 | 3003 |

## CLI Commands

### Start Node

```bash
# Start mining node
./build/blockchain startnode -port 3000 -miner <WALLET_ADDRESS>

# Start regular node
./build/blockchain startnode -port 3000
```

**Flags:**
- `-port` - Port to listen on (default: 3000)
- `-miner` - Optional wallet address for mining rewards

### Add Peer

```bash
./build/blockchain addpeer -address localhost:3001
```

### List Peers

```bash
./build/blockchain peers
```

## Mining in Network Mode

When a node is started with the `-miner` flag:

1. Node receives transactions via P2P
2. Transactions are added to mempool
3. When mempool reaches threshold (2+ transactions), mining starts
4. New block is mined and broadcast to all peers
5. Peers validate and add block to their chains

## Troubleshooting

### Node Not Connecting

**Issue:** Node can't connect to seed node

**Solution:**
```bash
# Check if seed node is running
nc -zv localhost 3000

# Check logs
docker logs blockchain-seed
```

### Blockchain Out of Sync

**Issue:** Node has different blockchain than peers

**Solution:**
```bash
# Reindex UTXO set
./build/blockchain reindexutxo

# Or restart node (will sync from seed)
docker-compose restart node-miner1
```

### Port Already in Use

**Issue:** `address already in use`

**Solution:**
```bash
# Find process using port
lsof -i :3000

# Kill process or use different port
./build/blockchain startnode -port 3005
```

## Network Security Considerations

⚠️ **This is an educational implementation. For production:**

1. **Add TLS/SSL** - Encrypt peer connections
2. **Authentication** - Verify peer identities
3. **DDoS Protection** - Rate limiting, connection limits
4. **Sybil Attack Prevention** - Proof of work for peer discovery
5. **Eclipse Attack Prevention** - Diverse peer selection
6. **Input Validation** - Sanitize all network messages

## Performance Tuning

### Mempool Size

Edit `internal/network/server.go`:
```go
// Mine when mempool has N transactions
if len(memoryPool) >= 2 {  // Change this value
    s.mineTransactions()
}
```

### Block Propagation

Optimize by:
- Using compact block relay (send only transaction IDs)
- Implementing bloom filters
- Adding peer scoring

### Network Topology

Current: Star topology (all nodes → seed)

For better resilience:
- Implement peer discovery
- Allow mesh topology
- Add multiple seed nodes

## Bitcoin Protocol Comparison

| Feature | Bitcoin | This Implementation |
|---------|---------|---------------------|
| Protocol | Custom binary | Gob encoding |
| Discovery | DNS seeds, peers.dat | Static seed node |
| Connection | Persistent | Per-message |
| Mempool | Yes | Yes (in-memory) |
| Block relay | Compact blocks | Full blocks |
| Transaction relay | Yes | Yes |
| SPV support | Yes | No |

## Next Steps

To make this network more robust:

1. **Persistent Connections** - Keep peer connections open
2. **Peer Discovery** - Automatic peer finding
3. **Block Headers First** - Faster sync
4. **SPV Support** - Lightweight clients
5. **Network Statistics** - Monitor peer health
6. **Automatic Reconnection** - Handle network failures

## References

- [Bitcoin P2P Protocol](https://en.bitcoin.it/wiki/Protocol_documentation)
- [Bitcoin Developer Guide - P2P Network](https://developer.bitcoin.org/devguide/p2p_network.html)
- [Mastering Bitcoin - Chapter 8](https://github.com/bitcoinbook/bitcoinbook/blob/develop/ch08.asciidoc)

