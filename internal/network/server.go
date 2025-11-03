package network

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/marcocsrachid/blockchain-go/internal/api"
	"github.com/marcocsrachid/blockchain-go/internal/blockchain"
)

const (
	protocol      = "tcp"
	version       = 1
	commandLength = 12
	// Target block time for difficulty adjustment (not used as a timer!)
	targetBlockTime = 60 * time.Second // 1 minute target (Bitcoin = 10 min)
)

var (
	nodeAddress     string
	miningAddress   string
	knownNodes      = initKnownNodes()
	blocksInTransit = [][]byte{}
	memoryPool      = make(map[string]*blockchain.Transaction)
	mempoolMux      sync.RWMutex
)

// initKnownNodes initializes known nodes from environment or default
func initKnownNodes() []string {
	if seedNode := os.Getenv("SEED_NODE"); seedNode != "" {
		return []string{seedNode}
	}
	return []string{"localhost:3000"} // Default seed node
}

// Server represents the network server
type Server struct {
	Address         string
	Blockchain      *blockchain.Blockchain
	Peers           *PeerList
	IsMining        bool
	stopMining      chan bool
	miningInterrupt chan bool
	APIServer       *api.Server
	Wallets         *blockchain.Wallets
}

// NewServer creates a new network server
func NewServer(address string, bc *blockchain.Blockchain, wallets *blockchain.Wallets) *Server {
	// Extract port from address for API
	parts := strings.Split(address, ":")
	apiPort := "8080" // Default API port
	if len(parts) == 2 {
		// Use P2P port + 1000 for API
		// 3000 -> 4000, 3001 -> 4001, etc.
		port := parts[1]
		switch port {
		case "3000":
			apiPort = "4000"
		case "3001":
			apiPort = "4001"
		case "3002":
			apiPort = "4002"
		case "3003":
			apiPort = "4003"
		default:
			apiPort = "4000"
		}
	}

	apiServer := api.NewServer(bc, wallets, apiPort)

	server := &Server{
		Address:         address,
		Blockchain:      bc,
		Peers:           NewPeerList(),
		IsMining:        false,
		stopMining:      make(chan bool),
		miningInterrupt: make(chan bool, 10), // Buffered to not block
		APIServer:       apiServer,
		Wallets:         wallets,
	}

	// Set network server reference in API for broadcasting transactions
	apiServer.SetNetworkServer(server)

	return server
}

