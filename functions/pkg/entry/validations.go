package entry

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/models"
)

func CheckValidEntry(e models.Entry, salt string) error {
	sum := sha256.Sum256([]byte(e.Name + strconv.Itoa(int(e.Score)) + salt))
	if hex.EncodeToString(sum[:]) != e.Hash {
		return models.ErrInvalidHash
	}
	return nil
}
