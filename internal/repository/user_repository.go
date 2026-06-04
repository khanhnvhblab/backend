package repository

import (
	"context"
	"todolist/backend/internal/db"
	"todolist/backend/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, doc *model.User) error
	FindByID(ctx context.Context, id bson.ObjectID) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, id bson.ObjectID, update bson.M) error
	Delete(ctx context.Context, id bson.ObjectID) error
}

type mongoUserRepository struct{}

func NewUserRepository() UserRepository {
	return &mongoUserRepository{}
}

func (r *mongoUserRepository) Create(ctx context.Context, doc *model.User) error {
	_, err := db.Col("users").InsertOne(ctx, doc)
	return err
}

func (r *mongoUserRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.User, error) {
	var user model.User
	err := db.Col("users").FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *mongoUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := db.Col("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *mongoUserRepository) Update(ctx context.Context, id bson.ObjectID, update bson.M) error {
	_, err := db.Col("users").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *mongoUserRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := db.Col("users").DeleteOne(ctx, bson.M{"_id": id})
	return err
}
