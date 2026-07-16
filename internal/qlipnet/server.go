package qlipnet

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Server struct {
	mu       sync.Mutex
	clients  map[string]net.Conn
	listener net.Listener
}

// Create a new server instance.
func NewServer() *Server {
	return &Server{
		clients:  make(map[string]net.Conn),
		listener: nil,
	}
}

// Serve on a specified port.
func (s *Server) Serve(port int) {
	// Listen on TCP port specified by user
	addrString := fmt.Sprintf(":%v", port)
	l, err := net.Listen("tcp", addrString)
	if err != nil {
		log.Fatal(err)
	}

	s.listener = l

	defer s.listener.Close()
	for {
		// Wait for a connection.
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("shutting down server: %v", err)
			return
		}

		// Handle the connection in a new goroutine.
		go func(conn net.Conn) {
			defer s.handleClose(conn)
			for s.handle(conn) {
			}
		}(conn)
	}
}

// Gracefully shutdown the server.
func (s *Server) Close() {
	s.listener.Close()
	for _, conn := range s.clients {
		conn.Close()
	}
}

// Dispatcher for server. Handles all protocols and interactions
// between the clients and the server.
func (s *Server) handle(c net.Conn) bool {
	packet, err := readPacket(c)
	if err != nil {
		if !errors.Is(err, net.ErrClosed) && err != io.EOF {
			log.Printf("read error: %v", err)
		}
		s.handleClose(c)
		return false
	}

	switch packet.Type {
	case InitMsgType:
		s.handleInit(c)
	case QlipBoardTextMsgType, QlipBoardImgMsgType:
		s.handleClipboard(c, packet)
	case CloseMsgType:
		return false
	default:
		log.Printf("unknown packet type %d", packet.Type)
	}

	return true
}

// Close a connection in the server client lookup.
func (s *Server) handleClose(c net.Conn) {
	addr := c.RemoteAddr().String()
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, addr)
	c.Close()
}

// Handle an incoming initialization handshake.
func (s *Server) handleInit(c net.Conn) {
	// Get the IP of the incoming address to use
	// as a lookup key for this connection.
	addr := c.RemoteAddr().String()
	p := Packet{
		Type: InitMsgType,
		Len:  0,
		Data: nil,
	}

	// Register the client connection
	s.mu.Lock()
	defer s.mu.Unlock()
	err := writePacket(c, p)
	if err != nil {
		log.Printf("could not register addr %v: %v", addr, err)
		return
	}
	s.clients[addr] = c
}

// Handle an incoming clipboard message.
func (s *Server) handleClipboard(c net.Conn, p Packet) {
	// Get the IP of the incoming address to
	// ensure data is only sent to other hosts.
	sender := c.RemoteAddr().String()

	// Get a list of connections to send to
	s.mu.Lock()
	var conns []net.Conn
	for addr, conn := range s.clients {
		if addr != sender {
			conns = append(conns, conn)
		}
	}
	s.mu.Unlock()

	// Broadcast clipboard data
	for _, conn := range conns {
		err := writePacket(conn, p)
		if err != nil {
			log.Printf("write failed: %v", err)
		}
	}
}
