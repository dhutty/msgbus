package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"

	"github.com/prologic/msgbus"
)

const defaultTopic = "hello"

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

	if flag.Arg(0) == "sub" {
		subscribe(flag.Arg(1))
	} else if flag.Arg(0) == "pub" {
		publish(flag.Arg(1), flag.Arg(2))
	} else {
		log.Fatalf("invalid command %s", flag.Arg(0))
	}
}

func publish(topic string, message string) {
	var payload bytes.Buffer

	if topic == "" {
		topic = defaultTopic
	}

	if message == "" || message == "-" {
		log.Printf("Reading message from stdin...\n")
		buf, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		payload.Write(buf)
	} else {
		payload.Write([]byte(message))
	}

	url := fmt.Sprintf("http://%s:%d/%s", host, port, topic)

	client := &http.Client{}

	req, err := http.NewRequest("PUT", url, &payload)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}

func subscribe(topic string) {
	if topic == "" {
		topic = defaultTopic
	}

	origin := "http://localhost/"
	url := fmt.Sprintf("ws://%s:%d/push/%s", host, port, topic)
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
