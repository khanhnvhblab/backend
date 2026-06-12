package model

type DashboardStats struct {
	Total      int64 `json:"total"`
	Pending    int64 `json:"pending"`
	InProgress int64 `json:"in_progress"`
	Done       int64 `json:"done"`
	Overdue    int64 `json:"overdue"`
	DueSoon    int64 `json:"due_soon"`
}
