package persistence

import (
	"log"
	"time"

	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/persistence/dynamo"
	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/persistence/memory"
	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/utils"
)

const (
	DYNAMO = "DYNAMO"
	MEMORY = "MEMORY"
)

func GetEntryDao(driver string) EntryDAO {
	log.Printf("Getting entry dao of with %v driver\n", driver)
	switch driver {
	case DYNAMO:
		return dynamo.NewEntryDAO(dynamo.DEFAULT_ENTRY_TABLE_NAME,
			utils.NewDynamoClient(), 5*time.Second)
	case MEMORY:
		return memory.NewEntryDao()
	}
	return nil
}
