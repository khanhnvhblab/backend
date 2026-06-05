package repository

import (
	"context"
	"time"
	"todolist/backend/internal/db"
	"todolist/backend/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TodoFilter struct {
	Status     string
	Priority   string
	CategoryID string
	Search     string
	SortBy     string
	Order      string
	Page       int
	Limit      int
}

type TodoRepository interface {
	Create(ctx context.Context, doc *model.Todo) error
	FindByID(ctx context.Context, id, userID bson.ObjectID) (*model.Todo, error)
	FindAll(ctx context.Context, userID bson.ObjectID, filter TodoFilter) ([]model.Todo, int64, error)
	Update(ctx context.Context, id, userID bson.ObjectID, update bson.M) error
	Delete(ctx context.Context, id, userID bson.ObjectID) error
}

type mongoTodoRepository struct{}

func NewTodoRepository() TodoRepository {
	return &mongoTodoRepository{}
}

func (r *mongoTodoRepository) Create(ctx context.Context, doc *model.Todo) error {
	doc.ID = bson.NewObjectID()
	now := time.Now()
	doc.CreatedAt = now
	doc.UpdatedAt = now
	_, err := db.Col("todos").InsertOne(ctx, doc)
	return err
}

func (r *mongoTodoRepository) FindByID(ctx context.Context, id, userID bson.ObjectID) (*model.Todo, error) {
	var todo model.Todo
	err := db.Col("todos").FindOne(ctx, bson.M{"_id": id, "user_id": userID}).Decode(&todo)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &todo, err
}

func (r *mongoTodoRepository) FindAll(ctx context.Context, userID bson.ObjectID, f TodoFilter) ([]model.Todo, int64, error) {
	filter := bson.M{"user_id": userID}
	if f.Status != "" {
		filter["status"] = f.Status
	}
	if f.Priority != "" {
		filter["priority"] = f.Priority
	}
	if f.CategoryID != "" {
		if catID, err := bson.ObjectIDFromHex(f.CategoryID); err == nil {
			filter["category_id"] = catID
		}
	}
	if f.Search != "" {
		filter["title"] = bson.M{"$regex": f.Search, "$options": "i"}
	}

	total, err := db.Col("todos").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Page <= 0 {
		f.Page = 1
	}

	sortOrder := -1
	if f.Order == "asc" {
		sortOrder = 1
	}
	sortField := "created_at"
	if f.SortBy == "deadline" {
		sortField = "deadline"
	}

	opts := options.Find().
		SetSort(bson.D{{Key: sortField, Value: sortOrder}}).
		SetSkip(int64((f.Page - 1) * f.Limit)).
		SetLimit(int64(f.Limit))

	cursor, err := db.Col("todos").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var todos []model.Todo
	if err := cursor.All(ctx, &todos); err != nil {
		return nil, 0, err
	}
	return todos, total, nil
}

func (r *mongoTodoRepository) Update(ctx context.Context, id, userID bson.ObjectID, update bson.M) error {
	_, err := db.Col("todos").UpdateOne(ctx, bson.M{"_id": id, "user_id": userID}, bson.M{"$set": update})
	return err
}

func (r *mongoTodoRepository) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	_, err := db.Col("todos").DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
	return err
}
