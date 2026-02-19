package models

import (
	"time"
)

type Query struct {
	ID              string    `json:"id"`
	CadastralNumber string    `json:"cadastral_number"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	Status          string    `json:"status"`
	Result          *bool     `json:"result,omitempty"`
	UserID          string    `json:"user_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	CompletedAt     time.Time `json:"completed_at,omitempty"`
}

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
