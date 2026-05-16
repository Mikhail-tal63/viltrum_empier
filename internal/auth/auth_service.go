package auth

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"time"

	"github.com/Mikhail-tal63/viltrum_empier/config"
	utils "github.com/Mikhail-tal63/viltrum_empier/utils/JWT"
	"github.com/Mikhail-tal63/viltrum_empier/utils/password"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var usernameRe = regexp.MustCompile(`^[a-z0-9_]{3,20}$`)

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

	username := strings.ToLower(strings.TrimSpace(user.Username))
	if !usernameRe.MatchString(username) {
		return nil, fmt.Errorf("username must be 3-20 chars, lowercase letters, numbers or underscore")
	}
	name := strings.TrimSpace(user.Name)
	if name == "" {
		name = username
	}

	existing, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("user with %s already exists", user.Email)
	}

	existingByUsername, err := s.repo.GetUserByUsername(context.TODO(), username)
	if err != nil {
		return nil, err
	}
	if existingByUsername != nil {
		return nil, fmt.Errorf("username %q is already taken", username)
	}

	hash, err := password.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	newUser := &User{
		ID:           primitive.NewObjectID(),
		Name:         name,
		Username:     username,
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

func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}
