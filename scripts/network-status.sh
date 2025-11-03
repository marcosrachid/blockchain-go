#!/bin/bash

# Comprehensive network status script
# Shows a complete overview of the blockchain network

clear

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║         BLOCKCHAIN NETWORK - STATUS DASHBOARD                 ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo

# Get network info from seed node
echo "📊 NETWORK OVERVIEW"
echo "$(printf '═%.0s' {1..64})"
network_info=$(curl -s http://localhost:4000/api/networkinfo 2>/dev/null)

if [ $? -eq 0 ] && [ ! -z "$network_info" ]; then
    height=$(echo $network_info | python3 -c "import sys, json; print(json.load(sys.stdin).get('height', 'N/A'))" 2>/dev/null)
    difficulty=$(echo $network_info | python3 -c "import sys, json; print(json.load(sys.stdin).get('difficulty', 'N/A'))" 2>/dev/null)
    total_supply=$(echo $network_info | python3 -c "import sys, json; print(json.load(sys.stdin).get('total_supply', 'N/A'))" 2>/dev/null)
    max_supply=$(echo $network_info | python3 -c "import sys, json; print(json.load(sys.stdin).get('max_supply', 'N/A'))" 2>/dev/null)
    block_reward=$(echo $network_info | python3 -c "import sys, json; print(json.load(sys.stdin).get('current_block_reward', 'N/A'))" 2>/dev/null)
    blocks_until_halving=$(echo $network_info | python3 -c "import sys, json; print(json.load(sys.stdin).get('blocks_until_halving', 'N/A'))" 2>/dev/null)
    
    echo "  Current Height:        $height"
    echo "  Mining Difficulty:     $difficulty"
    echo "  Total Supply:          $total_supply coins"
    echo "  Max Supply:            $max_supply coins"
    echo "  Current Block Reward:  $block_reward coins"
    echo "  Blocks Until Halving:  $blocks_until_halving"
else
    echo "  ❌ Network info not available"
fi

echo
echo "🔗 NODE STATUS"
echo "$(printf '═%.0s' {1..64})"

# Array of nodes
declare -A NODES=(
    ["Seed"]="4000"
    ["Miner1"]="4001"
    ["Miner2"]="4002"
    ["Regular"]="4003"
)

# Check each node
for node_name in "Seed" "Miner1" "Miner2" "Regular"; do
    port="${NODES[$node_name]}"
    response=$(curl -s http://localhost:$port/api/lastblock 2>/dev/null)
    
    if [ $? -eq 0 ] && [ ! -z "$response" ]; then
        node_height=$(echo $response | python3 -c "import sys, json; print(json.load(sys.stdin).get('height', 'N/A'))" 2>/dev/null)
        echo "  ✅ $node_name (port $port) - Height: $node_height"
    else
        echo "  ❌ $node_name (port $port) - OFFLINE"
    fi
done

echo
echo "💰 WALLET BALANCES"
echo "$(printf '═%.0s' {1..64})"

total_network_balance=0

for node_name in "Seed" "Miner1" "Miner2" "Regular"; do
    port="${NODES[$node_name]}"
    addresses_response=$(curl -s http://localhost:$port/api/addresses 2>/dev/null)
    
    if [ $? -eq 0 ] && [ ! -z "$addresses_response" ]; then
        addresses=$(echo $addresses_response | python3 -c "import sys, json; data=json.load(sys.stdin); print('\n'.join(data.get('addresses', [])))" 2>/dev/null)
        
        if [ ! -z "$addresses" ]; then
            node_balance=0
            
            while IFS= read -r address; do
                if [ ! -z "$address" ]; then
                    balance_response=$(curl -s "http://localhost:$port/api/balance/$address" 2>/dev/null)
                    balance=$(echo $balance_response | python3 -c "import sys, json; print(json.load(sys.stdin).get('balance', 0))" 2>/dev/null)
                    
                    if [ ! -z "$balance" ]; then
                        node_balance=$((node_balance + balance))
                    fi
                fi
            done <<< "$addresses"
            
            total_network_balance=$((total_network_balance + node_balance))
            echo "  💎 $node_name: $node_balance coins"
        else
            echo "  📭 $node_name: 0 coins (no wallets)"
        fi
    else
        echo "  ❌ $node_name: unavailable"
    fi
done

echo "  $(printf '─%.0s' {1..62})"
echo "  💰 TOTAL NETWORK: $total_network_balance coins"

echo
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║  Last updated: $(date '+%Y-%m-%d %H:%M:%S')                              ║"
echo "╚════════════════════════════════════════════════════════════════╝"

