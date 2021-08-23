package service

import (
	"context"
	"errors"
	"github.com/Roman-Shine/value_backend/internal/domain"
	"github.com/Roman-Shine/value_backend/internal/repository"
	"github.com/Roman-Shine/value_backend/pkg/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UsersService struct {
	repo   repository.Users
	hasher hash.PasswordHasher

	emailService Emails

	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
	verificationCodeLength int
	domain                 string
}

func NewUsersService(repo repository.Users, hasher hash.PasswordHasher, accessTTL, refreshTTL time.Duration, verificationCodeLength int, domain string) *UsersService {
	return &UsersService{
		repo:                   repo,
		hasher:                 hasher,
		accessTokenTTL:         accessTTL,
		refreshTokenTTL:        refreshTTL,
		verificationCodeLength: verificationCodeLength,
		domain:                 domain,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input UserSignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	//verificationCode := s.otpGenerator.RandomSecret(s.verificationCodeLength)

	user := domain.User{
		Name:         input.Name,
		Password:     passwordHash,
		Phone:        input.Phone,
		Email:        input.Email,
		RegisteredAt: time.Now(),
		LastVisitAt:  time.Now(),
		//Verification: domain.Verification{
		//	Code: verificationCode,
		//},
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return err
		}

		return err
	}

	// todo. DECIDE ON EMAIL MARKETING STRATEGY
	return s.emailService.SendUserVerificationEmail(VerificationEmailInput{
		Email: user.Email,
		Name:  user.Name,
		//VerificationCode: verificationCode,
	})
}

func (s *UsersService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	//passwordHash, err := s.hasher.Hash(input.Password)
	//if err != nil {
	//	return Tokens{}, err
	//}

	//user, err := s.repo.GetByCredentials(ctx, input.Email, passwordHash)
	user, err := s.repo.GetByCredentials(ctx, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return Tokens{}, err
		}

		return Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	student, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, student.ID)
}

func (s *UsersService) Verify(ctx context.Context, userId primitive.ObjectID, hash string) error {
	err := s.repo.Verify(ctx, userId, hash)
	if err != nil {
		if errors.Is(err, domain.ErrVerificationCodeInvalid) {
			return err
		}

		return err
	}

	return nil
}

func (s *UsersService) createSession(ctx context.Context, userId primitive.ObjectID) (Tokens, error) {
	var res Tokens
	res.AccessToken = string("123")
	res.RefreshToken = string("123")
	return res, nil
}

//func (s *UsersService) createSession(ctx context.Context, userId primitive.ObjectID) (Tokens, error) {
//	var (
//		res Tokens
//		err error
//	)
//
//	res.AccessToken, err = s.tokenManager.NewJWT(userId.Hex(), s.accessTokenTTL)
//	if err != nil {
//		return res, err
//	}
//
//	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
//	if err != nil {
//		return res, err
//	}
//
//	session := domain.Session{
//		RefreshToken: res.RefreshToken,
//		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
//	}
//
//	err = s.repo.SetSession(ctx, userId, session)
//
//	return res, err
//}
