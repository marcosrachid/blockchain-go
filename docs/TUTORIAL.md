# Tutorial: Getting Started with the Blockchain

> **Versão em Português**: [TUTORIAL.pt-br.md](TUTORIAL.pt-br.md)

This tutorial will show you how to use the blockchain step by step.

## Step 1: Build the Project

```bash
# Navigate to project directory
cd blockchain-go

# Install dependencies
make deps

# Build the project
make build
```

This creates the `./build/blockchain` executable.

## Step 2: Create Your First Wallet

```bash
./build/blockchain createwallet
```

**Output:**
```
New address is: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
```

💡 **Save this address!** You'll need it for the next steps.

### What Happened?

- Generated an ECDSA key pair
- Created a Bitcoin-like address with Base58 encoding
- Saved the wallet to `./tmp/wallets.dat`

## Step 3: Create More Wallets

Create at least 2 more wallets:

```bash
./build/blockchain createwallet
./build/blockchain createwallet
```

### List All Addresses

```bash
./build/blockchain listaddresses
```

**Output:**
```
1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2
1JfbZRwdDHKZmuiZgYArJZhcuuzuw2HuMu
```

## Step 4: Create the Blockchain

Use the first address to receive the genesis block reward:

```bash
./build/blockchain createblockchain -address 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
```

**Output:**
```
Mining genesis block...
[Lots of hash attempts...]
Genesis created
Finished!
```

### What Happened?

- Created genesis block with Proof of Work
- Coinbase transaction gave 50 coins to your address
- Saved blockchain to `./tmp/blocks/`
- Built UTXO set

## Step 5: Check Your Balance

```bash
./build/blockchain getbalance -address 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
```

**Output:**
```
Balance of 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa: 50
```

🎉 You have 50 coins from the genesis block!

## Step 6: Send Your First Transaction

Send coins to another address:

```bash
./build/blockchain send \
  -from 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa \
  -to 1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2 \
  -amount 10
```

**Output:**
```
Mining new block...
[Proof of Work...]
Success!
```

### What Happened?

1. Created transaction spending 10 coins
2. Found UTXOs to cover the amount
3. Created change output for remaining coins
4. Signed transaction with ECDSA
5. Mined new block with transaction + coinbase
6. Updated UTXO set

## Step 7: Verify Balances

Check sender balance:
```bash
./build/blockchain getbalance -address 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
```

**Output:** `90` (50 initial - 10 sent + 50 mining reward)

Check receiver balance:
```bash
./build/blockchain getbalance -address 1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2
```

**Output:** `10` (received from transaction)

## Step 8: View the Blockchain

```bash
./build/blockchain printchain
```

**Output:**
```
============ Block 000abc... ============
Height: 1
Prev. hash: 000def...
PoW: true
--- Transaction 123abc... ---
   Input 0:
      TXID: 456def
      Out: 0
      Signature: 789ghi...
   Output 0:
      Value: 10
      To: 1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2
   Output 1:
      Value: 40
      To: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa (change)

============ Block 000def... ============
Height: 0
Prev. hash:
PoW: true
--- Transaction 456def... ---
   Output 0:
      Value: 50
      To: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
```

## Step 9: Multiple Transactions

Send more transactions:

```bash
# Send from first to third address
./build/blockchain send \
  -from 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa \
  -to 1JfbZRwdDHKZmuiZgYArJZhcuuzuw2HuMu \
  -amount 25

# Send from second to third
./build/blockchain send \
  -from 1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2 \
  -to 1JfbZRwdDHKZmuiZgYArJZhcuuzuw2HuMu \
  -amount 5
```

Check all balances:
```bash
./build/blockchain getbalance -address 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
./build/blockchain getbalance -address 1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2
./build/blockchain getbalance -address 1JfbZRwdDHKZmuiZgYArJZhcuuzuw2HuMu
```

## Step 10: Reindex UTXO Set

If needed, rebuild the UTXO set from the blockchain:

```bash
./build/blockchain reindexutxo
```

**Output:**
```
Done! There are 6 transactions in the UTXO set.
```

## 🌐 Network Tutorial

### Start Multiple Nodes

**Terminal 1 - Seed + Miner:**
```bash
./build/blockchain startnode -port 3000 -miner 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
```

**Terminal 2 - Miner:**
```bash
./build/blockchain startnode -port 3001 -miner 1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2
```

**Terminal 3 - Send Transaction:**
```bash
./build/blockchain send \
  -from 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa \
  -to 1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2 \
  -amount 10
```

Watch the transaction propagate and get mined!

## 📊 Understanding the Output

### Transaction Structure

```
Transaction abc123...
  Inputs: [Reference to previous outputs]
  Outputs: [New outputs with values]
  Signature: [ECDSA signature]
```

### Block Structure

```
Block hash: 000abc...
  Timestamp: 1234567890
  Height: 5
  Previous Hash: 000def...
  Nonce: 123456
  Transactions: [...]
```

### UTXO Set

Collection of all unspent outputs:
```
{
  "tx_id:output_index": {value, address},
  ...
}
```

## 🎯 Exercises

1. **Create 5 wallets** and distribute coins among them
2. **Send circular transactions** (A→B, B→C, C→A)
3. **Monitor mining** by printing the chain after each transaction
4. **Test edge cases**:
   - Try sending more than you have
   - Send to an invalid address
   - Use a non-existent wallet

## 🐛 Troubleshooting

### "No existing blockchain found"

**Solution:** Create a blockchain first:
```bash
./build/blockchain createblockchain -address YOUR_ADDRESS
```

### "Address is not valid"

**Solution:** Check your address with:
```bash
./build/blockchain listaddresses
```

### "Not enough funds"

**Solution:** Check your balance:
```bash
./build/blockchain getbalance -address YOUR_ADDRESS
```

You need enough to cover the amount + it will be used for mining.

## 📚 Next Steps

- Read [Bitcoin Comparison](BITCOIN_COMPARISON.md) to understand similarities
- Try the [Network Quick Start](../QUICKSTART_NETWORK.md) guide
- Explore [Improvements](IMPROVEMENTS.md) for future enhancements
- Study the [Architecture](../ARCHITECTURE.md) document

## 💡 Key Concepts

### Proof of Work
- Mining finds a hash with leading zeros
- Difficulty determines number of zeros
- Nonce is incremented until valid hash found

### UTXO Model
- Outputs can only be spent once
- Transactions reference previous outputs
- Change is returned to sender

### Digital Signatures
- Private key signs transactions
- Public key verifies signatures
- Prevents transaction tampering

### Merkle Tree
- Efficient transaction hashing
- Root included in block header
- Allows SPV verification

---

**Congratulations! You've completed the tutorial!** 🎉

Now you understand the basics of blockchain. Try experimenting with the network features next!

