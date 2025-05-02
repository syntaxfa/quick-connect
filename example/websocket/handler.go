package main

import "github.com/gorilla/websocket"

func echoHandler(conn *websocket.Conn) {
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(mt, message)
	}
}
