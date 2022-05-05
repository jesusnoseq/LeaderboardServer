package persistence

import (
	"log"
	"time"

	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/persistence/dynamo"
	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/persistence/memory"
	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/utils"
)

type DatabaseEngine string

const (
	Dynamo DatabaseEngine = "DYNAMO"
	Memory DatabaseEngine = "MEMORY"
)

func GetDatabaseEngine(des string) DatabaseEngine {
	m := map[string]DatabaseEngine{
		string(Dynamo): Dynamo,
		string(Memory): Memory,
	}

	return m[des]
}

func GetEntryDAO(de DatabaseEngine) EntryDAO {
	log.Printf("Getting entry dao of with %v driver\n", de)
	switch de {
	case Dynamo:
		return dynamo.NewEntryDAO(dynamo.DEFAULT_ENTRY_TABLE_NAME,
			utils.NewDynamoClient(), 5*time.Second)
	case Memory:
		return memory.NewEntryDao()
	}
	return nil
}
