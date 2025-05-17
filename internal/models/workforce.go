package models


// MachineGrading represents the machine_grading table in the database
type Workforce struct {
	ID                     string    `json:"id"`
	Name            	   string    `json:"name"`
	Aadhaar                float64   `json:"aadhaar"`
	Addresss 			   string    `json:"address"`
}