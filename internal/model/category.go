package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Category struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectID `bson:"user_id"       json:"-"`
	Name      string        `bson:"name"          json:"name"`
	Color     string        `bson:"color"         json:"color"`
	CreatedAt time.Time     `bson:"created_at"    json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name  string `json:"name"  binding:"required"`
	Color string `json:"color"`
}

type UpdateCategoryRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}
