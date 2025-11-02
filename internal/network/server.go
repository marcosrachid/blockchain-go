package network

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/marcocsrachid/blockchain-go/internal/blockchain"
)

const (
	protocol      = "tcp"
	version       = 1
	commandLength = 12
)

var (
	nodeAddress     string
	miningAddress   string
	knownNodes      = []string{"localhost:3000"} // Seed node
	blocksInTransit = [][]byte{}
	memoryPool      = make(map[string]*blockchain.Transaction)
)

// Server represents the network server
type Server struct {
	Address    string
	Blockchain *blockchain.Blockchain
	Peers      *PeerList
	IsMining   bool
}

// NewServer creates a new network server
func NewServer(address string, bc *blockchain.Blockchain) *Server {
	return &Server{
		Address:    address,
		Blockchain: bc,
		Peers:      NewPeerList(),
		IsMining:   false,
	}
}

// Start starts the network server
func (s *Server) Start() error {
	nodeAddress = s.Address
	
	ln, err := net.Listen(protocol, s.Address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer ln.Close()

	log.Printf("Node server started on %s", s.Address)

	// Connect to seed nodes if not seed
	if s.Address != knownNodes[0] {
		s.sendVersion(knownNodes[0])
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
}

// handleConnection handles incoming connections
func (s *Server) handleConnection(conn net.Conn) {
	request, err := io.ReadAll(conn)
	if err != nil {
		log.Printf("Error reading request: %v", err)
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
		s.sendGetBlocks(payload.AddrFrom)
	} else if bestHeight > otherHeight {
		s.sendVersion(payload.AddrFrom)
	}

	if !s.nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
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
	memoryPool[hex.EncodeToString(tx.ID)] = &tx

	log.Printf("Received transaction %x", tx.ID)

	// If this node is mining and we have enough transactions, mine a block
	if s.IsMining && len(memoryPool) >= 2 {
		s.mineTransactions()
	}
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
		if !s.nodeIsKnown(addr) {
			knownNodes = append(knownNodes, addr)
			log.Printf("Added new peer: %s", addr)
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
	for _, node := range knownNodes {
		if node != nodeAddress {
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
	// In a real implementation, we should validate the block before adding
	// For now, we'll just log it
	log.Printf("Block %x added to chain", block.Hash)
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
	var txs []*blockchain.Transaction

	for id := range memoryPool {
		tx := memoryPool[id]
		if s.Blockchain.VerifyTransaction(tx) {
			txs = append(txs, tx)
		}
	}

	if len(txs) == 0 {
		log.Println("No valid transactions to mine")
		return
	}

	cbTx := blockchain.CoinbaseTX(miningAddress, "")
	txs = append(txs, cbTx)

	newBlock := s.Blockchain.MineBlock(txs)
	
	UTXOSet := blockchain.UTXOSet{Blockchain: s.Blockchain}
	UTXOSet.Reindex()

	log.Printf("New block mined: %x", newBlock.Hash)

	// Clear mined transactions from mempool
	for _, tx := range txs {
		txID := hex.EncodeToString(tx.ID)
		delete(memoryPool, txID)
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

