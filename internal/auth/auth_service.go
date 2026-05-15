package auth

import (
	"fmt"

	"time"

	"github.com/Mikhail-tal63/viltrum_empier/config"
	utils "github.com/Mikhail-tal63/viltrum_empier/utils/JWT"
	"github.com/Mikhail-tal63/viltrum_empier/utils/password"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	repo *AuthRepository
}

func NewAuthService(repo *AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(user *CreateUserPayload) (*AuthResponse, error) {
	now := time.Now()

	existing, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("user with %s already exists", user.Email)
	}
	hash, err := password.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	newUser := &User{
		ID:           primitive.NewObjectID(),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: hash,
		Avatar:       "",
		IsOnline:     false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	createduser, err := s.repo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	secret := []byte(config.ENVs.JWTSecret)
	token, err := utils.CreatJWT(secret, newUser.ID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateRefreshToken(secret, newUser.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		User:         *createduser,
	}, nil
}

func (s *AuthService) Login(email, HashPassword string) (*AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	if !password.ComparePassword(user.PasswordHash, []byte(HashPassword)) {
		return nil, fmt.Errorf("invalid password")
	}

	secret := []byte(config.ENVs.JWTSecret)
	token, err := utils.CreatJWT(secret, user.ID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateRefreshToken(secret, user.ID)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil

}
