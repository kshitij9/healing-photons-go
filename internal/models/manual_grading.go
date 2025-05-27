package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// GradeCategory represents grading measurements for a specific size category
type GradeCategory struct {
	Whole float64 `json:"whole"` // Whole kernels
	A     float64 `json:"a"`     // K grade
	SW    float64 `json:"sw"`    // LWP grade
	SSW   float64 `json:"ssw"`   // SWP grade
	TW    float64 `json:"tw"`    // BB grade
	JB    float64 `json:"jb"`    // BBNP grade
	KW    float64 `json:"kw"`    // Husk weight
}

// ManualGrading represents the manual_grading table
type ManualGrading struct {
	ID                    string         `json:"id"`
	GraderMachineOutputsID string        `json:"grader_machine_outputs_id"`
	StockID              string         `json:"stock_id"`
	CategoryID           sql.NullInt64  `json:"category_id"`
	SizeID               int64          `json:"size_id"`
	PieceID              sql.NullInt64  `json:"piece_id"`
	Weight               int64          `json:"weight"`
	WorkerID             string         `json:"worker_id"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

// UnmarshalJSON implements custom JSON unmarshaling for ManualGrading
func (m *ManualGrading) UnmarshalJSON(data []byte) error {
	type Alias ManualGrading
	aux := &struct {
		CategoryID *int64 `json:"category_id"`
		PieceID    *int64 `json:"piece_id"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.CategoryID != nil {
		m.CategoryID = sql.NullInt64{Int64: *aux.CategoryID, Valid: true}
	} else {
		m.CategoryID = sql.NullInt64{Valid: false}
	}

	if aux.PieceID != nil {
		m.PieceID = sql.NullInt64{Int64: *aux.PieceID, Valid: true}
	} else {
		m.PieceID = sql.NullInt64{Valid: false}
	}

	return nil
}
