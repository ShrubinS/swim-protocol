package group

import (
	"github.com/google/uuid"
)

//
type Group struct {
	GroupID uuid.UUID
	View    []string
}

//
func NewGroup(id uuid.UUID, self string) *Group {
	return &Group{
		id,
		[]string{
			self,
		},
	}
}

func MakeGroup(id uuid.UUID, view []string) *Group {
	return &Group{
		id,
		view,
	}
}
