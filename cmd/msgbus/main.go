package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/net/websocket"

	"github.com/prologic/msgbus"
)

var (
	host string
	port int
	err  error
	msg  msgbus.Message
	ws   *websocket.Conn
)

func init() {
	flag.StringVar(&host, "host", "localhost", "host to connect to")
	flag.IntVar(&port, "port", 8000, "port to connect to")
}

func main() {
	flag.Parse()

	origin := "http://localhost/"
	url := fmt.Sprintf("ws://%s:%d/push/foo", host, port)
	ws, err = websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening for messages from %s", url)

	for {
		err = websocket.JSON.Receive(ws, &msg)
		if err != nil {
			log.Fatal(err)
		}

		ack := msgbus.Ack{Ack: msg.ID}
		err = websocket.JSON.Send(ws, &ack)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Received: %s\n", msg.Payload)
	}
}
