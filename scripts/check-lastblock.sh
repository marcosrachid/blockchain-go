#!/bin/bash

# Script to check the last block of all nodes
# Shows the current state of each node in the network

echo "======================================"
echo "   LAST BLOCK STATUS - ALL NODES"
echo "======================================"
echo

# Array of nodes with their names and ports
declare -A NODES=(
    ["Seed Node"]="4000"
    ["Miner 1"]="4001"
    ["Miner 2"]="4002"
    ["Regular Node"]="4003"
)

# Function to get last block info
get_last_block() {
    local node_name=$1
    local port=$2
    
    echo "📦 $node_name (localhost:$port)"
    echo "   $(printf '─%.0s' {1..50})"
    
    # Try to get last block info
    response=$(curl -s http://localhost:$port/api/lastblock 2>/dev/null)
    
    if [ $? -eq 0 ] && [ ! -z "$response" ]; then
        # Parse JSON response
        height=$(echo $response | python3 -c "import sys, json; print(json.load(sys.stdin).get('height', 'N/A'))" 2>/dev/null)
        hash=$(echo $response | python3 -c "import sys, json; print(json.load(sys.stdin).get('hash', 'N/A'))" 2>/dev/null)
        timestamp=$(echo $response | python3 -c "import sys, json; print(json.load(sys.stdin).get('timestamp', 'N/A'))" 2>/dev/null)
        difficulty=$(echo $response | python3 -c "import sys, json; print(json.load(sys.stdin).get('difficulty', 'N/A'))" 2>/dev/null)
        
        echo "   Height:     $height"
        echo "   Hash:       ${hash:0:64}"
        echo "   Timestamp:  $timestamp"
        echo "   Difficulty: $difficulty"
    else
        echo "   ❌ Node is not responding"
    fi
    
    echo
}

# Check all nodes
for node_name in "Seed Node" "Miner 1" "Miner 2" "Regular Node"; do
    get_last_block "$node_name" "${NODES[$node_name]}"
done

echo "======================================"
echo "   Check completed at $(date '+%Y-%m-%d %H:%M:%S')"
echo "======================================"

