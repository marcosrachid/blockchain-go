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
echo -e "${YELLOW}⏳ Waiting for nodes to initialize (30 seconds)...${NC}"
sleep 30

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
echo -e "${BLUE}📝 Seed node logs (last 10 lines):${NC}"
docker logs blockchain-seed --tail 10

echo ""
echo -e "${BLUE}📝 Miner 1 logs (last 10 lines):${NC}"
docker logs blockchain-miner1 --tail 10

echo ""
echo "=========================================="
echo "   Network is Running!"
echo "=========================================="
echo ""
echo "Available commands:"
echo "  docker-compose logs -f node-seed       # View seed node logs"
echo "  docker-compose logs -f node-miner1     # View miner 1 logs"
echo "  docker-compose logs -f node-miner2     # View miner 2 logs"
echo "  docker-compose logs -f node-regular    # View regular node logs"
echo "  docker-compose down                     # Stop all nodes"
echo "  docker-compose down -v                  # Stop and remove volumes"
echo ""
echo "Execute commands in a container:"
echo "  docker exec -it blockchain-seed /app/blockchain listaddresses"
echo "  docker exec -it blockchain-miner1 /app/blockchain getbalance -address ADDRESS"
echo "  docker exec -it blockchain-seed /app/blockchain printchain"
echo ""
echo "Press Ctrl+C to stop watching logs"
echo ""

# Follow logs
docker-compose logs -f

