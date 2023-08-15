package main

import (
	"fmt"
	"tcp_server/server"
)

func main() {

	server := server.CreateTCPServer(":3000")

	go func() {
		for msg := range server.MessageChannel {
			fmt.Printf("< Message \n < Headers: address: %s > \n < Payload: %s > \n >", msg.Header.FromAddress, msg.Payload)
		}
	}()
	server.Listen()
}
