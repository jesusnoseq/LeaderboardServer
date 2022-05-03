package persistence

import (
	"log"

	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/persistence/memory"
)

func GetEntryDao(driver string) EntryDAO {
	log.Printf("Getting entry dao of with %v driver\n", driver)
	return &memory.MemoryDAO{}
}
