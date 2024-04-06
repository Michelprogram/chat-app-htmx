package repository

import (
	"github.com/michelprogram/htmx-go/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUserRepository interface {
	Insert(username, password string) (*models.User, error)
	Find(username, password string) *models.User
	FindById(id string) *models.User
	FindByUsername(username string) *models.User
	FindAndUpdate(id string, user models.User) error
	UpdateToken(id primitive.ObjectID, token string) error
	UsersnamesByIDs(ids []string) ([]models.User, error)
}

type IEmoteRepository interface {
	Insert(messageid string, emote *models.Emote) error
	AddCounter(emoteid string) (models.Emote, error)
}

type IMessageRepository interface {
	Insert(content string, author primitive.ObjectID) (*models.Message, error)
	Delete(id, userID string) (*models.Message, error)
	GetLast(n int64, from primitive.ObjectID) ([]models.EmbdedMessage, error)
	GetLastInsert() (models.Message, error)
	GetLasts(n int64) ([]models.EmbdedMessage, error)
}
