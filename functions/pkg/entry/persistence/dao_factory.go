package persistence

import (
	"log"

	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/persistence/memory"
)

func GetDao(modelName string, driver string) EntryDAO {
	log.Printf("Getting dao of %v with %v driver\n", modelName, driver)
	return &memory.MemoryDAO{}
}
