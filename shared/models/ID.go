package models

import "github.com/google/uuid"

// NewID generates a time-ordered (v7) UUID for a primary key. Unlike v4,
// v7 embeds a millisecond timestamp so IDs sort in creation order,
// avoiding random B-tree insert spots and enabling keyset pagination
// (`WHERE id < :cursor ORDER BY id DESC`). Use this instead of
// uuid.New() for any app-side-generated entity ID.
func NewID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}
