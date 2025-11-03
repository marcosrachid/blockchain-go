package blockchain

// Blockchain configuration constants
// All protocol parameters are centralized here for easy maintenance

const (
	// Mining and Reward Configuration
	InitialSubsidy  = 50       // Initial mining reward (similar to Bitcoin's initial 50 BTC)
	HalvingInterval = 210000   // Blocks until reward halving (same as Bitcoin ~4 years)
	MaxSupply       = 21000000 // Maximum supply of coins (21 million like Bitcoin)

	// Proof of Work Configuration
	Difficulty        = 22 // Mining difficulty (number of leading zeros required in hash)
	GenesisDifficulty = 16 // Lower difficulty for genesis block (faster initialization)

	// Genesis Block Configuration
	GenesisData = "First Transaction from Genesis" // Genesis block coinbase data

	// Database Configuration
	DBPath = "./tmp/blocks" // Default database path (can be overridden by env var)

	// Network Configuration (for reference)
	DefaultPort     = 3000 // Default network port
	ProtocolVersion = 1    // Protocol version for network communication
)

// GetBlockReward calculates the mining reward based on block height
// Implements halving every 210,000 blocks like Bitcoin
func GetBlockReward(height int) int {
	reward := InitialSubsidy

	// Calculate number of halvings
	halvings := height / HalvingInterval

	// Each halving divides reward by 2
	for i := 0; i < halvings; i++ {
		reward = reward / 2
	}

	// When reward becomes 0, no more coins are minted
	if reward < 1 {
		return 0
	}

	return reward
}

// GetMaxSupply returns the maximum supply
func GetMaxSupply() int {
	return MaxSupply
}

// GetTotalMinableBlocks returns the approximate number of blocks until max supply
func GetTotalMinableBlocks() int {
	totalBlocks := 0
	reward := InitialSubsidy

	for reward > 0 {
		totalBlocks += HalvingInterval
		reward = reward / 2
	}

	return totalBlocks
}
