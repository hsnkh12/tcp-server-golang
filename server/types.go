package server

type HeaderMessage struct {
	FromAddress string
}
type Message struct {
	Header  HeaderMessage
	Payload []byte
}
