package osc

import (
	"fmt"
	"log"

	gosc "github.com/hypebeast/go-osc/osc"
)

// Server represents an OSC server
type Server struct {
	server     *gosc.Server
	dispatcher *gosc.StandardDispatcher
	addr       string
	port       int
}

// NewServer creates a new OSC server
func NewServer(addr string, port int) *Server {
	dispatcher := gosc.NewStandardDispatcher()

	server := &gosc.Server{
		Addr:       fmt.Sprintf("%s:%d", addr, port),
		Dispatcher: dispatcher,
	}

	return &Server{
		server:     server,
		dispatcher: dispatcher,
		addr:       addr,
		port:       port,
	}
}

// AddHandler adds a message handler for a specific OSC address pattern
func (s *Server) AddHandler(pattern string, handler gosc.HandlerFunc) {
	err := s.dispatcher.AddMsgHandler(pattern, handler)
	if err != nil {
		log.Printf("Error adding handler for pattern %s: %v", pattern, err)
	}
}

// Start starts the OSC server
func (s *Server) Start() error {
	log.Printf("Starting OSC server on %s:%d", s.addr, s.port)
	return s.server.ListenAndServe()
}

// Stop stops the OSC server
func (s *Server) Stop() {
	log.Println("Stopping OSC server")
	if err := s.server.CloseConnection(); err != nil {
		log.Printf("Error closing server connection: %v", err)
	}
}
