package transport

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ServerHandler is a handler for server network event.
type ServerHandler interface {
	HandleClosed(session *Session)
	HandleMessage(session *Session, message []byte)
}

// InMessage represents a message from a client.
type InMessage struct {
	sender  *Session
	content []byte
}

// Server maintains the set of active sessions and handles all network
// logic.
type Server struct {
	// Registered sessions.
	sessions map[*Session]bool

	// Incoming message channel.
	message chan *InMessage

	// Register requests from the sessions.
	register chan *Session

	// Unregister requests from sessions.
	unregister chan *Session

	// Server handler
	handler ServerHandler
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NewServer creates a new server.
func NewServer(h ServerHandler) *Server {
	return &Server{
		message:    make(chan *InMessage),
		register:   make(chan *Session),
		unregister: make(chan *Session),
		sessions:   make(map[*Session]bool),
		handler:    h,
	}
}

// run handles main application logic.
func (s *Server) run() {
	for {
		select {
		case session := <-s.register:
			s.sessions[session] = true
		case session := <-s.unregister:
			if _, ok := s.sessions[session]; ok {
				delete(s.sessions, session)
				close(session.send)

				s.handler.HandleClosed(session)
			}
		case message := <-s.message:
			s.handler.HandleMessage(message.sender, message.content)
			// for session := range s.sessions {
			// 	select {
			// 	case session.send <- message:
			// 	default:
			// 		close(session.send)
			// 		delete(s.sessions, session)
			// 	}
			// }
		}
	}
}

// handleWs handles websocket requests from the peer.
func (s *Server) handleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Creates a session
	session := NewSession(conn, s)

	// Register a new session.
	s.register <- session

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go session.HandleWrite()
	go session.HandleRead()
}

// OnClose handles session close event.
func (s *Server) OnClose(session *Session) {
	s.unregister <- session
}

// OnReceived handles new messages from client.
func (s *Server) OnReceived(session *Session, message []byte) {
	s.message <- &InMessage{sender: session, content: message}
}

// Serve start running server.
func (s *Server) ListenAndServe(addr string) {
	go s.run()

	http.HandleFunc("/chat/player", servePlayer)
	http.HandleFunc("/chat/supporter", serveSupporter)
	http.HandleFunc("/chat/ws", func(w http.ResponseWriter, r *http.Request) {
		s.handleWs(w, r)
	})

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func servePlayer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/chat/player" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "player.html")
}

func serveSupporter(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/chat/supporter" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "supporter.html")
}
