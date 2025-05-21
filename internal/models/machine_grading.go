package models

import (
	"database/sql"
	"encoding/json"
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

// UnmarshalJSON implements custom JSON unmarshaling for MachineGrading
func (m *MachineGrading) UnmarshalJSON(data []byte) error {
	type Alias MachineGrading
	aux := &struct {
		SizeVariationsID *int64 `json:"size_variations_id,omitempty"`
		PiecesID         *int64 `json:"pieces_id,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.SizeVariationsID != nil {
		m.SizeVariationsID = sql.NullInt64{Int64: *aux.SizeVariationsID, Valid: true}
	}
	if aux.PiecesID != nil {
		m.PiecesID = sql.NullInt64{Int64: *aux.PiecesID, Valid: true}
	}
	return nil
}
