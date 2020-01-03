package chat

import (
	"log"

	"github.com/tinhatwork/go-chat/transport"
)

// Service is the core component handle all chat business logic.
type Service struct {
	session2user map[*transport.Session]*User
	id2session   map[string]*transport.Session
}

func NewService() *Service {
	return &Service{
		session2user: make(map[*transport.Session]*User),
		id2session:   make(map[string]*transport.Session),
	}
}

func (s *Service) HandleClosed(session *transport.Session) {
	log.Println("Chat::HandleClosed")
}

func (s *Service) HandleMessage(session *transport.Session, message []byte) {
	log.Println("Chat::HandleMessage: ", string(message))
	session.Close()
}
