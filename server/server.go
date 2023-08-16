package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Server struct {
	Address           string
	Listener          net.Listener
	QuitChannel       chan struct{}         // Channel for signaling server shutdown
	ReceiveBuffer     chan Message          // Channel for receiving messages from connections
	SendBuffer        chan string           // Channel for sending responses to connections
	Wg                sync.WaitGroup        // WaitGroup to manage goroutine lifecycles
	ActiveConnections map[net.Conn]struct{} // Map to track active connections
	ActiveConnsMux    sync.Mutex            // Mutex for concurrent map access
}

func CreateTCPServer(addr string) *Server {
	return &Server{
		Address:           addr,
		QuitChannel:       make(chan struct{}),
		ReceiveBuffer:     make(chan Message, 10),
		SendBuffer:        make(chan string, 10),
		ActiveConnections: make(map[net.Conn]struct{}),
	}
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	defer listener.Close()
	s.Listener = listener

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		fmt.Printf("Received signal %s, shutting down...\n", sig)

		// Close the QuitChannel to signal graceful shutdown
		close(s.QuitChannel)

		// Close all active connections before shutting down
		s.CloseAllConnections()
	}()

	// Start accepting connections in a separate goroutine
	s.Wg.Add(1)
	go s.AcceptConnections()

	fmt.Println("Server is listening on", s.Address)

	// Wait for QuitChannel to be closed (server shutdown)
	<-s.QuitChannel

	// Close the ReceiveBuffer and SendBuffer channels
	close(s.ReceiveBuffer)
	close(s.SendBuffer)

	return nil
}

func (s *Server) AcceptConnections() {
	defer s.Wg.Done()

	for {
		select {
		case <-s.QuitChannel:
			return
		default:
		}

		conn, err := s.Listener.Accept()

		if err != nil {
			log.Printf("error accepting connection: %v", err)
			continue
		}

		// Add the connection to the active connections map
		s.ActiveConnsMux.Lock()
		s.ActiveConnections[conn] = struct{}{}
		s.ActiveConnsMux.Unlock()

		log.Println("new connection:", conn.RemoteAddr())

		// Start handling the connection in a separate goroutine
		s.Wg.Add(1)
		go s.ReadConneciton(conn)
	}
}

func (s *Server) ReadConneciton(conn net.Conn) {
	defer conn.Close()
	defer s.Wg.Done()

	buffer := make([]byte, 2050)

	for {
		bytesRead, err := conn.Read(buffer)

		if err != nil {
			log.Printf("Connection closed: %s", conn.RemoteAddr().Network())

			// Remove the connection from the active connections map
			s.ActiveConnsMux.Lock()
			delete(s.ActiveConnections, conn)
			s.ActiveConnsMux.Unlock()

			return
		}

		header := HeaderMessage{
			FromAddress: conn.RemoteAddr().String(),
		}

		s.ReceiveBuffer <- Message{
			Header:  header,
			Payload: buffer[:bytesRead],
		}

		respond := <-s.SendBuffer

		conn.Write([]byte(respond))
	}
}

func (s *Server) CloseAllConnections() {
	s.ActiveConnsMux.Lock()
	defer s.ActiveConnsMux.Unlock()

	for conn := range s.ActiveConnections {
		conn.Close()
		// You can access additional connection information here
		delete(s.ActiveConnections, conn)
	}
}
