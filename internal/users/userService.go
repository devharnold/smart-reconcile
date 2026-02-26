package users

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/devharnold/smart-reconcile/internal/auth"
)

type UserService struct {
	Repo       *UsersRepository
	jwtService auth.JWTService
}

func NewUserService(repo *UsersRepository, jwtService auth.JWTService) *UserService {
	return &UserService{
		Repo:       repo,
		jwtService: jwtService,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, firstName, lastName, userEmail, phoneNumber, password string) (*Users, error) {

	if firstName == "" || lastName == "" || userEmail == "" || password == "" {
		return nil, errors.New("all fields are required")
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &Users{
		FirstName:   firstName,
		LastName:    lastName,
		Email:       userEmail,
		PhoneNumber: phoneNumber,
		Password:    hashedPassword,
	}

	id, err := s.Repo.RegisterUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}

	user.UserID = id

	token, err := s.jwtService.GenerateToken(
		fmt.Sprint(user.UserID),
		user.Email,
		"user",
	)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	user.Token = token
	return user, nil
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (*Users, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password required")
	}

	user, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !auth.VerifyPassword(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(
		fmt.Sprint(user.UserID),
		user.Email,
		"user",
	)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	user.Token = token

	log.Printf("login successful user_id=%d", user.UserID)
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*Users, error) {
	return s.Repo.GetUserByEmail(ctx, email)
}
