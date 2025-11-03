#!/bin/bash

# Script to check balances of all wallet addresses in all nodes
# Shows the current balance of each address in the network

echo "======================================"
echo "   WALLET BALANCES - ALL NODES"
echo "======================================"
echo

# Array of nodes with their names and ports
declare -A NODES=(
    ["Seed Node"]="4000"
    ["Miner 1"]="4001"
    ["Miner 2"]="4002"
    ["Regular Node"]="4003"
)

# Function to get balances for a node
get_node_balances() {
    local node_name=$1
    local port=$2
    
    echo "💰 $node_name (localhost:$port)"
    echo "   $(printf '─%.0s' {1..50})"
    
    # Try to get addresses
    addresses_response=$(curl -s http://localhost:$port/api/addresses 2>/dev/null)
    
    if [ $? -eq 0 ] && [ ! -z "$addresses_response" ]; then
        # Parse addresses array
        addresses=$(echo $addresses_response | python3 -c "import sys, json; data=json.load(sys.stdin); print('\n'.join(data.get('addresses', [])))" 2>/dev/null)
        
        if [ -z "$addresses" ]; then
            echo "   📭 No wallet addresses found"
        else
            total_balance=0
            
            # Check balance for each address
            while IFS= read -r address; do
                if [ ! -z "$address" ]; then
                    balance_response=$(curl -s "http://localhost:$port/api/balance/$address" 2>/dev/null)
                    balance=$(echo $balance_response | python3 -c "import sys, json; print(json.load(sys.stdin).get('balance', 0))" 2>/dev/null)
                    
                    if [ ! -z "$balance" ]; then
                        echo "   Address: $address"
                        echo "   Balance: $balance coins"
                        echo
                        total_balance=$((total_balance + balance))
                    fi
                fi
            done <<< "$addresses"
            
            echo "   $(printf '─%.0s' {1..50})"
            echo "   💎 Total Balance: $total_balance coins"
        fi
    else
        echo "   ❌ Node is not responding"
    fi
    
    echo
}

# Check all nodes
for node_name in "Seed Node" "Miner 1" "Miner 2" "Regular Node"; do
    get_node_balances "$node_name" "${NODES[$node_name]}"
done

echo "======================================"
echo "   Check completed at $(date '+%Y-%m-%d %H:%M:%S')"
echo "======================================"

