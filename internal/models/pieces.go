package models


// MachineGrading represents the machine_grading table in the database
type Pieces struct {
	PieceID                int    `json:"piece_id"`
	PieceCode              int    `json:"piece_code"`
	Description			   string `json:"description"`
}