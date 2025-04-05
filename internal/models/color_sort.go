package models

import "time"

type ColorSort struct {
	ID             string    `json:"id" db:"id"`
	PeelID         *string   `json:"peel_id,omitempty" db:"peel_id"`
	StockID        *string   `json:"stock_id,omitempty" db:"stock_id"`
	WeightTypeID   int       `json:"weight_type_id" db:"weight_type_id"`
	AcceptedWeight float64   `json:"accepted_weight" db:"accepted_weight"`
	SortCounter    int       `json:"sort_counter" db:"sort_counter"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
