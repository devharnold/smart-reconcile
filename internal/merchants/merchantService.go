// Change every user to use merchant

package merchants

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/devharnold/smart-reconcile/internal/auth"
)

type MerchantService struct {
	Repo       *MerchantRepository
	jwtService auth.JWTService
}

func NewMerchantService(repo *MerchantRepository, jwtService auth.JWTService) *MerchantService {
	return &MerchantService{
		Repo:       repo,
		jwtService: jwtService,
	}
}

func (s *MerchantService) RegisterMerchant(ctx context.Context, businessName, Email, phoneNumber, password string) (*Merchants, error) {

	if businessName == "" || Email == "" || password == "" {
		return nil, errors.New("all fields are required")
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	merchant := &Merchants{
		BusinessName: businessName,
		Email:        Email,
		PhoneNumber:  phoneNumber,
		Password:     hashedPassword,
	}

	id, err := s.Repo.RegisterMerchant(ctx, merchant)
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}

	merchant.UserID = id

	token, err := s.jwtService.GenerateToken(
		fmt.Sprint(merchant.UserID),
		merchant.Email,
		"user",
	)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	merchant.Token = token
	return merchant, nil
}

func (s *MerchantService) MerchantLogin(ctx context.Context, email, password string) (*Merchants, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password required")
	}

	merchant, err := s.Repo.GetMerchantByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !auth.VerifyPassword(password, merchant.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(
		fmt.Sprint(merchant.UserID),
		merchant.Email,
		"user",
	)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	merchant.Token = token

	log.Printf("login successful user_id=%d", merchant.UserID)
	return merchant, nil
}

func (s *MerchantService) GetUserByEmail(ctx context.Context, email string) (*Merchants, error) {
	return s.Repo.GetMerchantByEmail(ctx, email)
}
