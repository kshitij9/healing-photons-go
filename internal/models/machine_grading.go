package models

import (
	"database/sql"
	"time"
)

// MachineGrading represents the machine_grading table in the database
type MachineGrading struct {
	ID                string         `json:"id"`
	ColorSortID       string         `json:"color_sort_id"`
	StockID           string         `json:"stock_id"`
	SizeVariationsID  sql.NullInt64  `json:"size_variations_id,omitempty"`
	PiecesID          sql.NullInt64  `json:"pieces_id,omitempty"`
	Weight            float64        `json:"weight"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}
