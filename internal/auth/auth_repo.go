package auth

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	collection *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) *AuthRepository {
	return &AuthRepository{
		collection: db.Collection("users"),
	}
}

func (r *AuthRepository) CreateUser(user *User) (*User, error) {
	res, err := r.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	user.ID = insertedID

	return user, nil
}

func (r *AuthRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	filter := bson.M{"email": email}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	filter := bson.M{"username": username}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
