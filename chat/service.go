package chat

import (
	"encoding/json"
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
	if user, ok := s.session2user[session]; ok {
		delete(s.session2user, session)
		delete(s.id2session, user.ID)
	}
}

func (s *Service) HandleMessage(session *transport.Session, message []byte) {
	log.Println("Chat::HandleMessage: ", string(message))
	user, ok := s.session2user[session]

	// New user
	if !ok {
		var helloMsg HelloMessage
		if err := json.Unmarshal(message, &helloMsg); err != nil {
			log.Println("Chat::HandleMessage parse hello failed, ", err)
			session.Close()
			return
		}

		s.session2user[session] = &User{
			ID:   helloMsg.ID,
			Role: helloMsg.Role,
		}

		s.id2session[helloMsg.ID] = session

		return
	}

	var chatMsg ChatMessage
	if err := json.Unmarshal(message, &chatMsg); err != nil {
		log.Println("Chat::HandleMessage parse chat failed, ", err)
		session.Close()
		return
	}

	// Send back message.
	session.Send(message)

	if user.Role == "player" {
		chatMsg.Sender = user.ID
		if outMessage, err := json.MarshalIndent(chatMsg, "", "  "); err != nil {
			// Send error.
			log.Println("Chat::HandleMessage encode chat failed, ", err)
			return
		} else {
			// Broadcast to all supporters.
			isSent := false
			for userSession, userInfo := range s.session2user {
				if userInfo.Role == "supporter" {
					userSession.Send(outMessage)
					isSent = true
				}
			}

			if !isSent {
				log.Println("Chat::HandleMessage all supporters are offline")
			}
		}
	} else if user.Role == "supporter" {
		if receiver, ok := s.id2session[chatMsg.Receiver]; ok {
			receiver.Send(message)
		} else {
			// This user is not online.
			log.Println("Chat::HandleMessage receiver offline")
			return
		}
	} else {
		log.Println("Chat::HandleMessage unsupported role")
		session.Close()
	}

}
