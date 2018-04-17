package node

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shrubins/swim-protocol/group"
)

const heartbeatInterval = 500 * time.Millisecond

type Node struct {
	nodeID    uuid.UUID
	groupView *group.Group
}

func NewLeaderNode() Node {
	newNodeUUID, err := uuid.NewRandom()
	if err != nil {
		panic("Error generating Node UUID:\n" + err.Error())
	}
	newGroupUUID, err := uuid.NewRandom()
	if err != nil {
		panic("Error generating Group UUID:\n" + err.Error())
	}
	return Node{
		newNodeUUID,
		group.NewGroup(newGroupUUID, 20),
	}
}

func (node Node) StartNode(done chan bool) {
	ticker := time.NewTicker(heartbeatInterval)
	go func() {
		for t := range ticker.C {
			// fmt.Println("Heartbeat", node.nodeID.String())
			t.String()
		}
	}()
	time.Sleep(5 * time.Second)
	ticker.Stop()
	fmt.Println("Terminating Node...")
	done <- true
}
