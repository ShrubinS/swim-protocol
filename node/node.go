package node

import (
	"github.com/google/uuid"
	"github.com/shrubins/swim-protocol/group"
)

type Node struct {
	ID        uuid.UUID
	groupView *group.Group
}

func NewLeaderNode() *Node {
	newNodeUUID, err := uuid.NewRandom()
	if err != nil {
		panic("Error generating Node UUID:\n" + err.Error())
	}
	newGroupUUID, err := uuid.NewRandom()
	if err != nil {
		panic("Error generating Group UUID:\n" + err.Error())
	}
	return &Node{
		newNodeUUID,
		group.NewGroup(newGroupUUID, 20),
	}
}
