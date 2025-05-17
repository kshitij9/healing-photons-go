package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllMachineGradings - Get all machine grading records
func GetAllMachineGradings(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
        SELECT id, color_sort_id, stock_id, size_variations_id, pieces_id,
               weight, created_at, updated_at 
        FROM machine_grading`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var gradings []models.MachineGrading
	for rows.Next() {
		var grading models.MachineGrading
		if err := rows.Scan(
			&grading.ID,
			&grading.ColorSortID,
			&grading.StockID,
			&grading.SizeVariationsID,
			&grading.PiecesID,
			&grading.Weight,
			&grading.CreatedAt,
			&grading.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		gradings = append(gradings, grading)
	}
	c.JSON(http.StatusOK, gradings)
}

// GetMachineGrading - Get single machine grading record
func GetMachineGrading(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var grading models.MachineGrading
	err := db.QueryRow(`
        SELECT id, color_sort_id, stock_id, size_variations_id, pieces_id,
               weight, created_at, updated_at 
        FROM machine_grading WHERE id = ?`, id).Scan(
		&grading.ID,
		&grading.ColorSortID,
		&grading.StockID,
		&grading.SizeVariationsID,
		&grading.PiecesID,
		&grading.Weight,
		&grading.CreatedAt,
		&grading.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, grading)
}

// CreateMachineGrading - Create new machine grading record
func CreateMachineGrading(c *gin.Context, db *sql.DB) {
	var grading models.MachineGrading
	if err := c.ShouldBindJSON(&grading); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRow(`
        INSERT INTO machine_grading (
            id, color_sort_id, stock_id, size_variations_id, pieces_id,
            weight, created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
        RETURNING created_at, updated_at`,
		grading.ID,
		grading.ColorSortID,
		grading.StockID,
		grading.SizeVariationsID,
		grading.PiecesID,
		grading.Weight,
	).Scan(&grading.CreatedAt, &grading.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, grading)
}

// UpdateMachineGrading - Update existing machine grading record
func UpdateMachineGrading(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var grading models.MachineGrading
	if err := c.ShouldBindJSON(&grading); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        UPDATE machine_grading
        SET color_sort_id = ?,
            stock_id = ?,
            size_variations_id = ?,
            pieces_id = ?,
            weight = ?,
            updated_at = NOW()
        WHERE id = ?`,
		grading.ColorSortID,
		grading.StockID,
		grading.SizeVariationsID,
		grading.PiecesID,
		grading.Weight,
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

// DeleteMachineGrading - Delete machine grading record
func DeleteMachineGrading(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM machine_grading WHERE id = ?", id)
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

// GetMachineGradingsByStock - Get machine grading records for a specific stock ID
func GetMachineGradingsByStock(c *gin.Context, db *sql.DB) {
	stockID := c.Param("stockId")

	rows, err := db.Query(`
        SELECT id, color_sort_id, stock_id, size_variations_id, pieces_id,
               weight, created_at, updated_at 
        FROM machine_grading 
        WHERE stock_id = ?
        ORDER BY created_at DESC`, stockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var gradings []models.MachineGrading
	for rows.Next() {
		var grading models.MachineGrading
		if err := rows.Scan(
			&grading.ID,
			&grading.ColorSortID,
			&grading.StockID,
			&grading.SizeVariationsID,
			&grading.PiecesID,
			&grading.Weight,
			&grading.CreatedAt,
			&grading.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		gradings = append(gradings, grading)
	}

	if len(gradings) == 0 {
		c.JSON(http.StatusOK, []models.MachineGrading{}) // Return empty array instead of null
		return
	}

	c.JSON(http.StatusOK, gradings)
}

// GetWeightSummary - Get summary of weights for a stock ID
func GetWeightSummary(c *gin.Context, db *sql.DB) {
	stockID := c.Param("stockId")

	var summary struct {
		StockID     string  `json:"stock_id"`
		TotalWeight float64 `json:"total_weight"`
		RecordCount int     `json:"record_count"`
	}

	err := db.QueryRow(`
        SELECT 
            stock_id,
            COALESCE(SUM(weight), 0) as total_weight,
            COUNT(*) as record_count
        FROM machine_grading 
        WHERE stock_id = ?
        GROUP BY stock_id`,
		stockID,
	).Scan(
		&summary.StockID,
		&summary.TotalWeight,
		&summary.RecordCount,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{
			"stock_id":     stockID,
			"total_weight": 0,
			"record_count": 0,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// SetupMachineGradingRoutes - Setup all routes for machine grading
func SetupMachineGradingRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/machine-gradings", func(c *gin.Context) { GetAllMachineGradings(c, db) })
	router.GET("/machine-gradings/:id", func(c *gin.Context) { GetMachineGrading(c, db) })
	router.POST("/machine-gradings", func(c *gin.Context) { CreateMachineGrading(c, db) })
	router.PUT("/machine-gradings/:id", func(c *gin.Context) { UpdateMachineGrading(c, db) })
	router.DELETE("/machine-gradings/:id", func(c *gin.Context) { DeleteMachineGrading(c, db) })
	router.GET("/machine-gradings/stock/:stockId", func(c *gin.Context) { GetMachineGradingsByStock(c, db) })
	router.GET("/machine-gradings/stock/:stockId/summary", func(c *gin.Context) { GetWeightSummary(c, db) })
}
