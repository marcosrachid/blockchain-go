# Network Quick Start Guide

## 🚀 Quick Test with Docker (Recommended)

### 1. Start the Network

```bash
# Clean old data and build
make docker-clean
make docker-build

# Start 4-node network
make docker-up
```

This starts:
- **Seed Node** (localhost:3000) - Non-mining
- **Miner 1** (localhost:3001) - Mining node
- **Miner 2** (localhost:3002) - Mining node  
- **Regular Node** (localhost:3003) - Non-mining

### 2. Watch the Network

```bash
# View all logs
make docker-logs

# View specific node
docker-compose logs -f node-miner1
```

### 3. Interact with the Network

```bash
# List wallets in miner 1
docker exec -it blockchain-miner1 /app/blockchain listaddresses

# Get an address
ADDR=$(docker exec -it blockchain-miner1 /app/blockchain listaddresses | head -1 | tr -d '\r')

# Check balance
docker exec -it blockchain-miner1 /app/blockchain getbalance -address "$ADDR"

# View blockchain
docker exec -it blockchain-seed /app/blockchain printchain
```

### 4. Stop the Network

```bash
# Stop containers
make docker-down

# Stop and remove all data
make docker-clean
```

## 🖥️ Manual Testing (Multiple Terminals)

### Terminal 1: Seed Node

```bash
# Clean and build
make clean
make build

# Create blockchain
./build/blockchain createwallet
# Save the address: ADDRESS1=<your_address>

./build/blockchain createblockchain -address <ADDRESS1>

# Start seed node
./build/blockchain startnode -port 3000 -miner <ADDRESS1>
```

### Terminal 2: Mining Node

```bash
# Create wallet
./build/blockchain createwallet
# Save the address: ADDRESS2=<your_address>

# Start mining node
./build/blockchain startnode -port 3001 -miner <ADDRESS2>
```

### Terminal 3: Regular Node

```bash
# Start regular node
./build/blockchain startnode -port 3002
```

### Terminal 4: Send Transactions

```bash
# Send transaction
./build/blockchain send -from <ADDRESS1> -to <ADDRESS2> -amount 10

# Check balances
./build/blockchain getbalance -address <ADDRESS1>
./build/blockchain getbalance -address <ADDRESS2>

# View blockchain
./build/blockchain printchain
```

## 📊 What to Expect

1. **Node Startup**
   - Nodes connect to seed node (localhost:3000)
   - Exchange version and blockchain height
   - Synchronize blockchain

2. **Transaction Flow**
   - Transaction created on any node
   - Broadcast to all connected peers
   - Added to mempool

3. **Mining**
   - Miners collect transactions from mempool
   - Mine new block with PoW
   - Broadcast new block to network
   - All nodes validate and add block

4. **Consensus**
   - All nodes maintain same blockchain
   - Longest chain rule (like Bitcoin)
   - Automatic fork resolution

## 🔍 Debugging Tips

### Check Node Connectivity

```bash
# Test if seed is listening
nc -zv localhost 3000

# See connected peers
./build/blockchain peers
```

### View Docker Logs

```bash
# All nodes
docker-compose logs

# Specific time range (last 10 minutes)
docker-compose logs --since 10m

# Follow live
docker-compose logs -f node-miner1
```

### Access Docker Container

```bash
# Interactive shell
docker exec -it blockchain-seed sh

# Run commands
docker exec -it blockchain-seed /app/blockchain printchain
```

## 🎯 Test Scenarios

### Scenario 1: Basic Network

1. Start seed node + 2 miners
2. Create transaction
3. Watch blocks being mined
4. Verify all nodes have same blockchain

### Scenario 2: Late Joiner

1. Start seed + miner 1
2. Mine several blocks
3. Start miner 2 (late joiner)
4. Verify miner 2 syncs blockchain

### Scenario 3: Multiple Transactions

1. Create 3 wallets
2. Send multiple transactions
3. Watch miners compete
4. Verify all UTXOs are correct

## 📝 Common Issues

### Port Already in Use

```bash
# Find process
lsof -i :3000

# Kill process
kill -9 <PID>

# Or use different port
./build/blockchain startnode -port 3005
```

### Blockchain Not Syncing

```bash
# Reindex UTXO
./build/blockchain reindexutxo

# Or clean and restart
rm -rf ./tmp
./build/blockchain createblockchain -address <ADDRESS>
```

### Docker Build Fails

```bash
# Clean Docker cache
docker system prune -a

# Rebuild
make docker-build
```

## 🎓 Learning Exercises

1. **Modify Mining Difficulty**
   - Edit `internal/blockchain/proof.go`
   - Change `Difficulty` constant
   - Observe mining time changes

2. **Change Mempool Threshold**
   - Edit `internal/network/server.go`
   - Modify `len(memoryPool) >= 2`
   - Test different transaction batching

3. **Add Network Statistics**
   - Track messages sent/received
   - Monitor peer connections
   - Log synchronization time

4. **Implement Persistence**
   - Save peer list to disk
   - Restore connections on restart
   - Add peer reputation scoring

## 📚 Next Steps

- Read [docs/NETWORK.md](docs/NETWORK.md) for detailed documentation
- Explore Bitcoin P2P protocol differences
- Implement additional features:
  - Compact block relay
  - Transaction mempool priority
  - Peer discovery protocol
  - Network statistics dashboard

## 🐛 Troubleshooting

If something doesn't work:

1. Check if blockchain exists: `ls -la tmp/`
2. Verify ports are available: `netstat -an | grep 3000`
3. Look at logs: `make docker-logs`
4. Clean everything: `make docker-clean && rm -rf tmp/`
5. Start fresh: Follow "Quick Test with Docker" from step 1

---

**Happy Blockchain Networking! 🎉**

