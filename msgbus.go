package msgbus

import (
	"log"
	"time"
)

// Message ...
type Message struct {
	ID      uint64    `json:"id"`
	Payload []byte    `json:"payload"`
	Created time.Time `json:"created"`
	Acked   time.Time `json:"acked"`
}

// Ack ...
type Ack struct {
	Ack uint64 `json:"ack"`
}

// Listeners ...
type Listeners struct {
	ids map[string]bool
	chs map[string]chan *Message
}

// NewListeners ...
func NewListeners() *Listeners {
	return &Listeners{
		ids: make(map[string]bool),
		chs: make(map[string]chan *Message),
	}
}

// Add ...
func (ls *Listeners) Add(id string) chan *Message {
	log.Printf("Listeners.Add(%s)\n", id)
	ls.ids[id] = true
	ls.chs[id] = make(chan *Message)
	return ls.chs[id]
}

// Remove ...
func (ls *Listeners) Remove(id string) {
	log.Printf("Listeners.Remove(%s)\n", id)
	delete(ls.ids, id)

	close(ls.chs[id])
	delete(ls.chs, id)
}

// Exists ...
func (ls *Listeners) Exists(id string) bool {
	_, ok := ls.ids[id]
	return ok
}

// Get ...
func (ls *Listeners) Get(id string) (chan *Message, bool) {
	ch, ok := ls.chs[id]
	if !ok {
		return nil, false
	}
	return ch, true
}

// NotifyAll ...
func (ls *Listeners) NotifyAll(message *Message) {
	log.Printf("Listeners.NotifyAll(%v)\n", message)
	for _, ch := range ls.chs {
		ch <- message
	}
}

// MessageBus ...
type MessageBus struct {
	seqid     uint64
	topics    map[string]*Queue
	listeners map[string]*Listeners
}

// NewMessageBus ...
func NewMessageBus() *MessageBus {
	return &MessageBus{
		topics:    make(map[string]*Queue),
		listeners: make(map[string]*Listeners),
	}
}

// Len ...
func (mb *MessageBus) Len() int {
	return len(mb.topics)
}

// NewMessage ...
func (mb *MessageBus) NewMessage(payload []byte) *Message {
	message := &Message{
		ID:      mb.seqid,
		Payload: payload,
		Created: time.Now(),
	}

	mb.seqid++

	return message
}

// Put ...
func (mb *MessageBus) Put(topic string, message *Message) {
	log.Printf("MessageBus.Put(%s, %v)\n", topic, message)
	q, ok := mb.topics[topic]
	if !ok {
		q = &Queue{}
		mb.topics[topic] = q
	}
	q.Push(message)

	mb.NotifyAll(topic, message)
}

// Get ...
func (mb *MessageBus) Get(topic string) (*Message, bool) {
	log.Printf("MessageBus.Get(%s)\n", topic)
	q, ok := mb.topics[topic]
	if !ok {
		return &Message{}, false
	}

	m := q.Pop()
	if m == nil {
		return &Message{}, false
	}
	return m.(*Message), true
}

// NotifyAll ...
func (mb *MessageBus) NotifyAll(topic string, message *Message) {
	log.Printf("MessageBus.NotifyAll(%s, %v)\n", topic, message)
	ls, ok := mb.listeners[topic]
	if !ok {
		return
	}
	ls.NotifyAll(message)
}

// Subscribe ...
func (mb *MessageBus) Subscribe(id, topic string) chan *Message {
	log.Printf("MessageBus.Subscribe(%s, %s)\n", id, topic)
	ls, ok := mb.listeners[topic]
	if !ok {
		ls = NewListeners()
		mb.listeners[topic] = ls
	}

	if ls.Exists(id) {
		// Already verified th listener exists
		ch, _ := ls.Get(id)
		return ch
	}
	return ls.Add(id)
}

// Unsubscribe ...
func (mb *MessageBus) Unsubscribe(id, topic string) {
	log.Printf("MessageBus.Unsubscribe(%s, %s)\n", id, topic)
	ls, ok := mb.listeners[topic]
	if !ok {
		return
	}

	if ls.Exists(id) {
		// Already verified th listener exists
		ls.Remove(id)
	}
}
