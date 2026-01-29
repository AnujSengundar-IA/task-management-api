package store

import (
	"context"
	"task-management-api/internal/models"
)

type TaskRepository interface {
	Create(ctx context.Context, task models.Task) error
	GetAll(ctx context.Context) ([]models.Task, error)
	GetByID(ctx context.Context, id string) (models.Task, error)
	Update(ctx context.Context, task models.Task) error
	Delete(ctx context.Context, id string) error
}
