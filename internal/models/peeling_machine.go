package models

import "time"

type PeelingMachine struct {
	ID           string    `json:"id" db:"id"`
	HumidifierID string    `json:"humidifier_id" db:"humidifier_id"`
	StockID      *string   `json:"stock_id,omitempty" db:"stock_id"`
	WeightTypeID int       `json:"weight_type_id" db:"weight_type_id"`
	Weight       float64   `json:"weight" db:"weight"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
