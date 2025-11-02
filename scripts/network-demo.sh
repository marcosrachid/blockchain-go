#!/bin/bash

# Simple Network Demo Script
# Demonstrates blockchain network with multiple nodes

set -e

echo "=========================================="
echo "   Blockchain Network Demo"
echo "=========================================="
echo ""

# Clean old data
echo "🧹 Cleaning old data..."
rm -rf ./tmp_node*
rm -rf ./build
echo ""

# Build project
echo "🔨 Building project..."
make build
echo ""

echo "=========================================="
echo "   Starting Multi-Node Network"
echo "=========================================="
echo ""

# Create data directories
mkdir -p ./tmp_node1
mkdir -p ./tmp_node2
mkdir -p ./tmp_node3

echo "📝 This demo will start 3 nodes:"
echo "   - Node 1 (Port 3000): Seed node + Miner"
echo "   - Node 2 (Port 3001): Miner"
echo "   - Node 3 (Port 3002): Regular node"
echo ""

# Create wallets for each node
echo "👛 Creating wallets..."
WALLET1=$(./build/blockchain createwallet 2>&1 | grep "New address is:" | cut -d' ' -f4)
echo "   Node 1 Wallet: $WALLET1"

WALLET2=$(./build/blockchain createwallet 2>&1 | grep "New address is:" | cut -d' ' -f4)
echo "   Node 2 Wallet: $WALLET2"

WALLET3=$(./build/blockchain createwallet 2>&1 | grep "New address is:" | cut -d' ' -f4)
echo "   Node 3 Wallet: $WALLET3"
echo ""

# Create blockchain
echo "⛓️  Creating blockchain..."
./build/blockchain createblockchain -address $WALLET1
echo ""

echo "=========================================="
echo "   Network Nodes Ready!"
echo "=========================================="
echo ""
echo "To start the network, run these commands in separate terminals:"
echo ""
echo "Terminal 1 (Seed + Miner):"
echo "  ./build/blockchain startnode -port 3000 -miner $WALLET1"
echo ""
echo "Terminal 2 (Miner):"
echo "  ./build/blockchain startnode -port 3001 -miner $WALLET2"
echo ""
echo "Terminal 3 (Regular node):"
echo "  ./build/blockchain startnode -port 3002"
echo ""
echo "To send a transaction:"
echo "  ./build/blockchain send -from $WALLET1 -to $WALLET2 -amount 10"
echo ""
echo "To check balances:"
echo "  ./build/blockchain getbalance -address $WALLET1"
echo "  ./build/blockchain getbalance -address $WALLET2"
echo ""
echo "To view the blockchain:"
echo "  ./build/blockchain printchain"
echo ""
echo "Wallets created:"
echo "  Node 1: $WALLET1"
echo "  Node 2: $WALLET2"
echo "  Node 3: $WALLET3"
echo ""

