package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/marcocsrachid/blockchain-go/internal/blockchain"
)

// Server represents the HTTP API server
type Server struct {
	Blockchain    *blockchain.Blockchain
	Wallets       *blockchain.Wallets
	Port          string
	NetworkServer interface{} // Reference to network server for broadcasting
}

// Response structures
type BalanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type AddressesResponse struct {
	Addresses []string `json:"addresses"`
}

type BlockResponse struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prev_hash"`
	Height       int    `json:"height"`
	Timestamp    int64  `json:"timestamp"`
	Transactions int    `json:"transactions"`
	Nonce        int    `json:"nonce"`
}

type SendRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type SendResponse struct {
	Success bool   `json:"success"`
	TxID    string `json:"tx_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type DifficultyResponse struct {
	Difficulty      int    `json:"difficulty"`
	Target          string `json:"target"`
	HashRate        string `json:"hash_rate_info"`
	TargetBlockTime int    `json:"target_block_time_seconds"`
}

type NetworkInfoResponse struct {
	Height        int `json:"height"`
	Difficulty    int `json:"difficulty"`
	TotalSupply   int `json:"total_supply"`
	MaxSupply     int `json:"max_supply"`
	CurrentReward int `json:"current_block_reward"`
	NextHalving   int `json:"blocks_until_halving"`
}

type LastBlockResponse struct {
	Hash         string `json:"hash"`
	Height       int    `json:"height"`
	Timestamp    int64  `json:"timestamp"`
	Transactions int    `json:"transactions"`
	Nonce        int    `json:"nonce"`
	PrevHash     string `json:"prev_hash"`
}

type CreateWalletResponse struct {
	Address string `json:"address"`
	Message string `json:"message"`
}

// NewServer creates a new API server
func NewServer(chain *blockchain.Blockchain, wallets *blockchain.Wallets, port string) *Server {
	return &Server{
		Blockchain:    chain,
		Wallets:       wallets,
		Port:          port,
		NetworkServer: nil, // Will be set later to avoid circular dependency
	}
}

// SetNetworkServer sets the network server reference for broadcasting transactions
func (s *Server) SetNetworkServer(networkServer interface{}) {
	s.NetworkServer = networkServer
}

// Start starts the HTTP API server
func (s *Server) Start() error {
	http.HandleFunc("/api/balance/", s.handleGetBalance)
	http.HandleFunc("/api/addresses", s.handleGetAddresses)
	http.HandleFunc("/api/createwallet", s.handleCreateWallet)
	http.HandleFunc("/api/send", s.handleSend)
	http.HandleFunc("/api/height", s.handleGetHeight)
	http.HandleFunc("/api/difficulty", s.handleGetDifficulty)
	http.HandleFunc("/api/networkinfo", s.handleGetNetworkInfo)
	http.HandleFunc("/api/lastblock", s.handleGetLastBlock)
	http.HandleFunc("/api/block/", s.handleGetBlockByHash)
	http.HandleFunc("/health", s.handleHealth)

	addr := fmt.Sprintf(":%s", s.Port)
	log.Printf("API server started on http://0.0.0.0%s", addr)
	return http.ListenAndServe(addr, nil)
}

// handleGetBalance returns the balance of an address
// GET /api/balance/:address
func (s *Server) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract address from URL path
	address := r.URL.Path[len("/api/balance/"):]
	if address == "" {
		s.sendError(w, "Address is required", http.StatusBadRequest)
		return
	}

	// Validate address
	if !blockchain.ValidateAddress(address) {
		s.sendError(w, "Invalid address format", http.StatusBadRequest)
		return
	}

	// Get balance
	pubKeyHash := blockchain.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	UTXOSet := blockchain.UTXOSet{Blockchain: s.Blockchain}
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	balance := 0
	for _, out := range UTXOs {
		balance += out.Value
	}

	response := BalanceResponse{
		Address: address,
		Balance: balance,
	}

	s.sendJSON(w, response, http.StatusOK)
}

// handleGetAddresses returns all wallet addresses
// GET /api/addresses
func (s *Server) handleGetAddresses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	addresses := s.Wallets.GetAllAddresses()

	response := AddressesResponse{
		Addresses: addresses,
	}

	s.sendJSON(w, response, http.StatusOK)
}

// handleCreateWallet creates a new wallet and returns the address
// POST /api/createwallet
func (s *Server) handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create new wallet
	address := s.Wallets.AddWallet()

	// Save wallets to file
	s.Wallets.SaveFile()

	response := CreateWalletResponse{
		Address: address,
		Message: "Wallet created successfully",
	}

	log.Printf("✅ New wallet created: %s", address)
	s.sendJSON(w, response, http.StatusCreated)
}

// handleGetBlockByHash returns a specific block by its hash
// GET /api/block/:hash
func (s *Server) handleGetBlockByHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract hash from URL path
	hashStr := r.URL.Path[len("/api/block/"):]
	if hashStr == "" {
		s.sendError(w, "Block hash is required", http.StatusBadRequest)
		return
	}

	// Decode hex hash
	blockHash, err := hex.DecodeString(hashStr)
	if err != nil {
		s.sendError(w, "Invalid block hash format", http.StatusBadRequest)
		return
	}

	// Get block from blockchain
	block, err := s.Blockchain.GetBlock(blockHash)
	if err != nil {
		s.sendError(w, "Block not found", http.StatusNotFound)
		return
	}

	response := BlockResponse{
		Hash:         fmt.Sprintf("%x", block.Hash),
		PrevHash:     fmt.Sprintf("%x", block.PrevHash),
		Height:       block.Height,
		Timestamp:    block.Timestamp,
		Transactions: len(block.Transactions),
		Nonce:        block.Nonce,
	}

	s.sendJSON(w, response, http.StatusOK)
}

// handleSend creates and broadcasts a new transaction
// POST /api/send
func (s *Server) handleSend(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== HANDLER SEND CALLED ===")

	if r.Method != http.MethodPost {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate inputs
	if req.From == "" || req.To == "" || req.Amount <= 0 {
		s.sendError(w, "From, To, and Amount are required", http.StatusBadRequest)
		return
	}

	if !blockchain.ValidateAddress(req.From) {
		s.sendError(w, "Invalid 'from' address", http.StatusBadRequest)
		return
	}

	if !blockchain.ValidateAddress(req.To) {
		s.sendError(w, "Invalid 'to' address", http.StatusBadRequest)
		return
	}

	// Get wallet to verify it exists
	wallet := s.Wallets.GetWallet(req.From)

	// Check if wallet exists by verifying if public key is empty
	if len(wallet.PublicKey) == 0 {
		s.sendError(w, "Wallet not found for 'from' address", http.StatusNotFound)
		return
	}

	log.Printf("🔵 API: Received send request - From: %s, To: %s, Amount: %d", req.From, req.To, req.Amount)

	// Create transaction using addresses
	tx := blockchain.NewTransaction(req.From, req.To, req.Amount, s.Blockchain)
	if tx == nil {
		log.Printf("❌ API: Transaction creation failed - insufficient funds")
		s.sendError(w, "Failed to create transaction - insufficient funds", http.StatusBadRequest)
		return
	}

	log.Printf("✅ API: Transaction created successfully: %x", tx.ID)

	// Add transaction to local mempool first
	if s.NetworkServer != nil {
		// Type assert to add to local mempool
		type MempoolManager interface {
			AddToMempool(tx *blockchain.Transaction)
			BroadcastTx(tx *blockchain.Transaction)
		}
		if manager, ok := s.NetworkServer.(MempoolManager); ok {
			manager.AddToMempool(tx)
			log.Printf("📥 API: Added transaction to local mempool")
			manager.BroadcastTx(tx)
			log.Printf("📤 API: Transaction broadcasted: %x", tx.ID)
		} else {
			log.Printf("⚠️  API: NetworkServer does not implement required methods!")
		}
	} else {
		log.Printf("⚠️  API: NetworkServer is nil - transaction will NOT be broadcasted!")
	}

	response := SendResponse{
		Success: true,
		TxID:    fmt.Sprintf("%x", tx.ID),
	}

	log.Printf("🔵 API: Sending response to client")
	s.sendJSON(w, response, http.StatusOK)
}

// handleGetHeight returns the current blockchain height
// GET /api/height
func (s *Server) handleGetHeight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	height := s.Blockchain.GetBestHeight()

	response := map[string]int{
		"height": height,
	}

	s.sendJSON(w, response, http.StatusOK)
}

// handleGetLastBlock returns information about the last block
// GET /api/lastblock
func (s *Server) handleGetLastBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lastBlock := s.Blockchain.GetLastBlock()

	response := LastBlockResponse{
		Hash:         fmt.Sprintf("%x", lastBlock.Hash),
		Height:       lastBlock.Height,
		Timestamp:    lastBlock.Timestamp,
		Transactions: len(lastBlock.Transactions),
		Nonce:        lastBlock.Nonce,
		PrevHash:     fmt.Sprintf("%x", lastBlock.PrevHash),
	}

	s.sendJSON(w, response, http.StatusOK)
}

// handleGetDifficulty returns the current network difficulty
// GET /api/difficulty
func (s *Server) handleGetDifficulty(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := DifficultyResponse{
		Difficulty:      blockchain.Difficulty,
		Target:          fmt.Sprintf("2^(256-%d) = %d leading zeros required", blockchain.Difficulty, blockchain.Difficulty),
		HashRate:        "Higher difficulty = more computational work required",
		TargetBlockTime: 60, // 1 minute target
	}

	s.sendJSON(w, response, http.StatusOK)
}

// handleGetNetworkInfo returns comprehensive network information
// GET /api/networkinfo
func (s *Server) handleGetNetworkInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	height := s.Blockchain.GetBestHeight()
	currentReward := blockchain.GetBlockReward(height)

	// Calculate blocks until next halving
	blocksUntilHalving := blockchain.HalvingInterval - (height % blockchain.HalvingInterval)

	// Estimate current supply (simplified - doesn't account for lost coins)
	// This is an approximation
	totalSupply := calculateTotalSupply(height)

	response := NetworkInfoResponse{
		Height:        height,
		Difficulty:    blockchain.Difficulty,
		TotalSupply:   totalSupply,
		MaxSupply:     blockchain.MaxSupply,
		CurrentReward: currentReward,
		NextHalving:   blocksUntilHalving,
	}

	s.sendJSON(w, response, http.StatusOK)
}

// calculateTotalSupply estimates the total supply based on current height
func calculateTotalSupply(height int) int {
	totalSupply := 0
	currentReward := blockchain.InitialSubsidy
	blocksProcessed := 0

	for blocksProcessed <= height && currentReward > 0 {
		blocksInThisEra := blockchain.HalvingInterval
		if blocksProcessed+blocksInThisEra > height {
			blocksInThisEra = height - blocksProcessed + 1
		}

		totalSupply += blocksInThisEra * currentReward
		blocksProcessed += blockchain.HalvingInterval
		currentReward = currentReward / 2
	}

	return totalSupply
}

// handleHealth is a health check endpoint
// GET /health
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	s.sendJSON(w, response, http.StatusOK)
}

// Helper functions

func (s *Server) sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func (s *Server) sendError(w http.ResponseWriter, message string, status int) {
	response := ErrorResponse{
		Error: message,
	}
	s.sendJSON(w, response, status)
}

// ParseIntParam parses an integer parameter from the request
func ParseIntParam(r *http.Request, param string, defaultValue int) int {
	value := r.URL.Query().Get(param)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}
