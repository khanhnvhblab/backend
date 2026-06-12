package service

import (
	"context"
	"todolist/backend/internal/model"
	"todolist/backend/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DashboardService interface {
	GetStats(ctx context.Context, userID bson.ObjectID) (*model.DashboardStats, error)
}

type dashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetStats(ctx context.Context, userID bson.ObjectID) (*model.DashboardStats, error) {
	return s.repo.GetStats(ctx, userID)
}
