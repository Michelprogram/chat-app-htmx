package room

import (
	"bytes"
	"context"
	"github.com/gofiber/contrib/websocket"
	"github.com/michelprogram/htmx-go/web/htmx/account"
	"log"
	"sync"
)

const (
	Connect    Event = "connect"
	Disconnect Event = "disconnect"
	Message    Event = "message"
	Remove     Event = "remove"
	Self       Event = "self"
	Emote      Event = "emote"
	Http       Event = "http"
)

var (
	chatrooms *Chatroom = nil
	once      sync.Once
)

func GetInstance() *Chatroom {

	once.Do(func() {
		chatrooms = newChatroom()
	})
	return chatrooms

}

func (c *Chatroom) HandleEvent() {

	for {
		select {
		case event := <-c.Events:
			switch event.Event {
			case Connect:
				c.Connection(event)
			case Message:
				_ = c.Broadcast(event)
			case Disconnect:
				c.RemoveClient(event)
			case Emote:
				_ = c.Emote(event)
			case Self:
				c.Self(event)
			case Http:
				c.BroadcastAll(event)
			case Remove:
				_ = c.RemoveMessage(event)
			}

		}
	}

}

func (c *Chatroom) BroadcastAll(event SenderEvent) {
	for _, client := range c.Clients {
		err := client.WriteMessage(websocket.TextMessage, event.Data)
		if err != nil {
			log.Println("Can't broadcast message: ", err)
		}
	}
}

func (c *Chatroom) Self(event SenderEvent) {
	err := event.Conn.WriteMessage(websocket.TextMessage, event.Data)
	if err != nil {
		log.Println("Can't send message: ", err)
	}
}

func (c *Chatroom) Connection(event SenderEvent) {

	var buffer bytes.Buffer

	c.Clients[event.Client.User.ID.Hex()] = event.Conn

	_ = account.Connect(*event.Client.User).Render(context.TODO(), &buffer)

	event.Data = buffer.Bytes()

	_ = c.Broadcast(event)
}

func (c *Chatroom) RemoveClient(event SenderEvent) {
	delete(c.Clients, event.Client.User.ID.Hex())
}

func (c *Chatroom) Broadcast(event SenderEvent) error {

	for key, client := range c.Clients {
		if key != event.User.ID.Hex() {
			err := client.WriteMessage(websocket.TextMessage, event.Data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Chatroom) Emote(event SenderEvent) error {
	return nil
}

func (c *Chatroom) RemoveMessage(event SenderEvent) error {

	for _, client := range c.Clients {
		err := client.WriteMessage(websocket.TextMessage, event.Data)
		if err != nil {
			return err
		}
	}

	return nil
}
