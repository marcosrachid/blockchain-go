package blockchain

import (
	"log"
	"os"
)

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// DBexists checks if the database exists
func DBexists() bool {
	// LevelDB creates a CURRENT file in the database directory
	if _, err := os.Stat(dbPath + "/CURRENT"); os.IsNotExist(err) {
		return false
	}

	return true
}
