package service

import (
	"context"
	"net/http"
	"time"
	"todolist/backend/internal/model"
	"todolist/backend/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TodoService interface {
	Create(ctx context.Context, userID bson.ObjectID, req model.CreateTodoRequest) (*model.Todo, error)
	GetByID(ctx context.Context, userID bson.ObjectID, id string) (*model.Todo, error)
	GetAll(ctx context.Context, userID bson.ObjectID, filter repository.TodoFilter) ([]model.Todo, int64, error)
	Update(ctx context.Context, userID bson.ObjectID, id string, req model.UpdateTodoRequest) (*model.Todo, error)
	UpdateStatus(ctx context.Context, userID bson.ObjectID, id string, status string) (*model.Todo, error)
	Delete(ctx context.Context, userID bson.ObjectID, id string) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) Create(ctx context.Context, userID bson.ObjectID, req model.CreateTodoRequest) (*model.Todo, error) {
	priority := req.Priority
	if priority == "" {
		priority = model.TodoPriorityMedium
	}

	todo := &model.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      model.TodoStatusPending,
		Priority:    priority,
		Deadline:    req.Deadline,
	}

	if req.CategoryID != "" {
		catID, err := bson.ObjectIDFromHex(req.CategoryID)
		if err != nil {
			return nil, &ServiceError{Code: http.StatusBadRequest, Message: "invalid category_id"}
		}
		todo.CategoryID = &catID
	}

	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *todoService) GetByID(ctx context.Context, userID bson.ObjectID, id string) (*model.Todo, error) {
	todoID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, &ServiceError{Code: http.StatusBadRequest, Message: "invalid id"}
	}

	todo, err := s.repo.FindByID(ctx, todoID, userID)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, &ServiceError{Code: http.StatusNotFound, Message: "todo not found"}
	}
	return todo, nil
}

func (s *todoService) GetAll(ctx context.Context, userID bson.ObjectID, filter repository.TodoFilter) ([]model.Todo, int64, error) {
	return s.repo.FindAll(ctx, userID, filter)
}

func (s *todoService) Update(ctx context.Context, userID bson.ObjectID, id string, req model.UpdateTodoRequest) (*model.Todo, error) {
	todoID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, &ServiceError{Code: http.StatusBadRequest, Message: "invalid id"}
	}

	todo, err := s.repo.FindByID(ctx, todoID, userID)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, &ServiceError{Code: http.StatusNotFound, Message: "todo not found"}
	}

	update := bson.M{"updated_at": time.Now()}
	if req.Title != "" {
		update["title"] = req.Title
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.Priority != "" {
		update["priority"] = req.Priority
	}
	if req.CategoryID != "" {
		catID, err := bson.ObjectIDFromHex(req.CategoryID)
		if err != nil {
			return nil, &ServiceError{Code: http.StatusBadRequest, Message: "invalid category_id"}
		}
		update["category_id"] = catID
	}
	if req.Deadline != nil {
		update["deadline"] = req.Deadline
	}

	if err := s.repo.Update(ctx, todoID, userID, update); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, todoID, userID)
}

func (s *todoService) UpdateStatus(ctx context.Context, userID bson.ObjectID, id string, status string) (*model.Todo, error) {
	if status != model.TodoStatusPending && status != model.TodoStatusInProgress && status != model.TodoStatusDone {
		return nil, &ServiceError{Code: http.StatusBadRequest, Message: "invalid status"}
	}

	todoID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, &ServiceError{Code: http.StatusBadRequest, Message: "invalid id"}
	}

	todo, err := s.repo.FindByID(ctx, todoID, userID)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, &ServiceError{Code: http.StatusNotFound, Message: "todo not found"}
	}

	update := bson.M{
		"status":     status,
		"updated_at": time.Now(),
	}
	if status == model.TodoStatusDone {
		now := time.Now()
		update["completed_at"] = now
	} else {
		update["completed_at"] = nil
	}

	if err := s.repo.Update(ctx, todoID, userID, update); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, todoID, userID)
}

func (s *todoService) Delete(ctx context.Context, userID bson.ObjectID, id string) error {
	todoID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return &ServiceError{Code: http.StatusBadRequest, Message: "invalid id"}
	}

	todo, err := s.repo.FindByID(ctx, todoID, userID)
	if err != nil {
		return err
	}
	if todo == nil {
		return &ServiceError{Code: http.StatusNotFound, Message: "todo not found"}
	}

	return s.repo.Delete(ctx, todoID, userID)
}
