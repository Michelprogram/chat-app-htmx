package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michelprogram/htmx-go/internal/models"
	"github.com/michelprogram/htmx-go/internal/pkg"
	"github.com/michelprogram/htmx-go/internal/repository"
	"github.com/michelprogram/htmx-go/internal/server/room"
	htmx "github.com/michelprogram/htmx-go/web/htmx/chat/emote"
)

type EmoteController struct {
	Emote repository.IEmoteRepository
	User  repository.IUserRepository
}

func NewEmoteController(emoteRepository repository.IEmoteRepository, userRepository repository.IUserRepository) *EmoteController {
	return &EmoteController{
		Emote: emoteRepository,
		User:  userRepository,
	}
}

func (ec EmoteController) Add(c *fiber.Ctx) error {

	messageID := c.Params("messageid")

	userID := c.Locals("id").(string)

	emoji, err := pkg.GetRandomEmoji()

	if err != nil {
		return err
	}

	emote := models.NewEmote(emoji.DownloadURL)

	err = ec.Emote.Insert(messageID, emote)

	if err != nil {
		return err
	}

	buffer, err := pkg.ParseTemplToBuffer(c, htmx.Insert(messageID, *emote))

	if err != nil {
		return err
	}

	room.GetInstance().Events <- room.SenderEvent{Event: room.Message, Data: buffer, Client: room.Client{
		User: ec.User.FindById(userID),
		Conn: nil,
	}}

	return c.Send(buffer)
}

func (ec EmoteController) Counter(c *fiber.Ctx) error {

	emoteID := c.Params("emoteid")

	userID := c.Locals("id").(string)

	emote, err := ec.Emote.AddCounter(emoteID)

	if err != nil {
		return err
	}

	buffer, err := pkg.ParseTemplToBuffer(c, htmx.Emote(emote))

	if err != nil {
		return err
	}

	room.GetInstance().Events <- room.SenderEvent{Event: room.Message, Data: buffer, Client: room.Client{
		User: ec.User.FindById(userID),
		Conn: nil,
	}}

	return c.Send(buffer)
}
