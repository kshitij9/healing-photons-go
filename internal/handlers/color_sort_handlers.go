package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllColorSorts - Get all color sort records
func GetAllColorSorts(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
        SELECT id, peel_id, stock_id, weight_type_id, accepted_weight, 
               sort_counter, created_at, updated_at 
        FROM color_sort`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var colorSorts []models.ColorSort
	for rows.Next() {
		var colorSort models.ColorSort
		if err := rows.Scan(
			&colorSort.ID,
			&colorSort.PeelID,
			&colorSort.StockID,
			&colorSort.WeightTypeID,
			&colorSort.AcceptedWeight,
			&colorSort.SortCounter,
			&colorSort.CreatedAt,
			&colorSort.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		colorSorts = append(colorSorts, colorSort)
	}
	c.JSON(http.StatusOK, colorSorts)
}

// GetColorSort - Get single color sort record
func GetColorSort(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var colorSort models.ColorSort
	err := db.QueryRow(`
        SELECT id, peel_id, stock_id, weight_type_id, accepted_weight, 
               sort_counter, created_at, updated_at 
        FROM color_sort WHERE id = ?`, id).Scan(
		&colorSort.ID,
		&colorSort.PeelID,
		&colorSort.StockID,
		&colorSort.WeightTypeID,
		&colorSort.AcceptedWeight,
		&colorSort.SortCounter,
		&colorSort.CreatedAt,
		&colorSort.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, colorSort)
}

// CreateColorSort - Create new color sort record
func CreateColorSort(c *gin.Context, db *sql.DB) {
	var colorSort models.ColorSort
	if err := c.ShouldBindJSON(&colorSort); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRow(`
        INSERT INTO color_sort (
            id, peel_id, stock_id, weight_type_id, accepted_weight, 
            sort_counter, created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
        RETURNING created_at, updated_at`,
		colorSort.ID,
		colorSort.PeelID,
		colorSort.StockID,
		colorSort.WeightTypeID,
		colorSort.AcceptedWeight,
		colorSort.SortCounter,
	).Scan(&colorSort.CreatedAt, &colorSort.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, colorSort)
}

// UpdateColorSort - Update existing color sort record
func UpdateColorSort(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var colorSort models.ColorSort
	if err := c.ShouldBindJSON(&colorSort); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        UPDATE color_sort
        SET peel_id = ?,
            stock_id = ?,
            weight_type_id = ?,
            accepted_weight = ?,
            sort_counter = ?,
            updated_at = NOW()
        WHERE id = ?`,
		colorSort.PeelID,
		colorSort.StockID,
		colorSort.WeightTypeID,
		colorSort.AcceptedWeight,
		colorSort.SortCounter,
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record updated successfully"})
}

// DeleteColorSort - Delete color sort record
func DeleteColorSort(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM color_sort WHERE id = ?", id)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// GetColorSortsByStock - Get color sort records for a specific stock ID with optional counter filter
func GetColorSortsByStock(c *gin.Context, db *sql.DB) {
	stockID := c.Param("stockId")
	counter := c.Query("counter") // Optional query parameter

	var query string
	var args []interface{}

	if counter != "" {
		query = `
            SELECT id, peel_id, stock_id, weight_type_id, accepted_weight, 
                   sort_counter, created_at, updated_at 
            FROM color_sort 
            WHERE stock_id = ? AND sort_counter = ?
            ORDER BY created_at DESC`
		args = []interface{}{stockID, counter}
	} else {
		query = `
            SELECT id, peel_id, stock_id, weight_type_id, accepted_weight, 
                   sort_counter, created_at, updated_at 
            FROM color_sort 
            WHERE stock_id = ?
            ORDER BY created_at DESC`
		args = []interface{}{stockID}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var colorSorts []models.ColorSort
	for rows.Next() {
		var colorSort models.ColorSort
		if err := rows.Scan(
			&colorSort.ID,
			&colorSort.PeelID,
			&colorSort.StockID,
			&colorSort.WeightTypeID,
			&colorSort.AcceptedWeight,
			&colorSort.SortCounter,
			&colorSort.CreatedAt,
			&colorSort.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		colorSorts = append(colorSorts, colorSort)
	}

	if len(colorSorts) == 0 {
		c.JSON(http.StatusOK, []models.ColorSort{}) // Return empty array instead of null
		return
	}

	c.JSON(http.StatusOK, colorSorts)
}

// GetAcceptedWeightSummary - Get summary of accepted weights for a stock ID and counter
func GetAcceptedWeightSummary(c *gin.Context, db *sql.DB) {
	stockID := c.Param("stockId")
	counter := c.Param("counter")

	var summary struct {
		StockID       string  `json:"stock_id"`
		SortCounter   int     `json:"sort_counter"`
		TotalAccepted float64 `json:"total_accepted_weight"`
		RecordCount   int     `json:"record_count"`
	}

	err := db.QueryRow(`
        SELECT 
            stock_id,
            sort_counter,
            COALESCE(SUM(accepted_weight), 0) as total_accepted_weight,
            COUNT(*) as record_count
        FROM color_sort 
        WHERE stock_id = ? AND sort_counter = ?
        GROUP BY stock_id, sort_counter`,
		stockID, counter,
	).Scan(
		&summary.StockID,
		&summary.SortCounter,
		&summary.TotalAccepted,
		&summary.RecordCount,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{
			"stock_id":              stockID,
			"sort_counter":          counter,
			"total_accepted_weight": 0,
			"record_count":          0,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetColorSortsByStockAndCounter - Get color sort records for a specific stock ID and sort counter
func GetColorSortsByStockAndCounter(c *gin.Context, db *sql.DB) {
	stockID := c.Param("stockId")
	counter := c.Param("counter")

	rows, err := db.Query(`
        SELECT id, peel_id, stock_id, weight_type_id, accepted_weight, 
               sort_counter, created_at, updated_at 
        FROM color_sort 
        WHERE stock_id = ? AND sort_counter = ?
        ORDER BY created_at DESC`,
		stockID, counter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var colorSorts []models.ColorSort
	for rows.Next() {
		var colorSort models.ColorSort
		if err := rows.Scan(
			&colorSort.ID,
			&colorSort.PeelID,
			&colorSort.StockID,
			&colorSort.WeightTypeID,
			&colorSort.AcceptedWeight,
			&colorSort.SortCounter,
			&colorSort.CreatedAt,
			&colorSort.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		colorSorts = append(colorSorts, colorSort)
	}

	if len(colorSorts) == 0 {
		c.JSON(http.StatusOK, []models.ColorSort{}) // Return empty array instead of null
		return
	}

	c.JSON(http.StatusOK, colorSorts)
}

// SetupColorSortRoutes - Setup all routes for color sort
func SetupColorSortRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/color-sorts", func(c *gin.Context) { GetAllColorSorts(c, db) })
	router.GET("/color-sorts/:id", func(c *gin.Context) { GetColorSort(c, db) })
	router.POST("/color-sorts", func(c *gin.Context) { CreateColorSort(c, db) })
	router.PUT("/color-sorts/:id", func(c *gin.Context) { UpdateColorSort(c, db) })
	router.DELETE("/color-sorts/:id", func(c *gin.Context) { DeleteColorSort(c, db) })
	router.GET("/color-sorts/stock/:stockId", func(c *gin.Context) { GetColorSortsByStock(c, db) })
	router.GET("/color-sorts/stock/:stockId/counter/:counter", func(c *gin.Context) { GetColorSortsByStockAndCounter(c, db) })
	router.GET("/color-sorts/stock/:stockId/counter/:counter/summary", func(c *gin.Context) { GetAcceptedWeightSummary(c, db) })
}
