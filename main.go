package main

import (
	"fmt"
	"tcp_server/server"
)

func main() {

	server := server.CreateTCPServer("localhost:3000")

	go func() {
		for msg := range server.ReceiveBuffer {

			fmt.Printf("< Message \n < Headers: address: %s > \n < Payload: %s > \n >", msg.Header.FromAddress, msg.Payload)
			response := "< Message from " + msg.Header.FromAddress + " Received>"

			select {
			case server.SendBuffer <- response:
			case <-server.QuitChannel:
				return
			}

		}
	}()

	server.Listen()

}
