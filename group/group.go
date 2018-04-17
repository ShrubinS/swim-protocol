package group

import (
	"github.com/google/uuid"
)

//
type Group struct {
	ID   uuid.UUID
	View []int
}

//
func NewGroup(id uuid.UUID, size int) *Group {
	return &Group{
		id,
		make([]int, size),
	}
}
