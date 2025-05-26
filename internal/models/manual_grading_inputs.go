package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// ManualGradingInput represents the machine_grading_inputs table in the database
type ManualGradingInput struct {
	ID               int            `json:"id"`
	StockID          string         `json:"stock_id"`
	WorkerID         string         `json:"worker_id"`
	SizeVariationsID sql.NullInt64  `json:"size_variations_id,omitempty"`
	Weight           float64        `json:"weight"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// UnmarshalJSON implements custom JSON unmarshaling for ManualGradingInput
func (m *ManualGradingInput) UnmarshalJSON(data []byte) error {
	type Alias ManualGradingInput
	aux := &struct {
		SizeVariationsID *int64 `json:"size_variations_id,omitempty"`
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
	return nil
} 