package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures the API routes
func SetupRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/stocks", func(c *gin.Context) { GetAllStocks(c, db) })
	router.GET("/stocks/:id", func(c *gin.Context) { GetStock(c, db) })
	router.POST("/stocks", func(c *gin.Context) { CreateStock(c, db) })
	router.PUT("/stocks/:id", func(c *gin.Context) { UpdateStock(c, db) })
	router.DELETE("/stocks/:id", func(c *gin.Context) { DeleteStock(c, db) })
}

// GetStock - Get single stock
func GetStock(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var stock models.Stock
	err := db.QueryRow(`
		SELECT stock_id, seller_name, origin_country, weight, date, created_at, updated_at 
		FROM stocks WHERE stock_id = $1`, id).Scan(
		&stock.StockID,
		&stock.SellerName,
		&stock.OriginCountry,
		&stock.Weight,
		&stock.Date,
		&stock.CreatedAt,
		&stock.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stock)
}

// CreateStock - Create new stock
func CreateStock(c *gin.Context, db *sql.DB) {
	var stock models.Stock
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRow(`
		INSERT INTO stock (stock_id, seller_name, origin_country, weight, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING stock_id, created_at, updated_at`,
		stock.SellerName,
		stock.OriginCountry,
		stock.Weight,
		stock.Date,
	).Scan(&stock.StockID, &stock.CreatedAt, &stock.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, stock)
}

// UpdateStock - Update existing stock
func UpdateStock(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var stock models.Stock
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
		UPDATE stock
		SET seller_name = $1, 
			origin_country = $2, 
			weight = $3, 
			date = $4, 
			updated_at = NOW()
		WHERE stock_id = $6`,
		stock.SellerName,
		stock.OriginCountry,
		stock.Weight,
		stock.Date,
		id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}

// DeleteStock - Delete stock
func DeleteStock(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM stock WHERE stock_id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Stock deleted successfully"})
}

// GetAllStocks - Get all stocks
func GetAllStocks(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
		SELECT stock_id, seller_name, origin_country, weight, date, created_at, updated_at 
		FROM stock`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.StockID,
			&stock.SellerName,
			&stock.OriginCountry,
			&stock.Weight,
			&stock.Date,
			&stock.CreatedAt,
			&stock.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		stocks = append(stocks, stock)
	}

	c.JSON(http.StatusOK, stocks)
}
