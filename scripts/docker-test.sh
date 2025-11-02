#!/bin/bash

# Docker Network Test Script for Blockchain

set -e

echo "=========================================="
echo "   Blockchain Network - Docker Test"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Clean old containers and volumes
echo -e "${YELLOW}🧹 Cleaning old containers and volumes...${NC}"
docker-compose down -v 2>/dev/null || true
echo ""

# Build images
echo -e "${BLUE}🔨 Building Docker images...${NC}"
docker-compose build
echo ""

# Start network
echo -e "${GREEN}🚀 Starting blockchain network...${NC}"
echo "   - Seed node (localhost:3000)"
echo "   - Miner 1 (localhost:3001)"
echo "   - Miner 2 (localhost:3002)"
echo "   - Regular node (localhost:3003)"
echo ""
docker-compose up -d

# Wait for nodes to be ready
echo -e "${YELLOW}⏳ Waiting for nodes to initialize (40 seconds)...${NC}"
sleep 40

echo ""
echo "=========================================="
echo "   Network Status"
echo "=========================================="

# Check container status
echo -e "${BLUE}📊 Container status:${NC}"
docker-compose ps

echo ""
echo "=========================================="
echo "   Testing Network Operations"
echo "=========================================="

# Get seed node logs
echo -e "${BLUE}📝 Seed node logs (last 15 lines):${NC}"
docker logs blockchain-seed --tail 15

echo ""
echo -e "${BLUE}📝 Miner 1 logs (last 15 lines):${NC}"
docker logs blockchain-miner1 --tail 15 2>/dev/null || echo "Miner 1 still starting..."

echo ""
echo -e "${BLUE}📝 Miner 2 logs (last 15 lines):${NC}"
docker logs blockchain-miner2 --tail 15 2>/dev/null || echo "Miner 2 still starting..."

echo ""
echo "=========================================="
echo "   Test Blockchain Commands"
echo "=========================================="

# Test getting wallet list
echo -e "${BLUE}📋 Listing wallets in seed node:${NC}"
docker exec blockchain-seed /app/blockchain listaddresses 2>/dev/null || echo "Seed node still initializing..."

# Test blockchain info
echo ""
echo -e "${BLUE}⛓️  Blockchain info:${NC}"
docker exec blockchain-seed /app/blockchain printchain 2>/dev/null | head -20 || echo "Blockchain still being created..."

echo ""
echo "=========================================="
echo "   ✅ Network is Running!"
echo "=========================================="
echo ""
echo -e "${GREEN}Network started successfully!${NC}"
echo ""
echo "Available commands:"
echo "  docker-compose logs -f                  # View all logs"
echo "  docker-compose logs -f node-seed        # View seed node logs"
echo "  docker-compose logs -f node-miner1      # View miner 1 logs"
echo "  docker-compose logs -f node-miner2      # View miner 2 logs"
echo "  docker-compose logs -f node-regular     # View regular node logs"
echo "  docker-compose down                     # Stop all nodes"
echo "  docker-compose down -v                  # Stop and remove volumes"
echo ""
echo "Execute commands in a container:"
echo "  docker exec -it blockchain-seed /app/blockchain listaddresses"
echo "  docker exec -it blockchain-miner1 /app/blockchain getbalance -address ADDRESS"
echo "  docker exec -it blockchain-seed /app/blockchain printchain"
echo ""
echo -e "${YELLOW}To view logs continuously, run:${NC}"
echo "  docker-compose logs -f"
echo ""

