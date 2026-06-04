package service

import (
	"context"
	"errors"
	"net/http"
	"time"
	"todolist/backend/internal/model"
	"todolist/backend/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type ServiceError struct {
	Code    int
	Message string
}

func (e *ServiceError) Error() string { return e.Message }

type UserService interface {
	Create(ctx context.Context, req model.CreateUserRequest) (*model.User, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*model.User, error)
	Update(ctx context.Context, id bson.ObjectID, req model.UpdateUserRequest) (*model.User, error)
	Delete(ctx context.Context, id bson.ObjectID) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, req model.CreateUserRequest) (*model.User, error) {
	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, &ServiceError{Code: http.StatusConflict, Message: "email already exists"}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &model.User{
		ID:        bson.NewObjectID(),
		Email:     req.Email,
		Password:  string(hashed),
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id bson.ObjectID) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &ServiceError{Code: http.StatusNotFound, Message: "user not found"}
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, id bson.ObjectID, req model.UpdateUserRequest) (*model.User, error) {
	user, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	update := bson.M{"updated_at": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}

	if err := s.repo.Update(ctx, user.ID, update); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *userService) Delete(ctx context.Context, id bson.ObjectID) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return &ServiceError{Code: http.StatusNotFound, Message: "user not found"}
	}
	return s.repo.Delete(ctx, id)
}

func IsNotFound(err error) bool {
	var svcErr *ServiceError
	return errors.As(err, &svcErr) && svcErr.Code == http.StatusNotFound
}

func StatusCode(err error) int {
	var svcErr *ServiceError
	if errors.As(err, &svcErr) {
		return svcErr.Code
	}
	return http.StatusInternalServerError
}
