package task

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	ColumnID primitive.ObjectID `bson:"column_id" json:"column_id"`

	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`

	Priority string `bson:"priority" json:"priority"`

	Position int `bson:"position" json:"position"`

	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by"`

	AssignedMembers []primitive.ObjectID `bson:"assigned_members" json:"assigned_members"`

	Labels []string `bson:"labels" json:"labels"`

	DueDate *time.Time `bson:"due_date,omitempty" json:"due_date,omitempty"`

	IsArchived bool `bson:"is_archived" json:"is_archived"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
