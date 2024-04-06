package sockets

import (
	"bytes"
	"context"
	"errors"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/michelprogram/htmx-go/internal/repository"
	"github.com/michelprogram/htmx-go/internal/server/room"
	"github.com/michelprogram/htmx-go/web/htmx/account"
	"github.com/michelprogram/htmx-go/web/htmx/chat"
	"log"
)

type Chat struct {
	UserRepository    repository.IUserRepository
	MessageRepository repository.IMessageRepository
}

func NewChat(userRepository repository.IUserRepository, messageRepository repository.IMessageRepository) Chat {
	return Chat{
		userRepository,
		messageRepository,
	}
}

func (ct Chat) read(client room.Client, msg []byte) error {

	var buffer bytes.Buffer

	//Parse message
	data, err := room.ParseData(msg)

	if err != nil {
		return errors.New("Error while reading message: " + err.Error())
	}

	//Save message in database
	message, err := ct.MessageRepository.Insert(data.Message, client.User.ID)

	if err != nil {
		return errors.New("Error while inserting message: " + err.Error())
	}

	//Set message in chatroom
	err = chat.Message(*client.User, *message, false).Render(context.TODO(), &buffer)

	if err != nil {
		return errors.New("Error while rendering message: " + err.Error())
	}

	dest := make([]byte, len(buffer.Bytes()))

	copy(dest, buffer.Bytes())

	//Send message to all clients
	room.GetInstance().Events <- room.SenderEvent{Event: room.Message, Client: client, Data: dest}

	buffer.Reset()

	//Set message in your account section
	err = chat.Message(*client.User, *message, true).Render(context.TODO(), &buffer)

	if err != nil {
		return errors.New("Error while rendering message: " + err.Error())
	}

	err = client.Conn.WriteMessage(websocket.TextMessage, buffer.Bytes())

	if err != nil {
		return errors.New("Error while sending message: " + err.Error())
	}

	return nil
}

func (ct Chat) connect(userID string, client room.Client) ([]byte, error) {

	var buffer bytes.Buffer

	//Get BUCKET_MESSAGE_SIZE last messages
	messages, err := ct.MessageRepository.GetLasts(repository.BUCKET_MESSAGE_SIZE)

	if err != nil {
		return nil, err
	}

	//Set last messages
	err = chat.LastMessage(userID, messages).Render(context.TODO(), &buffer)

	if err != nil {
		return nil, err
	}

	//Add client to chatroom, send message to all clients
	room.GetInstance().Events <- room.SenderEvent{Event: room.Connect, Client: client}

	//Set username in your account section
	err = account.Index(client.User.Username).Render(context.TODO(), &buffer)

	if err != nil {
		return nil, err
	}

	//Get users already connected
	users, err := ct.UserRepository.UsersnamesByIDs(room.GetInstance().GetUsers())

	if err != nil {
		return nil, err
	}

	err = account.AlreadyConnected(users).Render(context.TODO(), &buffer)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (ct Chat) disconnect(client room.Client) {
	room.GetInstance().Events <- room.SenderEvent{Event: room.Disconnect, Client: client}
}

func (ct Chat) Chats() fiber.Handler {

	return websocket.New(func(c *websocket.Conn) {

		var (
			msg []byte
			err error
		)

		userID := c.Locals("id").(string)

		client := room.Client{
			User: ct.UserRepository.FindById(userID),
			Conn: c,
		}

		data, err := ct.connect(userID, client)

		if err != nil {
			_ = c.Close()
		}

		_ = c.Conn.WriteMessage(websocket.TextMessage, data)

		for {

			if _, msg, err = c.ReadMessage(); err != nil {
				ct.disconnect(client)
				break
			} else {

				err = ct.read(client, msg)

				if err != nil {
					log.Println(err)
				}
			}
		}
	})

}
