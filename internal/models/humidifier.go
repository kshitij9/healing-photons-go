package models

import "time"

// User represents the user data model
type Humidifier struct {
	ID        string    `json:"id"`
	StockID   string    `json:"stock_id"`
	Weight    float32   `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