// Start starts the network server
func (s *Server) Start() error {
	// Use environment variable NODE_ADDR for P2P identification (Docker)
	// Fall back to s.Address for standalone mode
	if envAddr := os.Getenv("NODE_ADDR"); envAddr != "" {
		nodeAddress = envAddr
		log.Printf("Using P2P address from env: %s", nodeAddress)
	} else {
		nodeAddress = s.Address
	}

	// Start API server in background
	go func() {
		log.Printf("Starting API server...")
		if err := s.APIServer.Start(); err != nil {
			log.Printf("API server error: %v", err)
		}
	}()

	ln, err := net.Listen(protocol, s.Address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer ln.Close()

	log.Printf("Node server started on %s", s.Address)
	log.Printf("Node identifies as: %s", nodeAddress)

	// Connect to seed nodes if not seed
	seedNode := knownNodes[0]
	if nodeAddress != seedNode {
		log.Printf("Connecting to seed node: %s", seedNode)
		s.sendVersion(seedNode)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// StartMining enables mining on this node
func (s *Server) StartMining(address string) {
	s.IsMining = true
	miningAddress = address
	log.Printf("Mining enabled. Rewards will go to %s", address)

	// Start continuous mining loop
	go s.miningLoop()
}

// miningLoop continuously mines new blocks
// Real PoW mining - no timers, works continuously until finding valid block
func (s *Server) miningLoop() {
	log.Println("🔨 Starting continuous mining (real PoW)...")

	for {
		select {
		case <-s.stopMining:
			log.Println("Mining stopped")
			return
		default:
			// Check if we have transactions to mine (or just mine empty block with coinbase)
			mempoolMux.RLock()
			hasTxs := len(memoryPool) > 0
			mempoolMux.RUnlock()

			if hasTxs || true { // Always mine (even empty blocks with coinbase)
				s.mineTransactions()
			} else {
				// Small sleep to avoid CPU spinning when no txs
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// handleConnection handles incoming connections
func (s *Server) handleConnection(conn net.Conn) {
	request, err := io.ReadAll(conn)
	if err != nil {
		log.Printf("Error reading request: %v", err)
		conn.Close()
		return
	}

	// Validate request length
	if len(request) < commandLength {
		log.Printf("Request too short: %d bytes", len(request))
		conn.Close()
		return
	}

	command := BytesToCmd(request[:commandLength])
	log.Printf("Received %s command", command)

	switch command {
	case CmdVersion:
		s.handleVersion(request, conn)
	case CmdGetBlocks:
		s.handleGetBlocks(request, conn)
	case CmdInv:
		s.handleInv(request, conn)
	case CmdGetData:
		s.handleGetData(request, conn)
	case CmdBlock:
		s.handleBlock(request, conn)
	case CmdTx:
		s.handleTx(request, conn)
	case CmdAddr:
		s.handleAddr(request, conn)
	case CmdPing:
		s.handlePing(conn)
	default:
		log.Printf("Unknown command: %s", command)
	}

	conn.Close()
}

// sendVersion sends version message to peer
func (s *Server) sendVersion(addr string) {
	bestHeight := s.getBestHeight()

	payload := GobEncode(Version{
		Version:    version,
		BestHeight: bestHeight,
		AddrFrom:   nodeAddress,
	})

	request := append(CmdToBytes(CmdVersion), payload...)
	s.sendData(addr, request)
}

// handleVersion handles version message
func (s *Server) handleVersion(request []byte, conn net.Conn) {
	var buff bytes.Buffer
	var payload Version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Printf("Error decoding version: %v", err)
		return
	}

	bestHeight := s.getBestHeight()
	otherHeight := payload.BestHeight

	// Add peer
	s.Peers.Add(payload.AddrFrom, conn)

	log.Printf("Received version from %s: height %d (ours: %d)",
		payload.AddrFrom, otherHeight, bestHeight)

	if bestHeight < otherHeight {
		log.Printf("Peer has longer chain, requesting blocks...")
		s.sendGetBlocks(payload.AddrFrom)
	} else if bestHeight > otherHeight {
		s.sendVersion(payload.AddrFrom)
	}

	if !s.nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
		log.Printf("Added new peer: %s (total peers: %d)", payload.AddrFrom, len(knownNodes))
	}

	// Share our peer list with the new node
	s.sendAddr(payload.AddrFrom)
}

// sendGetBlocks sends getblocks request
func (s *Server) sendGetBlocks(addr string) {
	payload := GobEncode(GetBlocks{AddrFrom: nodeAddress})
	request := append(CmdToBytes(CmdGetBlocks), payload...)
	s.sendData(addr, request)
}

// handleGetBlocks handles getblocks request
func (s *Server) handleGetBlocks(request []byte, conn net.Conn) {
	var buff bytes.Buffer
	var payload GetBlocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Printf("Error decoding getblocks: %v", err)
		return
	}

	blocks := s.getBlocks()
	s.sendInv(payload.AddrFrom, InvTypeBlock, blocks)
}

// sendInv sends inventory message
func (s *Server) sendInv(addr, kind string, items [][]byte) {
	inventory := Inv{
		AddrFrom: nodeAddress,
		Type:     kind,
		Items:    items,
	}
	payload := GobEncode(inventory)
	request := append(CmdToBytes(CmdInv), payload...)

	s.sendData(addr, request)
}

// handleInv handles inventory message
func (s *Server) handleInv(request []byte, conn net.Conn) {
	var buff bytes.Buffer
	var payload Inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Printf("Error decoding inv: %v", err)
		return
	}

	log.Printf("Received inventory with %d %s", len(payload.Items), payload.Type)

	if payload.Type == InvTypeBlock {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		s.sendGetData(payload.AddrFrom, InvTypeBlock, blockHash)

		var newInTransit [][]byte
		for _, b := range blocksInTransit {
			if !bytes.Equal(b, blockHash) {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == InvTypeTx {
		txID := payload.Items[0]

		if memoryPool[hex.EncodeToString(txID)].ID == nil {
			s.sendGetData(payload.AddrFrom, InvTypeTx, txID)
		}
	}
}

// sendGetData sends getdata request
func (s *Server) sendGetData(addr, kind string, id []byte) {
	payload := GobEncode(GetData{
		AddrFrom: nodeAddress,
		Type:     kind,
		ID:       id,
	})
	request := append(CmdToBytes(CmdGetData), payload...)

	s.sendData(addr, request)
}

// handleGetData handles getdata request
func (s *Server) handleGetData(request []byte, conn net.Conn) {
	var buff bytes.Buffer
	var payload GetData

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Printf("Error decoding getdata: %v", err)
		return
	}

	if payload.Type == InvTypeBlock {
		block, err := s.getBlock(payload.ID)
		if err != nil {
			log.Printf("Error getting block: %v", err)
			return
		}

		s.sendBlock(payload.AddrFrom, block)
	}

	if payload.Type == InvTypeTx {
		txID := hex.EncodeToString(payload.ID)
		tx := memoryPool[txID]

		s.sendTx(payload.AddrFrom, tx)
	}
}

// sendBlock sends block to peer
func (s *Server) sendBlock(addr string, block *blockchain.Block) {
	data := BlockMsg{
		AddrFrom: nodeAddress,
		Block:    block.Serialize(),
	}
	payload := GobEncode(data)
	request := append(CmdToBytes(CmdBlock), payload...)

	s.sendData(addr, request)
}

// sendAddr sends known peer addresses to a node
func (s *Server) sendAddr(addr string) {
	data := Addr{AddrList: knownNodes}
	payload := GobEncode(data)
	request := append(CmdToBytes(CmdAddr), payload...)

	s.sendData(addr, request)
}

// handleBlock handles block message
func (s *Server) handleBlock(request []byte, conn net.Conn) {
	var buff bytes.Buffer
	var payload BlockMsg

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Printf("Error decoding block: %v", err)
		return
	}

	blockData := payload.Block
	block := blockchain.Deserialize(blockData)

	log.Printf("Received a new block height %d", block.Height)

	// Add block to blockchain (validation should be done here)
	s.addBlock(block)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		s.sendGetData(payload.AddrFrom, InvTypeBlock, blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := blockchain.UTXOSet{Blockchain: s.Blockchain}
		UTXOSet.Reindex()
	}
}

// sendTx sends transaction to peer
func (s *Server) sendTx(addr string, tx *blockchain.Transaction) {
	data := TxMsg{
		AddrFrom:    nodeAddress,
		Transaction: tx.Serialize(),
	}
	payload := GobEncode(data)
	request := append(CmdToBytes(CmdTx), payload...)

	s.sendData(addr, request)
}

// handleTx handles transaction message
func (s *Server) handleTx(request []byte, conn net.Conn) {
	var buff bytes.Buffer
	var payload TxMsg

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Printf("Error decoding tx: %v", err)
		return
	}

	txData := payload.Transaction
	tx := blockchain.DeserializeTransaction(txData)

	mempoolMux.Lock()
	memoryPool[hex.EncodeToString(tx.ID)] = &tx
	mempoolMux.Unlock()

	log.Printf("📥 Received transaction %x (mempool size: %d)", tx.ID, len(memoryPool))

	// Mining happens automatically every 60 seconds via miningLoop
}

// handleAddr handles addr message
func (s *Server) handleAddr(request []byte, conn net.Conn) {
	var buff bytes.Buffer
	var payload Addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Printf("Error decoding addr: %v", err)
		return
	}

	for _, addr := range payload.AddrList {
		if !s.nodeIsKnown(addr) && addr != nodeAddress {
			knownNodes = append(knownNodes, addr)
			log.Printf("🌐 Discovered new peer: %s (total: %d)", addr, len(knownNodes))

			// Try to connect to the new peer
			go func(peerAddr string) {
				s.sendVersion(peerAddr)
			}(addr)
		}
	}
}

// handlePing handles ping message
func (s *Server) handlePing(conn net.Conn) {
	payload := GobEncode(Pong{})
	request := append(CmdToBytes(CmdPong), payload...)
	conn.Write(request)
}

// BroadcastTx broadcasts transaction to all known peers
func (s *Server) BroadcastTx(tx *blockchain.Transaction) {
	for _, node := range knownNodes {
		if node != nodeAddress {
			s.sendTx(node, tx)
		}
	}
}

// BroadcastBlock broadcasts block to all known peers
func (s *Server) BroadcastBlock(block *blockchain.Block) {
	log.Printf("📡 Broadcasting block %d to %d peers", block.Height, len(knownNodes)-1)
	for _, node := range knownNodes {
		if node != nodeAddress {
			log.Printf("   → Sending to %s", node)
			s.sendInv(node, InvTypeBlock, [][]byte{block.Hash})
		}
	}
}

// sendData sends data to address
func (s *Server) sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		log.Printf("Error connecting to %s: %v", addr, err)
		s.removeNode(addr)
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Printf("Error sending data to %s: %v", addr, err)
	}
}

