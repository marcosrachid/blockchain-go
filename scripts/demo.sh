#!/bin/bash

# Blockchain Demo Script

set -e

echo "=========================================="
echo "   Blockchain in Go - Demo"
echo "=========================================="
echo ""

# Clean old data
echo "🧹 Cleaning old data..."
rm -rf ./tmp
rm -rf ./build
echo ""

# Build project
echo "🔨 Building project..."
make build
echo ""

# Create tmp directory for wallet storage
echo "📁 Creating data directory..."
mkdir -p ./tmp
echo ""

echo "=========================================="
echo "   Step 1: Create Wallets"
echo "=========================================="

# Create wallets
echo "👛 Creating wallet for Alice..."
ALICE=$(./build/blockchain createwallet 2>&1 | grep "New address is:" | cut -d' ' -f4)
echo "   Alice: $ALICE"
echo ""

echo "👛 Creating wallet for Bob..."
BOB=$(./build/blockchain createwallet 2>&1 | grep "New address is:" | cut -d' ' -f4)
echo "   Bob: $BOB"
echo ""

echo "👛 Creating wallet for Charlie..."
CHARLIE=$(./build/blockchain createwallet 2>&1 | grep "New address is:" | cut -d' ' -f4)
echo "   Charlie: $CHARLIE"
echo ""

echo "📋 Listing all addresses:"
./build/blockchain listaddresses
echo ""

echo "=========================================="
echo "   Step 2: Create Blockchain"
echo "=========================================="
echo "⛏️  Creating blockchain with reward to Alice..."
./build/blockchain createblockchain -address $ALICE 2>&1 | grep -v "badger" | grep -v "DEBUG" | grep -v "INFO"
echo ""

echo "=========================================="
echo "   Step 3: Initial Balances"
echo "=========================================="
echo "💰 Checking initial balances:"
./build/blockchain getbalance -address $ALICE 2>&1 | grep "Balance"
./build/blockchain getbalance -address $BOB 2>&1 | grep "Balance"
./build/blockchain getbalance -address $CHARLIE 2>&1 | grep "Balance"
echo ""

echo "=========================================="
echo "   Step 4: First Transaction"
echo "=========================================="
echo "💸 Alice sends 10 coins to Bob..."
./build/blockchain send -from $ALICE -to $BOB -amount 10 2>&1 | grep -v "badger" | grep -v "DEBUG" | grep -v "INFO"
echo ""

echo "💰 Balances after first transaction:"
./build/blockchain getbalance -address $ALICE 2>&1 | grep "Balance"
./build/blockchain getbalance -address $BOB 2>&1 | grep "Balance"
./build/blockchain getbalance -address $CHARLIE 2>&1 | grep "Balance"
echo ""

echo "=========================================="
echo "   Step 5: Second Transaction"
echo "=========================================="
echo "💸 Alice sends 20 coins to Charlie..."
./build/blockchain send -from $ALICE -to $CHARLIE -amount 20 2>&1 | grep -v "badger" | grep -v "DEBUG" | grep -v "INFO"
echo ""

echo "💰 Balances after second transaction:"
./build/blockchain getbalance -address $ALICE 2>&1 | grep "Balance"
./build/blockchain getbalance -address $BOB 2>&1 | grep "Balance"
./build/blockchain getbalance -address $CHARLIE 2>&1 | grep "Balance"
echo ""

echo "=========================================="
echo "   Step 6: Bob sends to Charlie"
echo "=========================================="
echo "💸 Bob sends 5 coins to Charlie..."
./build/blockchain send -from $BOB -to $CHARLIE -amount 5 2>&1 | grep -v "badger" | grep -v "DEBUG" | grep -v "INFO"
echo ""

echo "💰 Final balances:"
echo "   Alice:"
./build/blockchain getbalance -address $ALICE 2>&1 | grep "Balance"
echo "   Bob:"
./build/blockchain getbalance -address $BOB 2>&1 | grep "Balance"
echo "   Charlie:"
./build/blockchain getbalance -address $CHARLIE 2>&1 | grep "Balance"
echo ""

echo "=========================================="
echo "   Step 7: View Blockchain"
echo "=========================================="
echo "📊 Printing blockchain..."
./build/blockchain printchain 2>&1 | grep -v "badger" | grep -v "DEBUG" | grep -v "INFO" | head -50
echo ""
echo "(Output truncated for readability)"
echo ""

echo "=========================================="
echo "   Step 8: UTXO Set"
echo "=========================================="
echo "🔍 Reindexing and counting UTXOs..."
./build/blockchain reindexutxo 2>&1 | grep -v "badger" | grep -v "DEBUG" | grep -v "INFO"
echo ""

echo "=========================================="
echo "   ✅ Demo Complete!"
echo "=========================================="
echo ""
echo "Operations summary:"
echo "1. Created 3 wallets (Alice, Bob, Charlie)"
echo "2. Initialized blockchain with 50 coin reward to Alice"
echo "3. Alice sent 10 to Bob (mined and earned +50)"
echo "4. Alice sent 20 to Charlie (mined and earned +50)"
echo "5. Bob sent 5 to Charlie (mined and earned +50)"
echo ""
echo "Expected final balances:"
echo "- Alice: 120 coins (50 initial - 10 - 20 + 50 + 50)"
echo "- Bob: 55 coins (10 received - 5 + 50)"
echo "- Charlie: 25 coins (20 + 5 received)"
echo ""
echo "Total blocks: 4 (1 genesis + 3 mined)"
echo ""
echo "To explore more, use these commands:"
echo "  ./build/blockchain listaddresses"
echo "  ./build/blockchain getbalance -address ADDRESS"
echo "  ./build/blockchain printchain"
echo ""
