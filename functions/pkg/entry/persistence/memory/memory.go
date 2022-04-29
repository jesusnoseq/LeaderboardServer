package memory

import (
	"context"

	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/models"
)

type MemoryDAO struct {
	entries []models.Entry
}

func (m *MemoryDAO) CreateEntry(ctx context.Context, entry models.Entry) (models.Entry, error) {
	m.entries = append(m.entries, entry)
	return entry, nil
}

func (m *MemoryDAO) GetTopEntries(ctx context.Context) ([]models.Entry, error) {
	return m.entries, nil
}
