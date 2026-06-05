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

type CategoryFilter struct {
	Page  int
	Limit int
}

type CategoryRepository interface {
	Create(ctx context.Context, doc *model.Category) error
	FindByID(ctx context.Context, id, userID bson.ObjectID) (*model.Category, error)
	FindAll(ctx context.Context, userID bson.ObjectID, filter CategoryFilter) ([]model.Category, int64, error)
	Update(ctx context.Context, id, userID bson.ObjectID, update bson.M) error
	Delete(ctx context.Context, id, userID bson.ObjectID) error
}

type mongoCategoryRepository struct{}

func NewCategoryRepository() CategoryRepository {
	return &mongoCategoryRepository{}
}

func (r *mongoCategoryRepository) Create(ctx context.Context, doc *model.Category) error {
	doc.ID = bson.NewObjectID()
	now := time.Now()
	doc.CreatedAt = now
	doc.UpdatedAt = now
	_, err := db.Col("categories").InsertOne(ctx, doc)
	return err
}

func (r *mongoCategoryRepository) FindByID(ctx context.Context, id, userID bson.ObjectID) (*model.Category, error) {
	var cat model.Category
	err := db.Col("categories").FindOne(ctx, bson.M{"_id": id, "user_id": userID}).Decode(&cat)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &cat, err
}

func (r *mongoCategoryRepository) FindAll(ctx context.Context, userID bson.ObjectID, f CategoryFilter) ([]model.Category, int64, error) {
	filter := bson.M{"user_id": userID}

	total, err := db.Col("categories").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Page <= 0 {
		f.Page = 1
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((f.Page - 1) * f.Limit)).
		SetLimit(int64(f.Limit))

	cursor, err := db.Col("categories").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var categories []model.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (r *mongoCategoryRepository) Update(ctx context.Context, id, userID bson.ObjectID, update bson.M) error {
	_, err := db.Col("categories").UpdateOne(ctx, bson.M{"_id": id, "user_id": userID}, bson.M{"$set": update})
	return err
}

func (r *mongoCategoryRepository) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	_, err := db.Col("categories").DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
	return err
}
