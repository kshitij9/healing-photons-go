package models

import (
	"database/sql"
	"encoding/json"
)

// GradingCategory represents the grading_categories table
type GradingCategory struct {
	CategoryID   int64          `json:"category_id"`
	CategoryCode string         `json:"category_code"`
	Description  sql.NullString `json:"description"`
}

// UnmarshalJSON implements custom JSON unmarshaling for GradingCategory
func (g *GradingCategory) UnmarshalJSON(data []byte) error {
	type Alias GradingCategory
	aux := &struct {
		Description *string `json:"description"`
		*Alias
	}{
		Alias: (*Alias)(g),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Description != nil {
		g.Description = sql.NullString{String: *aux.Description, Valid: true}
	} else {
		g.Description = sql.NullString{Valid: false}
	}

	return nil
} 