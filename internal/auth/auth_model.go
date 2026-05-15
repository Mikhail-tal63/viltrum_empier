package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Avatar       string             `bson:"avatar,omitempty" json:"avatar,omitempty"`

	IsOnline bool `bson:"is_online" json:"is_online"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type CreateUserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}
