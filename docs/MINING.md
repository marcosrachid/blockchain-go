# Mining and Difficulty Adjustment

## Current Implementation

### Proof of Work (PoW)
- **Fixed Difficulty**: 18 (defined in `config.go`)
- Miners compute SHA-256 hashes continuously until finding a hash with 18 leading zeros
- No timers - mining happens **continuously** until a valid block is found

### How Mining Works

1. **Continuous Mining**: When a node is set as a miner, it starts an infinite loop
2. **Block Creation**: Collects transactions from mempool + coinbase transaction
3. **PoW Computation**: Tries different nonce values until hash meets difficulty target
4. **Broadcast**: Once found, immediately broadcasts to all peers
5. **Repeat**: Immediately starts mining the next block

```go
// Mining loop (simplified)
for {
    // Collect transactions
    txs := getTransactionsFromMempool()
    coinbase := createCoinbaseTransaction()
    
    // Mine block (PoW)
    block := mineBlock(txs + coinbase) // This takes time based on difficulty
    
    // Broadcast
    broadcastBlock(block)
}
```

## Bitcoin's Difficulty Adjustment

In Bitcoin, difficulty adjusts **automatically** every 2016 blocks (~2 weeks):

```
new_difficulty = current_difficulty * (actual_time / target_time)
```

- **Target**: 10 minutes per block
- **Adjustment Interval**: 2016 blocks
- If blocks were faster → increase difficulty
- If blocks were slower → decrease difficulty

### Example

If 2016 blocks took **1 week** instead of **2 weeks**:
- Blocks were mined **2x faster** than target
- New difficulty = current × (1 week / 2 weeks) = current × 0.5
- **Difficulty increases** to slow down mining

## Future Improvements

### 1. Dynamic Difficulty Adjustment

Add to `internal/blockchain/config.go`:

```go
const (
    DifficultyAdjustmentInterval = 100 // Adjust every 100 blocks
    TargetBlockTime = 60 // seconds
)
```

Implement in `blockchain.go`:

```go
func (chain *Blockchain) CalculateNewDifficulty() int {
    lastBlock := chain.GetLastBlock()
    
    // Get block from adjustment interval ago
    targetBlock := chain.GetBlockAtHeight(lastBlock.Height - DifficultyAdjustmentInterval)
    
    // Calculate actual time
    actualTime := lastBlock.Timestamp - targetBlock.Timestamp
    expectedTime := DifficultyAdjustmentInterval * TargetBlockTime
    
    // Adjust difficulty
    if actualTime < expectedTime / 2 {
        // Blocks too fast, increase difficulty
        return currentDifficulty + 1
    } else if actualTime > expectedTime * 2 {
        // Blocks too slow, decrease difficulty
        return max(currentDifficulty - 1, 1)
    }
    
    return currentDifficulty
}
```

### 2. Mining Competition

Current behavior with multiple miners:
- All miners compete to find the next block
- First to find valid PoW broadcasts and wins
- Others discard their work and start on next block
- This is **correct** Bitcoin-like behavior

### 3. Orphan Blocks

When two miners find blocks simultaneously:
- Both broadcast their blocks
- Network may temporarily split
- Next block determines winner (longest chain rule)
- Losing block becomes "orphan" and is discarded

## Why No Timer?

❌ **Wrong Approach** (previous implementation):
```go
ticker := time.NewTicker(60 * time.Second)
for {
    <-ticker.C
    mineBlock() // Forces mining every 60 seconds
}
```

✅ **Correct Approach** (current):
```go
for {
    block := mineBlock() // Takes time naturally based on difficulty
    broadcast(block)
    // Immediately start next block
}
```

The **difficulty** controls the time, not a timer!

## Configuration

To adjust target block time, modify in `config.go`:

```go
const (
    Difficulty = 18  // Higher = slower blocks, more security
                     // Lower = faster blocks, less security
)
```

- Difficulty 12: ~0.1 seconds per block
- Difficulty 18: ~1 minute per block (current)
- Difficulty 20: ~4 minutes per block
- Difficulty 24: ~1 hour per block

## References

- [Bitcoin Difficulty Adjustment](https://en.bitcoin.it/wiki/Difficulty)
- [Proof of Work](https://en.bitcoin.it/wiki/Proof_of_work)
- [Mining](https://en.bitcoin.it/wiki/Mining)

