// Right now ws using JSON for message encoding/decoding
// In the future, we can add other encoding formats like text-based or binary for better performance

package ws

import "encoding/json"

type Message struct {
	Event   string `json:"event"`
	Payload any    `json:"payload"`
}

func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}

type MessageSchema struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}
