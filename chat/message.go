package chat

import "time"

type HelloMessage struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type ChatMessage struct {
	Sender   string    `json:"sender,omitempty"`
	Receiver string    `json:"receiver,omitempty"`
	Name     string    `json:"name"`
	Data     string    `json:"data"`
	Time     time.Time `json:"time"`
}
