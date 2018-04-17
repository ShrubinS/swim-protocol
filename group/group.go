package group

import (
	"github.com/google/uuid"
)

//
type Group struct {
	groupID uuid.UUID
	View    []string
}

//
func NewGroup(id uuid.UUID, size int) *Group {
	return &Group{
		id,
		make([]string, 0),
	}
}
