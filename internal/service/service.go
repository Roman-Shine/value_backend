package service

import (
	"context"
	"github.com/Roman-Shine/value_backend/internal/repository"
	"github.com/Roman-Shine/value_backend/pkg/cache"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserSignUpInput struct {
	Name     string
	Email    string
	Phone    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	Verify(ctx context.Context, userId primitive.ObjectID, hash string) error
}

type Services struct {
	Users Users
}

type Deps struct {
	Repos                  *repository.Repositories
	Cache                  cache.Cache
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	PaymentCallbackURL     string
	PaymentResponseURL     string
	CacheTTL               int64
	VerificationCodeLength int
	Environment            string
	Domain                 string
}

func NewServices(deps Deps) *Services {
	usersService := NewUsersService(deps.Repos.Users, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.VerificationCodeLength, deps.Domain)

	return &Services{
		Users: usersService,
	}
}
