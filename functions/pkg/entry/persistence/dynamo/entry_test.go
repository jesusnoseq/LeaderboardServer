package dynamo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/models"
	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/persistence/dynamo"
	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/utils"
)

func TestGetAllEntries(t *testing.T) {
	testTableName := fmt.Sprintf("LeaderboardTest-%d", time.Now().Unix())
	dbClient := utils.NewLocalDynamoClient()
	entryDAO := dynamo.NewEntryDAO(testTableName, dbClient, 10*time.Second)
	ctx := context.Background()
	utils.MustCreateTableSync(ctx, dbClient, testTableName)
	result, err := entryDAO.GetEntries(ctx)
	if err != nil {
		t.Errorf("err got %v, want nil", err)
	}
	if len(result) != 0 {
		t.Errorf("output got %q, want empty", result)
	}
	newEntry := models.Entry{Name: "prueba", Score: 99}
	entryDAO.CreateEntry(ctx, newEntry)
	result, err = entryDAO.GetEntries(ctx)
	if err != nil {
		t.Errorf("err got %v, want nil", err)
	}
	if len(result) != 1 {
		t.Errorf("output got %+v, want empty", result)
	}
	if !(result[0].Name == newEntry.Name && result[0].Score == newEntry.Score) {
		t.Errorf("output got %+v, want %+v", result[0], newEntry)
	}
}
