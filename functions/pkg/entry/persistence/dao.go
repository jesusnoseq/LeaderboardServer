package persistence

import (
	"context"

	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/models"
)

type EntryDAO interface {
	CreateEntry(context.Context, models.Entry) (models.Entry, error)
	GetTopEntries(context.Context) ([]models.Entry, error)
}

// type LeaderboardDAO interface {
// 	CreateLeaderboard(context.Context, Entry) (Entry, error)
// 	GetLeaderboard(context.Context, Entry) (Entry, error)
// }
