package models

import "time"

type ColorSort struct {
	ID          string    `json:"id"`
	StockID     string    `json:"stock_id"`
	PeelId      string    `json:"peel_id"`
	AccWholes   float32   `json:"acc_wholes"`
	AccK        float32   `json:"acc_k"`
	AccLwp      float32   `json:"acc_lwp"`
	AccSwp      float32   `json:"acc_swp"`
	AccBb       float32   `json:"acc_bb"`
	AccBbnp     float32   `json:"acc_bbnp"`
	AccHusk     float32   `json:"acc_husk"`
	SortCounter int       `json:"sort_counter"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
