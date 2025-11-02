package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/marcocsrachid/blockchain-go/internal/blockchain"
	"github.com/marcocsrachid/blockchain-go/internal/network"
)

// CommandLine handles all CLI operations
type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - get the balance for an address")
	fmt.Println(" createblockchain -address ADDRESS - creates a blockchain and sends genesis reward to address")
	fmt.Println(" printchain - prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - send amount of coins")
	fmt.Println(" createwallet - creates a new wallet")
	fmt.Println(" listaddresses - lists the addresses in the wallet file")
	fmt.Println(" reindexutxo - rebuilds the UTXO set")
	fmt.Println(" startnode -miner ADDRESS -port PORT - starts a node with optional mining")
	fmt.Println(" addpeer -address ADDRESS - adds a peer to the network")
	fmt.Println(" peers - lists all known peers")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

// reindexUTXO rebuilds the UTXO set
func (cli *CommandLine) reindexUTXO() {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()
	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}

// listAddresses lists all addresses in the wallets
func (cli *CommandLine) listAddresses() {
	wallets, _ := blockchain.NewWallets()
	addresses := wallets.GetAllAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

// createWallet creates a new wallet
func (cli *CommandLine) createWallet() {
	wallets, _ := blockchain.NewWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n", address)
}

// printChain prints all blocks in the blockchain
func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

// createBlockchain creates a new blockchain
func (cli *CommandLine) createBlockchain(address string) {
	if !blockchain.ValidateAddress(address) {
		log.Panic("Address is not valid")
	}
	chain := blockchain.InitBlockchain(address)
	defer chain.Database.Close()

	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	UTXOSet.Reindex()

	fmt.Println("Finished!")
}

// getBalance gets the balance of an address
func (cli *CommandLine) getBalance(address string) {
	if !blockchain.ValidateAddress(address) {
		log.Panic("Address is not valid")
	}
	chain := blockchain.ContinueBlockchain(address)
	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	defer chain.Database.Close()

	balance := 0
	pubKeyHash := blockchain.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

// send sends coins from one address to another
func (cli *CommandLine) send(from, to string, amount int) {
	if !blockchain.ValidateAddress(from) {
		log.Panic("From address is not valid")
	}
	if !blockchain.ValidateAddress(to) {
		log.Panic("To address is not valid")
	}

	chain := blockchain.ContinueBlockchain(from)
	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	cbTx := blockchain.CoinbaseTX(from, "")
	txs := []*blockchain.Transaction{cbTx, tx}
	block := chain.MineBlock(txs)
	UTXOSet.Update(block)
	fmt.Println("Success!")
}

// startNode starts a network node
func (cli *CommandLine) startNode(minerAddress, nodeAddress string) {
	fmt.Printf("Starting node %s\n", nodeAddress)
	
	if len(minerAddress) > 0 {
		if blockchain.ValidateAddress(minerAddress) {
			fmt.Printf("Mining enabled. Rewards will go to %s\n", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	
	chain := blockchain.ContinueBlockchain(minerAddress)
	defer chain.Database.Close()
	
	server := network.NewServer(nodeAddress, chain)
	
	if len(minerAddress) > 0 {
		server.StartMining(minerAddress)
	}
	
	// Start server (blocking)
	if err := server.Start(); err != nil {
		log.Panic(err)
	}
}

// addPeer adds a peer to the network
func (cli *CommandLine) addPeer(peerAddress string) {
	network.AddKnownNode(peerAddress)
	fmt.Printf("Added peer: %s\n", peerAddress)
}

// listPeers lists all known peers
func (cli *CommandLine) listPeers() {
	nodes := network.GetKnownNodes()
	fmt.Println("Known peers:")
	for _, node := range nodes {
		fmt.Printf("  - %s\n", node)
	}
}

// Run processes command line arguments
func (cli *CommandLine) run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	addPeerCmd := flag.NewFlagSet("addpeer", flag.ExitOnError)
	peersCmd := flag.NewFlagSet("peers", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	startNodeMiner := startNodeCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")
	startNodePort := startNodeCmd.String("port", "3000", "Port to listen on")
	addPeerAddress := addPeerCmd.String("address", "", "Peer address to add")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addpeer":
		err := addPeerCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "peers":
		err := peersCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO()
	}

	if startNodeCmd.Parsed() {
		nodeAddress := fmt.Sprintf("localhost:%s", *startNodePort)
		cli.startNode(*startNodeMiner, nodeAddress)
	}

	if addPeerCmd.Parsed() {
		if *addPeerAddress == "" {
			addPeerCmd.Usage()
			runtime.Goexit()
		}
		cli.addPeer(*addPeerAddress)
	}

	if peersCmd.Parsed() {
		cli.listPeers()
	}
}

func main() {
	defer os.Exit(0)
	cmd := CommandLine{}
	cmd.run()
}
