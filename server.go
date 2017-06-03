package msgbus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/julienschmidt/httprouter"
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

type Server struct {
	bus    *msgbus.MsgBus
	router *httprouter.Router
}

func (s *Server) init() {
	s.bus = msgbus.NewMessageBus()
	s.router = httprouter.New()

	router.GET("/", IndexHandler)
	router.GET("/push/:topic", PushHandler)
	router.GET("/pull/:topic", PullHandler)
	router.PUT("/:topic", PutHandler)
}

func (s *Server) Run() {
	s.init()

	log.Fatal(http.ListenAndServe(s.bind, s.router))
}
