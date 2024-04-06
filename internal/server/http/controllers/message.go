package controllers

import (
	"bytes"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/michelprogram/htmx-go/internal/repository"
	"github.com/michelprogram/htmx-go/internal/server/room"
	"github.com/michelprogram/htmx-go/web/htmx/chat"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageController struct {
	repository.IMessageRepository
	repository.IUserRepository
}

func NewMessageController(messageRepository repository.IMessageRepository, userRepository repository.IUserRepository) *MessageController {
	return &MessageController{
		IMessageRepository: messageRepository,
		IUserRepository:    userRepository,
	}
}

func (mc MessageController) LastMessages(c *fiber.Ctx) error {

	messageID := c.Params("id")

	objId, err := primitive.ObjectIDFromHex(messageID)

	if err != nil {
		return err
	}

	messages, err := mc.GetLast(10, objId)

	if err != nil {
		return err
	}

	var buffer bytes.Buffer

	_ = chat.Pagination(c.Locals("id").(string), messages).Render(context.TODO(), &buffer)

	return c.Send(buffer.Bytes())

}

func (mc MessageController) Remove(c *fiber.Ctx) error {

	var buffer bytes.Buffer

	userID := c.Locals("id").(string)

	messageID := c.Params("id")

	message, err := mc.Delete(messageID, userID)

	if err != nil || message == nil {
		return c.SendStatus(fiber.StatusForbidden)
	}

	user := mc.FindById(userID)

	dest := make([]byte, len(buffer.Bytes()))

	copy(dest, buffer.Bytes())

	err = chat.EmptyMessage(*user, *message, false).Render(context.TODO(), &buffer)

	if err != nil {
		return err
	}

	room.GetInstance().Events <- room.NewSenderEvent(room.Remove, room.Client{}, buffer.Bytes())

	buffer.Reset()

	err = chat.EmptyMessage(*user, *message, true).Render(context.TODO(), &buffer)

	return c.Send(buffer.Bytes())

}
