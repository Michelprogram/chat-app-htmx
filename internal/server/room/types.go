package room

import (
	"encoding/json"
	"github.com/gofiber/contrib/websocket"
	"github.com/michelprogram/htmx-go/internal/models"
)

type Data struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	HEADERS struct {
		HXRequest     string      `json:"HX-Request"`
		HXTrigger     string      `json:"HX-Trigger"`
		HXTriggerName interface{} `json:"HX-Trigger-Name"`
		HXTarget      string      `json:"HX-Target"`
		HXCurrentURL  string      `json:"HX-Current-URL"`
	} `json:"HEADERS"`
}

type Event string

type Client struct {
	User *models.User
	Conn *websocket.Conn
}

// TODO: Maybe turn []byte to string
type SenderEvent struct {
	Event
	Client
	Data []byte
}

func NewSenderEvent(event Event, client Client, data []byte) SenderEvent {
	return SenderEvent{
		Event:  event,
		Client: client,
		Data:   data,
	}
}

type Chatroom struct {
	Clients map[string]*websocket.Conn
	Events  chan SenderEvent
}

func (c *Chatroom) GetUsers() []string {
	usernames := make([]string, 0, len(c.Clients))

	for k := range c.Clients {
		usernames = append(usernames, k)
	}

	return usernames
}

func newChatroom() *Chatroom {
	return &Chatroom{
		Clients: make(map[string]*websocket.Conn),
		Events:  make(chan SenderEvent),
	}
}

func ParseData(msg []byte) (*Data, error) {

	var data Data

	err := json.Unmarshal(msg, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}
