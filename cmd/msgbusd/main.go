package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/julienschmidt/httprouter"
	"github.com/prologic/msgbus"
)

// IndexHandler ...
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// PushHandler ...
func PushHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	topic := p.ByName("topic")
	websocket.Handler(PushWebSocketHandler(topic)).ServeHTTP(w, r)
}

// PullHandler ...
func PullHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	topic := p.ByName("topic")
	message, ok := mb.Get(topic)
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	out, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.Write(out)
}

// PutHandler ...
func PutHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	topic := p.ByName("topic")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	mb.Put(topic, mb.NewMessage(body))
}

// PushWebSocketHandler ...
func PushWebSocketHandler(topic string) websocket.Handler {
	return func(conn *websocket.Conn) {
		id := conn.Request().RemoteAddr
		ch := mb.Subscribe(id, topic)
		defer func() {
			mb.Unsubscribe(id, topic)
		}()

		var (
			err error
			ack msgbus.Ack
		)

		for {
			msg := <-ch
			err = websocket.JSON.Send(conn, msg)
			if err != nil {
				log.Printf("Error sending msg to %s", id)
				continue
			}

			err = websocket.JSON.Receive(conn, &ack)
			if err != nil {
				log.Printf("Error receiving ack from %s", id)
				continue
			}

			log.Printf("message %v acked %v by %s", msg, ack, id)
		}
	}
}

var mb = msgbus.NewMessageBus()

func main() {
	router := httprouter.New()

	// UI
	router.GET("/", IndexHandler)

	// Subscribe
	router.GET("/push/:topic", PushHandler)
	router.GET("/pull/:topic", PullHandler)

	// Publish
	router.PUT("/:topic", PutHandler)

	log.Fatal(http.ListenAndServe(":8000", router))
}
