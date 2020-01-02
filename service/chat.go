package service

import (
	"log"

	"github.com/tinhatwork/go-chat/transport"
)

// Chat is the core component handle business logic of this service.
type Chat struct {
	session2user map[*transport.Session]*User
	id2session   map[string]*transport.Session
}

func NewService() *Chat {
	return &Chat{
		session2user: make(map[*transport.Session]*User),
		id2session:   make(map[string]*transport.Session),
	}
}

func (c *Chat) HandleClosed(session *transport.Session) {
	log.Println("Chat::HandleClosed")
}

func (c *Chat) HandleMessage(session *transport.Session, message []byte) {
	log.Println("Chat::HandleMessage: ", string(message))
}
