package node

import (
	"github.com/google/uuid"
)

type Message struct {
	GroupID uuid.UUID `json:"group"`
	NodeID  uuid.UUID `json:"node"`
	Members []string  `json:"members"`
}

// func (message Message) encode() ([]byte, error) {
// 	return json.Marshal(message)
// }

// func
