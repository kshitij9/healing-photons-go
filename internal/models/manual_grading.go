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

// ToFlatMap converts the nested structure to a flat map for database operations
func (m *ManualGrading) ToFlatMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         m.ID,
		"peeling_id": m.PeelingID,
		"stock_id":   m.StockID,
		// 180
		"w180":   m.Size180.Whole,
		"a180":   m.Size180.A,
		"sw180":  m.Size180.SW,
		"ssw180": m.Size180.SSW,
		"tw180":  m.Size180.TW,
		"jb180":  m.Size180.JB,
		"kw180":  m.Size180.KW,
		// 210
		"w210":   m.Size210.Whole,
		"a210":   m.Size210.A,
		"sw210":  m.Size210.SW,
		"ssw210": m.Size210.SSW,
		"tw210":  m.Size210.TW,
		"jb210":  m.Size210.JB,
		"kw210":  m.Size210.KW,
		// 240
		"w240":   m.Size240.Whole,
		"a240":   m.Size240.A,
		"sw240":  m.Size240.SW,
		"ssw240": m.Size240.SSW,
		"tw240":  m.Size240.TW,
		"jb240":  m.Size240.JB,
		"kw240":  m.Size240.KW,
		// 280
		"w280":   m.Size280.Whole,
		"a280":   m.Size280.A,
		"sw280":  m.Size280.SW,
		"ssw280": m.Size280.SSW,
		"tw280":  m.Size280.TW,
		"jb280":  m.Size280.JB,
		"kw280":  m.Size280.KW,
		// 320
		"w320":   m.Size320.Whole,
		"a320":   m.Size320.A,
		"sw320":  m.Size320.SW,
		"ssw320": m.Size320.SSW,
		"tw320":  m.Size320.TW,
		"jb320":  m.Size320.JB,
		"kw320":  m.Size320.KW,
		// 400
		"w400":   m.Size400.Whole,
		"a400":   m.Size400.A,
		"sw400":  m.Size400.SW,
		"ssw400": m.Size400.SSW,
		"tw400":  m.Size400.TW,
		"jb400":  m.Size400.JB,
		"kw400":  m.Size400.KW,
	}
}

// FromScanValues populates the struct from database scan values
func (m *ManualGrading) FromScanValues(values ...interface{}) {
	m.ID = values[0].(string)
	m.PeelingID = values[1].(string)
	m.StockID = values[2].(string)

	// 180
	m.Size180.Whole = values[3].(float64)
	m.Size180.A = values[4].(float64)
	m.Size180.SW = values[5].(float64)
	m.Size180.SSW = values[6].(float64)
	m.Size180.TW = values[7].(float64)
	m.Size180.JB = values[8].(float64)
	m.Size180.KW = values[9].(float64)

	// 210
	m.Size210.Whole = values[10].(float64)
	m.Size210.A = values[11].(float64)
	m.Size210.SW = values[12].(float64)
	m.Size210.SSW = values[13].(float64)
	m.Size210.TW = values[14].(float64)
	m.Size210.JB = values[15].(float64)
	m.Size210.KW = values[16].(float64)

	// Continue for other sizes...
	// ... (similar pattern for 240, 280, 320, 400)

	m.CreatedAt = values[44].(time.Time)
	m.UpdatedAt = values[45].(time.Time)
}
