package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client: ", ws.RemoteAddr())

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) handleWSOrderBook(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client to orderbook feed: ", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 2)
	}
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buff := make([]byte, 1024)

	for {
		//looping through all connections
		n, err := ws.Read(buff)
		if err != nil {
			if err == io.EOF {
				break //connection has been closed
			}
			fmt.Println("error: ", err)
			continue //not to use break or return else it will break the connection
		}
		//displaying bytes we read through
		msg := buff[:n]
		fmt.Println(string(msg))

		//letting everybody connected know there is a message
		s.broaddcast(msg)

		//write back to connected conn
		ws.Write([]byte("thank you for the message"))
	}
}

func (s *Server) broaddcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error: ", err)
			}
		}(ws)
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/orderws", websocket.Handler(server.handleWSOrderBook))
	http.ListenAndServe(":8080", nil)
}
