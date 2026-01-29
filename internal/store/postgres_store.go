package store

import (
	"context"
	"database/sql"
	"errors"
	"task-management-api/internal/models"
)

type PostgresTaskStore struct {
	db *sql.DB
}

func NewPostgresTaskStore(db *sql.DB) *PostgresTaskStore {
	return &PostgresTaskStore{db: db}
}

func (s *PostgresTaskStore) Create(ctx context.Context, task models.Task) error {
	query := `
		INSERT INTO tasks (id, title, status, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := s.db.ExecContext(
		ctx,
		query,
		task.ID,
		task.Title,
		task.Status,
		task.CreatedAt,
	)
	return err

}

func (s *PostgresTaskStore) GetAll(ctx context.Context) ([]models.Task, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, title, status, created_at FROM tasks`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Status,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (s *PostgresTaskStore) GetByID(ctx context.Context, id string) (models.Task, error) {
	var t models.Task

	err := s.db.QueryRowContext(
		ctx,
		`SELECT id, title, status, created_at FROM tasks WHERE id = $1`,
		id,
	).Scan(&t.ID, &t.Title, &t.Status, &t.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Task{}, ErrNotFound
	}
	return t, err
}

func (s *PostgresTaskStore) Update(ctx context.Context, task models.Task) error {
	res, err := s.db.ExecContext(
		ctx,
		`UPDATE tasks SET title=$1, status=$2, WHERE id=$3`,
		task.Title,
		task.Status,
		task.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostgresTaskStore) Delete(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(
		ctx,
		`DELETE FROM tasks WHERE id=$1`,
		id,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
