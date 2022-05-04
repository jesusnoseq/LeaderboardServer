package memory

import (
	"context"

	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/models"
)

type EntryDAO struct {
	entries []models.Entry
}

func NewEntryDao() *EntryDAO {
	return &EntryDAO{}
}

func (e *EntryDAO) CreateEntry(ctx context.Context, entry models.Entry) (models.Entry, error) {
	e.entries = append(e.entries, entry)
	return entry, nil
}

func (e *EntryDAO) GetEntries(ctx context.Context) ([]models.Entry, error) {
	return e.entries, nil
}