// Helper functions

func (s *Server) getBestHeight() int {
	var lastBlock *blockchain.Block
	iter := s.Blockchain.Iterator()
	lastBlock = iter.Next()
	return lastBlock.Height
}

func (s *Server) getBlocks() [][]byte {
	var blocks [][]byte
	iter := s.Blockchain.Iterator()

	for {
		block := iter.Next()
		blocks = append(blocks, block.Hash)

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return blocks
}

func (s *Server) getBlock(hash []byte) (*blockchain.Block, error) {
	iter := s.Blockchain.Iterator()

	for {
		block := iter.Next()

		if bytes.Equal(block.Hash, hash) {
			return block, nil
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return nil, fmt.Errorf("block not found")
}

func (s *Server) addBlock(block *blockchain.Block) {
	// Get current best height
	currentHeight := s.Blockchain.GetBestHeight()

	// Check if block is next in sequence
	if block.Height == currentHeight+1 {
		// Validate block using the difficulty stored in the block
		pow := blockchain.NewProofWithDifficulty(block, block.Difficulty)

		// Debug: print all InitData components
		pow.DebugInitData(block.Nonce)

		// Recalculate hash
		data := pow.InitData(block.Nonce)
		log.Printf("🔍 Raw InitData (len=%d): %x", len(data), data)
		hash := sha256.Sum256(data)

		if !pow.Validate() {
			txHash := block.HashTransactions()
			log.Printf("❌ Invalid block received (PoW failed)")
			log.Printf("   Block Height: %d, Hash: %x", block.Height, block.Hash)
			log.Printf("   Recalculated Hash: %x", hash)
			log.Printf("   Hashes match: %v", bytes.Equal(block.Hash, hash[:]))
			log.Printf("   TxHash: %x", txHash)
			log.Printf("   PrevHash: %x", block.PrevHash)
			log.Printf("   Nonce: %d, Difficulty: %d, Timestamp: %d", block.Nonce, block.Difficulty, block.Timestamp)
			log.Printf("   pow.Difficulty: %d, pow.Block.Difficulty: %d", pow.Difficulty, pow.Block.Difficulty)
			log.Printf("   Num Transactions: %d", len(block.Transactions))
			log.Printf("   ❌ Block rejected!")
			return
		}
		log.Printf("✅ Block PoW validated successfully (difficulty: %d)", block.Difficulty)

		// Add block to blockchain
		err := s.Blockchain.Database.Put(block.Hash, block.Serialize(), nil)
		if err != nil {
			log.Printf("Error storing block: %v", err)
			return
		}

		err = s.Blockchain.Database.Put([]byte("lh"), block.Hash, nil)
		if err != nil {
			log.Printf("Error updating last hash: %v", err)
			return
		}

		s.Blockchain.LastHash = block.Hash
		log.Printf("✅ Block accepted! Height: %d, Hash: %x", block.Height, block.Hash)

		// Update UTXO set
		UTXOSet := blockchain.UTXOSet{Blockchain: s.Blockchain}
		UTXOSet.Reindex()

		// Interrupt any ongoing mining (non-blocking)
		select {
		case s.miningInterrupt <- true:
			log.Println("🛑 Signaled mining interrupt - new block accepted")
		default:
			// Channel full or no miner active, ignore
		}

	} else if block.Height > currentHeight+1 {
		// We're missing blocks, request them
		log.Printf("⚠️  Missing blocks! Our height: %d, received: %d", currentHeight, block.Height)
		// This should trigger a full sync, but for now just log
	} else {
		log.Printf("ℹ️  Block %d already known or outdated", block.Height)
	}
}

func (s *Server) nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}
	return false
}

func (s *Server) removeNode(addr string) {
	var newNodes []string
	for _, node := range knownNodes {
		if node != addr {
			newNodes = append(newNodes, node)
		}
	}
	knownNodes = newNodes
}

func (s *Server) mineTransactions() {
	mempoolMux.Lock()

	var txs []*blockchain.Transaction

	// Collect valid transactions from mempool
	for id := range memoryPool {
		tx := memoryPool[id]
		if s.Blockchain.VerifyTransaction(tx) {
			txs = append(txs, tx)
		}
	}

	// Get current height for coinbase reward calculation
	newHeight := s.Blockchain.GetBestHeight() + 1
	cbTx := blockchain.CoinbaseTX(miningAddress, "", newHeight)
	txs = append(txs, cbTx)

	// Always mine, even if only coinbase transaction exists
	if len(txs) == 1 {
		log.Println("⛏️  Mining block with only coinbase transaction (reward)")
	} else {
		log.Printf("⛏️  Mining block with %d transaction(s) + coinbase", len(txs)-1)
	}

	// Unlock during mining (long operation)
	mempoolMux.Unlock()

	// Mine with interrupt support
	newBlock := s.Blockchain.MineBlockWithInterrupt(txs, s.miningInterrupt)

	// If block is nil, mining was interrupted by a new block from network
	if newBlock == nil {
		log.Println("⚠️  Mining interrupted - new block received from network")
		return
	}

	// Lock again for mempool cleanup
	mempoolMux.Lock()
	defer mempoolMux.Unlock()

	UTXOSet := blockchain.UTXOSet{Blockchain: s.Blockchain}
	UTXOSet.Reindex()

	log.Printf("✅ New block mined! Height: %d, Hash: %x", newBlock.Height, newBlock.Hash)

	// Clear mined transactions from mempool
	for _, tx := range txs {
		if !tx.IsCoinbase() { // Don't try to delete coinbase from mempool
			txID := hex.EncodeToString(tx.ID)
			delete(memoryPool, txID)
		}
	}

	// Broadcast new block
	s.BroadcastBlock(newBlock)
}

// GetKnownNodes returns list of known nodes
func GetKnownNodes() []string {
	return knownNodes
}

// AddKnownNode adds a node to known nodes
func AddKnownNode(addr string) {
	for _, node := range knownNodes {
		if node == addr {
			return
		}
	}
	knownNodes = append(knownNodes, addr)
}
