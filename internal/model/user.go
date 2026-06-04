package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email"         json:"email"`
	Password  string             `bson:"password"      json:"-"`
	Name      string             `bson:"name"          json:"name"`
	CreatedAt time.Time          `bson:"created_at"    json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"    json:"updated_at"`
}

type CreateUserRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name"     binding:"required"`
}

type UpdateUserRequest struct {
	Name string `json:"name" binding:"omitempty,min=1"`
}
