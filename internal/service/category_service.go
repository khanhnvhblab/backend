package service

import (
	"context"
	"net/http"
	"time"
	"todolist/backend/internal/model"
	"todolist/backend/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const defaultCategoryColor = "#6366f1"

type CategoryService interface {
	Create(ctx context.Context, userID bson.ObjectID, req model.CreateCategoryRequest) (*model.Category, error)
	GetByID(ctx context.Context, userID bson.ObjectID, id string) (*model.Category, error)
	GetAll(ctx context.Context, userID bson.ObjectID, filter repository.CategoryFilter) ([]model.Category, int64, error)
	Update(ctx context.Context, userID bson.ObjectID, id string, req model.UpdateCategoryRequest) (*model.Category, error)
	Delete(ctx context.Context, userID bson.ObjectID, id string) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(ctx context.Context, userID bson.ObjectID, req model.CreateCategoryRequest) (*model.Category, error) {
	color := req.Color
	if color == "" {
		color = defaultCategoryColor
	}

	cat := &model.Category{
		UserID: userID,
		Name:   req.Name,
		Color:  color,
	}

	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) GetByID(ctx context.Context, userID bson.ObjectID, id string) (*model.Category, error) {
	catID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, &ServiceError{Code: http.StatusBadRequest, Message: "invalid id"}
	}

	cat, err := s.repo.FindByID(ctx, catID, userID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, &ServiceError{Code: http.StatusNotFound, Message: "category not found"}
	}
	return cat, nil
}

func (s *categoryService) GetAll(ctx context.Context, userID bson.ObjectID, filter repository.CategoryFilter) ([]model.Category, int64, error) {
	return s.repo.FindAll(ctx, userID, filter)
}

func (s *categoryService) Update(ctx context.Context, userID bson.ObjectID, id string, req model.UpdateCategoryRequest) (*model.Category, error) {
	catID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, &ServiceError{Code: http.StatusBadRequest, Message: "invalid id"}
	}

	cat, err := s.repo.FindByID(ctx, catID, userID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, &ServiceError{Code: http.StatusNotFound, Message: "category not found"}
	}

	update := bson.M{"updated_at": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Color != "" {
		update["color"] = req.Color
	}

	if err := s.repo.Update(ctx, catID, userID, update); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, catID, userID)
}

func (s *categoryService) Delete(ctx context.Context, userID bson.ObjectID, id string) error {
	catID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return &ServiceError{Code: http.StatusBadRequest, Message: "invalid id"}
	}

	cat, err := s.repo.FindByID(ctx, catID, userID)
	if err != nil {
		return err
	}
	if cat == nil {
		return &ServiceError{Code: http.StatusNotFound, Message: "category not found"}
	}

	return s.repo.Delete(ctx, catID, userID)
}
