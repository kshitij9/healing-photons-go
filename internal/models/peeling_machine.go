package models

import "time"

// User represents the user data model
type PeelingMachine struct {
	ID           int       `json:"id"`
	HumidifierID int       `json:"humidifier_id"`
	Wholes       float32   `json:"wholes"`
	K            float32   `json:"k"`
	Lwp          float32   `json:"lwp"`
	Swp          float32   `json:"swp"`
	Bb           float32   `json:"bb"`
	Bbnp         float32   `json:"bbnp"`
	Husk         float32   `json:"husk"`
	StockID      int       `json:"stock_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
