package notification

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`

	Type string `bson:"type" json:"type"`

	Content string `bson:"content" json:"content"`

	IsRead bool `bson:"is_read" json:"is_read"`

	RelatedID primitive.ObjectID `bson:"related_id,omitempty" json:"related_id,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
