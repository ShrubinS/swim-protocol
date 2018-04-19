package message

import (
	"fmt"

	"github.com/google/uuid"
)

type Member struct {
	NodeID  uuid.UUID `json:"node"`
	Address string
}

type Message struct {
	GroupID uuid.UUID `json:"group"`
	Members []Member  `json:"members"`
}

func (m1 Member) Is(m2 Member) bool {
	if m1.NodeID == m2.NodeID {
		return true
	}
	return false
}

func (message Message) String() {
	fmt.Println("Message is.... ")
	for _, m := range message.Members {
		fmt.Printf("NodeID:%s\nAddress:%s\n", m.NodeID, m.Address)
	}
}

// func
