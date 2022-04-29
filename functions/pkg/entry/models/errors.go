package models

import "errors"

var (
	ErrNotFound    = errors.New("entry item is not found")
	ErrInvalidHash = errors.New("invalid hash... Don't be evil")
)
