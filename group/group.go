package group

import (
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/shrubins/swim-protocol/message"
)

//
type Group struct {
	RealGroup bool
	GroupID   uuid.UUID
	View      []message.Member
}

//
func NewGroup(id uuid.UUID, self message.Member) *Group {
	return &Group{
		true,
		id,
		[]message.Member{
			self,
		},
	}
}

func NotGroup(self message.Member) *Group {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Println("can't stand thisssss.....")
	}

	return &Group{
		false,
		id,
		[]message.Member{
			self,
		},
	}
}

func MakeGroup(id uuid.UUID, view []message.Member) *Group {
	return &Group{
		true,
		id,
		view,
	}
}

func (group *Group) AddToGroup(member message.Member) {
	// change View to map instead of list
	for _, m := range group.View {
		if member.Is(m) {
			return
		}
	}
	group.View = append(group.View, member)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!New members appended below...")
	group.String()
}

func (group *Group) String() {
	fmt.Println("Group is...")
	for _, m := range group.View {
		fmt.Printf("ID:%s\nAddress:%s\n", m.NodeID, m.Address)
	}
}
