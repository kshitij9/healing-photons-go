package models

import "time"

// User represents the user data model
type Stock struct {
	StockID       int       `json:"stock_id"`
	SellerName    string    `json:"seller_name"`
	OriginCountry string    `json:"origin_country"`
	Weight        float32   `json:"weight"`
	Date          time.Time `json:"date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
