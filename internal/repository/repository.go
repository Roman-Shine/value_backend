package repository

import (
	"context"
	"github.com/Roman-Shine/value_backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	Verify(ctx context.Context, userId primitive.ObjectID, code string) error
	SetSession(ctx context.Context, userId primitive.ObjectID, session domain.Session) error
}

type Repositories struct {
	Users Users
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
	}
}
