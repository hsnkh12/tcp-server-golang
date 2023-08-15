package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Address        string
	Listener       net.Listener
	QuitChannel    chan struct{}
	MessageChannel chan Message
}

func CreateTCPServer(addr string) *Server {

	return &Server{
		Address:        addr,
		QuitChannel:    make(chan struct{}),
		MessageChannel: make(chan Message, 10),
	}
}

func (s *Server) Listen() error {

	listener, err := net.Listen("tcp", s.Address)

	if err != nil {
		return err
	}

	defer listener.Close()
	s.Listener = listener

	go s.AcceptConnections()

	fmt.Println("Server is listening on ", s.Address)

	<-s.QuitChannel
	close(s.MessageChannel)
	return nil

}

func (s *Server) AcceptConnections() {

	for {

		conn, err := s.Listener.Accept()

		if err != nil {
			log.Fatalf("error accepting connection: %v", err)
		}

		log.Println("new connection: ", conn.RemoteAddr())

		go s.ReadConneciton(conn)
	}
}

func (s *Server) ReadConneciton(conn net.Conn) {

	defer conn.Close()
	buffer := make([]byte, 2050)

	for {

		bytesRead, err := conn.Read(buffer)

		if err != nil {
			log.Fatalf("error reading to buffer: %v", err)
			continue
		}

		header := HeaderMessage{
			FromAddress: conn.RemoteAddr().String(),
		}

		s.MessageChannel <- Message{
			Header:  header,
			Payload: buffer[:bytesRead],
		}

		conn.Write([]byte("< Message received >"))
	}
}
