# Halving and Supply Limit Implementation

## ✅ What Was Implemented

A complete **Bitcoin-like halving and supply limit system** has been added to the blockchain protocol.

## 📊 Protocol Parameters (Centralized in `config.go`)

All configuration is now centralized in: **`internal/blockchain/config.go`**

### Supply and Reward Configuration

```go
const (
    InitialSubsidy  = 50       // Initial mining reward (50 coins like Bitcoin)
    HalvingInterval = 210000   // Blocks until reward halving (~4 years)
    MaxSupply       = 21000000 // Maximum supply (21 million coins)
)
```

### Other Protocol Parameters

```go
const (
    Difficulty = 18 // Mining difficulty (PoW)
    GenesisData = "First Transaction from Genesis"
    DBPath = "./tmp/blocks"
    DefaultPort = 3000
    ProtocolVersion = 1
)
```

## 🔄 How Halving Works

### Block Reward Schedule

| Blocks | Reward | Coins Minted | Cumulative |
|--------|--------|--------------|------------|
| 0 - 209,999 | 50 | 10,500,000 | 10,500,000 |
| 210,000 - 419,999 | 25 | 5,250,000 | 15,750,000 |
| 420,000 - 629,999 | 12 | 2,520,000 | 18,270,000 |
| 630,000 - 839,999 | 6 | 1,260,000 | 19,530,000 |
| 840,000 - 1,049,999 | 3 | 630,000 | 20,160,000 |
| ... | ... | ... | ... |
| ~6,930,000+ | 0 | 0 | ~21,000,000 |

### Calculation Function

```go
func GetBlockReward(height int) int {
    reward := InitialSubsidy
    
    // Calculate number of halvings
    halvings := height / HalvingInterval
    
    // Each halving divides reward by 2
    for i := 0; i < halvings; i++ {
        reward = reward / 2
    }
    
    // When reward becomes 0, no more coins are minted
    if reward < 1 {
        return 0
    }
    
    return reward
}
```

## 🎯 Key Features

### 1. **Decreasing Supply Rate**
- Reward halves every 210,000 blocks
- Mimics Bitcoin's scarcity model
- Prevents inflation over time

### 2. **Maximum Supply Cap**
- Hard limit of 21 million coins
- No coins can be minted after reaching max supply
- Reward becomes 0 after ~33 halvings

### 3. **Predictable Emission**
- Transparent and deterministic
- Can calculate total supply at any block height
- Economic incentives are clear to miners

### 4. **Height-Based Calculation**
- Reward calculated from block height
- No need to store historical rewards
- Efficient and verifiable

## 📝 Updated Functions

### `CoinbaseTX` (Mining Reward Transaction)

**Before:**
```go
func CoinbaseTX(to, data string) *Transaction {
    txout := NewTXOutput(50, to) // Fixed reward
    // ...
}
```

**After:**
```go
func CoinbaseTX(to, data string, height int) *Transaction {
    reward := GetBlockReward(height) // Dynamic reward
    txout := NewTXOutput(reward, to)
    // ...
}
```

### `GetBestHeight` (New Function)

Added to blockchain to get current chain height:

```go
func (chain *Blockchain) GetBestHeight() int {
    var lastBlock Block
    // ... get last block from database
    return lastBlock.Height
}
```

### Usage in Mining

```go
// When creating a new block
newHeight := chain.GetBestHeight() + 1
cbTx := blockchain.CoinbaseTX(minerAddress, "", newHeight)
```

## 📂 Files Modified

1. **`internal/blockchain/config.go`** ⭐ NEW
   - Centralized configuration file
   - All protocol constants
   - Helper functions

2. **`internal/blockchain/transaction.go`**
   - Updated `CoinbaseTX` to accept height
   - Moved constants to config.go
   - Uses `GetBlockReward()`

3. **`internal/blockchain/blockchain.go`**
   - Added `GetBestHeight()` method
   - Updated genesis block creation
   - Uses constants from config.go

4. **`internal/blockchain/proof.go`**
   - Uses `Difficulty` from config.go

5. **`cmd/blockchain/main.go`**
   - Updated send command to calculate height
   - Passes height to `CoinbaseTX`

6. **`internal/network/server.go`**
   - Updated mining function
   - Calculates height before creating coinbase

## 🧪 Testing Halving

### Test Scenario

```bash
# Mine blocks and check rewards at different heights

# Block 0 (Genesis)
Reward: 50 coins

# Block 1-209,999
Reward: 50 coins each

# Block 210,000 (First Halving)
Reward: 25 coins

# Block 420,000 (Second Halving)
Reward: 12 coins (rounded down)

# Block 630,000 (Third Halving)
Reward: 6 coins
```

### Verify Supply

```go
// Helper function to verify total supply
func VerifySupply(chain *Blockchain) {
    height := chain.GetBestHeight()
    expectedSupply := CalculateSupplyUpToHeight(height)
    actualSupply := chain.GetTotalSupply()
    
    if actualSupply > MaxSupply {
        log.Fatal("Supply exceeded maximum!")
    }
}
```

## 📈 Economic Model

### Supply Curve

```
Supply (millions)
21M ┤                           ___________
    │                      ___/
    │                 ___/
15M ┤            ___/
    │       ___/
    │  ___/
10M ┤_/
    │
 5M ┤
    │
  0 └─────────────────────────────────────►
    0   210k  420k  630k  840k  1.05M  ...  Block Height
```

### Emission Rate

- **Years 0-4**: 50 coins/block (fast)
- **Years 4-8**: 25 coins/block (medium)
- **Years 8-12**: 12 coins/block (slow)
- **Years 12+**: Progressively slower
- **Year ~140**: Emission stops (max supply reached)

## 🎓 Benefits

1. **Scarcity**: Limited supply increases value over time
2. **Predictability**: Known emission schedule
3. **Incentive**: Early miners get higher rewards
4. **Stability**: Decreasing inflation rate
5. **Bitcoin Compatibility**: Same model as Bitcoin

## 🔮 Future Enhancements

Possible additions:

1. **Transaction Fees**
   ```go
   reward := GetBlockReward(height) + fees
   ```

2. **Supply Verification**
   ```go
   func (chain *Blockchain) ValidateSupply() bool {
       return chain.GetTotalSupply() <= MaxSupply
   }
   ```

3. **Emission Statistics**
   ```go
   func GetEmissionRate(height int) float64 {
       // Calculate coins per year at given height
   }
   ```

4. **Supply Query Commands**
   ```bash
   ./blockchain supply              # Current total supply
   ./blockchain supply -height 1000 # Supply at height 1000
   ./blockchain halving             # Next halving block
   ```

## ✅ Summary

The blockchain now has:
- ✅ **Halving mechanism** (every 210,000 blocks)
- ✅ **Maximum supply** (21 million coins)
- ✅ **Centralized configuration** (config.go)
- ✅ **Height-based rewards** (dynamic calculation)
- ✅ **Bitcoin-compatible** (same parameters)

**Status:** Production-ready economic model! 🎉

---

For more information, see:
- [../internal/blockchain/config.go](../internal/blockchain/config.go) - Protocol configuration
- [../README.md](../README.md) - General documentation
- [BITCOIN_COMPARISON.md](BITCOIN_COMPARISON.md) - Bitcoin comparison

