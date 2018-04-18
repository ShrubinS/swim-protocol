package group

import (
	"github.com/google/uuid"
	"github.com/shrubins/swim-protocol/message"
)

//
type Group struct {
	GroupID uuid.UUID
	View    []message.Member
}

//
func NewGroup(id uuid.UUID, self message.Member) *Group {
	return &Group{
		id,
		[]message.Member{
			self,
		},
	}
}

func MakeGroup(id uuid.UUID, view []message.Member) *Group {
	return &Group{
		id,
		view,
	}
}
