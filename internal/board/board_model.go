package board

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Board struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	WorkspaceID primitive.ObjectID `bson:"workspace_id" json:"workspace_id"`

	Name string `bson:"name" json:"name"`

	Position int `bson:"position" json:"position"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
type Column struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	BoardID primitive.ObjectID `bson:"board_id" json:"board_id"`

	Name string `bson:"name" json:"name"`

	Position int `bson:"position" json:"position"`

	Color string `bson:"color,omitempty" json:"color,omitempty"`

	IsArchived bool `bson:"is_archived" json:"is_archived"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`

	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
