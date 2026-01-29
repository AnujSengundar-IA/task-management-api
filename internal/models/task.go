package models

import "time"

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"string"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
