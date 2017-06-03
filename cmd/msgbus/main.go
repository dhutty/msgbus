package main

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"

	"github.com/prologic/msgbus"
)

func main() {
	origin := "http://localhost/"
	url := "ws://localhost:8000/push/foo"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	var (
		err error
		msg msgbus.Message
	)

	for {
		err = websocket.JSON.Receive(ws, &msg)
		if err != nil {
			log.Fatal(err)
		}

		ack := Ack{Ack: msg.Id}
		err = websocket.JSON.Send(ws, &ack)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Received: %s.\n", msg.Payload)
	}
}
