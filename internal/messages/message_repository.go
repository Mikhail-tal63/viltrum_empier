package messages

import "go.mongodb.org/mongo-driver/mongo"

type Messagerepo struct {
	msgRepo   *mongo.Collection
	groupRepo *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) *Messagerepo {
	return &Messagerepo{
		msgRepo:   db.Collection("message"),
		groupRepo: db.Collection("group_chat"),
	}
}
