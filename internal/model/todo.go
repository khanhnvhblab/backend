package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	TodoStatusPending    = "pending"
	TodoStatusInProgress = "in_progress"
	TodoStatusDone       = "done"

	TodoPriorityLow    = "low"
	TodoPriorityMedium = "medium"
	TodoPriorityHigh   = "high"
)

type Todo struct {
	ID          bson.ObjectID  `bson:"_id,omitempty"          json:"id"`
	UserID      bson.ObjectID  `bson:"user_id"                json:"-"`
	CategoryID  *bson.ObjectID `bson:"category_id,omitempty"  json:"category_id"`
	Title       string         `bson:"title"                  json:"title"`
	Description string         `bson:"description"            json:"description"`
	Status      string         `bson:"status"                 json:"status"`
	Priority    string         `bson:"priority"               json:"priority"`
	Deadline    *time.Time     `bson:"deadline,omitempty"     json:"deadline"`
	CompletedAt *time.Time     `bson:"completed_at,omitempty" json:"completed_at"`
	CreatedAt   time.Time      `bson:"created_at"             json:"created_at"`
	UpdatedAt   time.Time      `bson:"updated_at"             json:"updated_at"`
}

type CreateTodoRequest struct {
	Title       string     `json:"title"       binding:"required"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	CategoryID  string     `json:"category_id"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateTodoRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	CategoryID  string     `json:"category_id"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateTodoStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
