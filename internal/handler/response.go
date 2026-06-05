package handler

import (
	"todolist/backend/internal/model"
	"todolist/backend/internal/service"
)

// Swagger response envelope types

type userResp struct {
	Data model.User `json:"data"`
}

type registerResp struct {
	Data registerData `json:"data"`
}

type registerData struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type loginResp struct {
	Data service.LoginResponse `json:"data"`
}

type refreshResp struct {
	Data service.RefreshResponse `json:"data"`
}

type errResp struct {
	Error string `json:"error"`
}

type todoResp struct {
	Data model.Todo `json:"data"`
}

type todoListResp struct {
	Data  []model.Todo `json:"data"`
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Limit int          `json:"limit"`
}

type categoryResp struct {
	Data model.Category `json:"data"`
}

type categoryListResp struct {
	Data  []model.Category `json:"data"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}
