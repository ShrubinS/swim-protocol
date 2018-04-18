package message

import (
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

// func (message Message) encode() ([]byte, error) {
// 	return json.Marshal(message)
// }

// func
