package store

import (
	"errors"
	"sync"
	"task-management-api/internal/models"
)

var ErrNotFound = errors.New("Task Not Found")

type TaskStore struct {
	mu    sync.Mutex
	tasks map[string]models.Task
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]models.Task),
	}
}

func (s *TaskStore) Create(task models.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
}

func (s *TaskStore) GetAll() []models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	return result
}

func (s *TaskStore) GetByID(id string) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, ErrNotFound
	}
	return task, nil
}

func (s *TaskStore) Update(id string, updated models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}
	s.tasks[id] = updated
	return nil
}

func (s *TaskStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}
