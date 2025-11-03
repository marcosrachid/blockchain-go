package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/marcocsrachid/blockchain-go/internal/blockchain"
	"github.com/marcocsrachid/blockchain-go/internal/network"
)

func printUsage() {
	fmt.Println("Blockchain Node")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  blockchain createwallet              - Creates a new wallet")
	fmt.Println("  blockchain listaddresses             - Lists all wallet addresses")
	fmt.Println("  blockchain createblockchain -address ADDRESS  - Creates initial blockchain (internal use)")
	fmt.Println("  blockchain startnode [options]       - Starts the blockchain node")
	fmt.Println("")
	fmt.Println("Start Node Options:")
	fmt.Println("  -miner ADDRESS    Enable mining and send rewards to ADDRESS")
	fmt.Println("  -port PORT        Port to listen on (default: 3000)")
	fmt.Println("")
	fmt.Println("HTTP API will be available on port 4000+ (node port + 1000)")
	fmt.Println("")
	fmt.Println("API Endpoints:")
	fmt.Println("  GET  /api/balance/:address    - Get address balance")
	fmt.Println("  GET  /api/addresses           - List all addresses")
	fmt.Println("  POST /api/createwallet        - Create new wallet")
	fmt.Println("  POST /api/send                - Send transaction")
	fmt.Println("  GET  /api/height              - Get blockchain height")
	fmt.Println("  GET  /api/difficulty          - Get current difficulty")
	fmt.Println("  GET  /api/networkinfo         - Get network information")
	fmt.Println("  GET  /api/lastblock           - Get last block info")
	fmt.Println("  GET  /api/block/:hash         - Get block by hash")
}

// createWallet creates a new wallet
func createWallet() {
	wallets, err := blockchain.NewWallets()
	if err != nil {
		log.Printf("Warning: Could not load existing wallets: %v", err)
		wallets = &blockchain.Wallets{Wallets: make(map[string]*blockchain.Wallet)}
	}

	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n", address)
}

// listAddresses lists all addresses in the wallets
func listAddresses() {
	wallets, err := blockchain.NewWallets()
	if err != nil {
		log.Printf("Error loading wallets: %v", err)
		return
	}

	addresses := wallets.GetAllAddresses()

	if len(addresses) == 0 {
		fmt.Println("No addresses found. Create one with 'createwallet'")
		return
	}

	for _, address := range addresses {
		fmt.Println(address)
	}
}

// createBlockchain creates a new blockchain (for initial setup only)
func createBlockchain(address string) {
	if !blockchain.ValidateAddress(address) {
		log.Panic("Address is not valid")
	}
	chain := blockchain.InitBlockchain(address)
	defer chain.Database.Close()

	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	UTXOSet.Reindex()

	fmt.Println("Blockchain created successfully!")
}

// startNode starts a network node
func startNode(minerAddress, nodeAddress string) {
	fmt.Printf("Starting node %s\n", nodeAddress)

	if len(minerAddress) > 0 {
		if blockchain.ValidateAddress(minerAddress) {
			fmt.Printf("Mining enabled. Rewards will go to %s\n", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}

	// Check if blockchain exists
	var chain *blockchain.Blockchain
	if !blockchain.DBexists() {
		fmt.Println("No blockchain found. Node will sync from network peers.")
		fmt.Println("For now, create blockchain first using 'createblockchain' command.")
		os.Exit(1)
	}

	chain = blockchain.ContinueBlockchain(minerAddress)
	defer chain.Database.Close()

	// Load wallets for API
	wallets, err := blockchain.NewWallets()
	if err != nil {
		log.Printf("Warning: Could not load wallets: %v", err)
		wallets = &blockchain.Wallets{Wallets: make(map[string]*blockchain.Wallet)}
	}

	server := network.NewServer(nodeAddress, chain, wallets)

	if len(minerAddress) > 0 {
		server.StartMining(minerAddress)
	}

	// Start server (blocking)
	if err := server.Start(); err != nil {
		log.Panic(err)
	}
}

func main() {
	defer os.Exit(0)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "createwallet":
		createWallet()

	case "listaddresses":
		listAddresses()

	case "createblockchain":
		createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
		createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")

		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		createBlockchain(*createBlockchainAddress)

	case "startnode":
		startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
		startNodeMiner := startNodeCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")
		startNodePort := startNodeCmd.String("port", "3000", "Port to listen on")

		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

		nodeAddress := fmt.Sprintf("0.0.0.0:%s", *startNodePort)
		startNode(*startNodeMiner, nodeAddress)

	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}
