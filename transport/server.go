package transport

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Server maintains the set of active sessions and handles all network
// logic.
type Server struct {
	// Registered sessions.
	sessions map[*Session]bool

	// 1-1 messages channel.
	message chan []byte

	// Register requests from the sessions.
	register chan *Session

	// Unregister requests from sessions.
	unregister chan *Session
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NewServer creates a new server.
func NewServer() *Server {
	return &Server{
		message:    make(chan []byte),
		register:   make(chan *Session),
		unregister: make(chan *Session),
		sessions:   make(map[*Session]bool),
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
			}
		case message := <-s.message:
			for session := range s.sessions {
				select {
				case session.send <- message:
				default:
					close(session.send)
					delete(s.sessions, session)
				}
			}
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

// OnSessionClose handles session close event.
func (s *Server) OnSessionClose(session *Session) {
	s.unregister <- session
}

// OnSessionMessage handles new messages from client.
func (s *Server) OnSessionMessage(session *Session, message []byte) {
	log.Println(message)
}

// Serve start running server.
func (s *Server) ListenAndServe(addr string) {
	go s.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.handleWs(w, r)
	})

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}
