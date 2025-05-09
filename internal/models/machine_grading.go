package models

import (
	"time"
)

// MachineGrading represents the machine_grading table in the database
type MachineGrading struct {
	ID                     string    `json:"id"`
	ColorSortID            string    `json:"color_sort_id"`
	StockID                string    `json:"stock_id"`
	GraderMachineOutputsID int       `json:"grader_machine_outputs_id"`
	Weight                 float64   `json:"weight"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}
