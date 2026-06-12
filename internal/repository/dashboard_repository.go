package repository

import (
	"context"
	"time"
	"todolist/backend/internal/db"
	"todolist/backend/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DashboardRepository interface {
	GetStats(ctx context.Context, userID bson.ObjectID) (*model.DashboardStats, error)
}

type mongoDashboardRepository struct{}

func NewDashboardRepository() DashboardRepository {
	return &mongoDashboardRepository{}
}

type facetCount struct {
	Count int64 `bson:"count"`
}

type facetResult struct {
	Total      []facetCount `bson:"total"`
	Pending    []facetCount `bson:"pending"`
	InProgress []facetCount `bson:"in_progress"`
	Done       []facetCount `bson:"done"`
	Overdue    []facetCount `bson:"overdue"`
	DueSoon    []facetCount `bson:"due_soon"`
}

func first(arr []facetCount) int64 {
	if len(arr) == 0 {
		return 0
	}
	return arr[0].Count
}

func (r *mongoDashboardRepository) GetStats(ctx context.Context, userID bson.ObjectID) (*model.DashboardStats, error) {
	now := time.Now()
	soon := now.AddDate(0, 0, 7)

	pipeline := bson.A{
		bson.M{"$match": bson.M{"user_id": userID}},
		bson.M{"$facet": bson.M{
			"total": bson.A{
				bson.M{"$count": "count"},
			},
			"pending": bson.A{
				bson.M{"$match": bson.M{"status": "pending"}},
				bson.M{"$count": "count"},
			},
			"in_progress": bson.A{
				bson.M{"$match": bson.M{"status": "in_progress"}},
				bson.M{"$count": "count"},
			},
			"done": bson.A{
				bson.M{"$match": bson.M{"status": "done"}},
				bson.M{"$count": "count"},
			},
			"overdue": bson.A{
				bson.M{"$match": bson.M{
					"deadline": bson.M{"$ne": nil, "$lt": now},
					"status":   bson.M{"$ne": "done"},
				}},
				bson.M{"$count": "count"},
			},
			"due_soon": bson.A{
				bson.M{"$match": bson.M{
					"deadline": bson.M{"$gte": now, "$lte": soon},
					"status":   bson.M{"$ne": "done"},
				}},
				bson.M{"$count": "count"},
			},
		}},
	}

	cursor, err := db.Col("todos").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []facetResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return &model.DashboardStats{}, nil
	}

	r2 := results[0]
	return &model.DashboardStats{
		Total:      first(r2.Total),
		Pending:    first(r2.Pending),
		InProgress: first(r2.InProgress),
		Done:       first(r2.Done),
		Overdue:    first(r2.Overdue),
		DueSoon:    first(r2.DueSoon),
	}, nil
}
