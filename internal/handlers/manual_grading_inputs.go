package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllManualGradingInputs - Get all machine grading input records
func GetAllManualGradingInputs(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
        SELECT id, stock_id, worker_id, size_variations_id, weight, 
               created_at, updated_at 
        FROM machine_grading_inputs`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var inputs []models.ManualGradingInput
	for rows.Next() {
		var input models.ManualGradingInput
		if err := rows.Scan(
			&input.ID,
			&input.StockID,
			&input.WorkerID,
			&input.SizeVariationsID,
			&input.Weight,
			&input.CreatedAt,
			&input.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		inputs = append(inputs, input)
	}
	c.JSON(http.StatusOK, inputs)
}

// GetManualGradingInput - Get single machine grading input record
func GetManualGradingInput(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var input models.ManualGradingInput
	err := db.QueryRow(`
        SELECT id, stock_id, worker_id, size_variations_id, weight, 
               created_at, updated_at 
        FROM machine_grading_inputs WHERE id = ?`, id).Scan(
		&input.ID,
		&input.StockID,
		&input.WorkerID,
		&input.SizeVariationsID,
		&input.Weight,
		&input.CreatedAt,
		&input.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, input)
}

// CreateManualGradingInput - Create new machine grading input record
func CreateManualGradingInput(c *gin.Context, db *sql.DB) {
	var input models.ManualGradingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the record
	_, err := db.Exec(`
        INSERT INTO machine_grading_inputs (
            id, stock_id, worker_id, size_variations_id, weight,
            created_at, updated_at
        )
        VALUES (?, ?, ?, ?, ?, NOW(), NOW())`,
		input.ID,
		input.StockID,
		input.WorkerID,
		input.SizeVariationsID,
		input.Weight,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch the created record to get timestamps
	err = db.QueryRow(`
        SELECT id, stock_id, worker_id, size_variations_id, weight,
               created_at, updated_at 
        FROM machine_grading_inputs WHERE id = ?`, input.ID).Scan(
		&input.ID,
		&input.StockID,
		&input.WorkerID,
		&input.SizeVariationsID,
		&input.Weight,
		&input.CreatedAt,
		&input.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// UpdateManualGradingInput - Update existing machine grading input record
func UpdateManualGradingInput(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var input models.ManualGradingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        UPDATE machine_grading_inputs
        SET stock_id = ?,
            worker_id = ?,
            size_variations_id = ?,
            weight = ?,
            updated_at = NOW()
        WHERE id = ?`,
		input.StockID,
		input.WorkerID,
		input.SizeVariationsID,
		input.Weight,
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

// DeleteManualGradingInput - Delete machine grading input record
func DeleteManualGradingInput(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM machine_grading_inputs WHERE id = ?", id)
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

// GetManualGradingInputsByStock - Get machine grading input records for a specific stock ID
func GetManualGradingInputsByStock(c *gin.Context, db *sql.DB) {
	stockID := c.Param("stockId")

	rows, err := db.Query(`
        SELECT id, stock_id, worker_id, size_variations_id, weight, 
               created_at, updated_at 
        FROM machine_grading_inputs 
        WHERE stock_id = ?
        ORDER BY created_at DESC`, stockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var inputs []models.ManualGradingInput
	for rows.Next() {
		var input models.ManualGradingInput
		if err := rows.Scan(
			&input.ID,
			&input.StockID,
			&input.WorkerID,
			&input.SizeVariationsID,
			&input.Weight,
			&input.CreatedAt,
			&input.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		inputs = append(inputs, input)
	}

	if len(inputs) == 0 {
		c.JSON(http.StatusOK, []models.ManualGradingInput{}) // Return empty array instead of null
		return
	}

	c.JSON(http.StatusOK, inputs)
}

// SetupManualGradingInputRoutes - Setup all routes for machine grading inputs
func SetupManualGradingInputRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/manual-grading-inputs", func(c *gin.Context) { GetAllManualGradingInputs(c, db) })
	router.GET("/manual-grading-inputs/:id", func(c *gin.Context) { GetManualGradingInput(c, db) })
	router.POST("/manual-grading-inputs", func(c *gin.Context) { CreateManualGradingInput(c, db) })
	router.PUT("/manual-grading-inputs/:id", func(c *gin.Context) { UpdateManualGradingInput(c, db) })
	router.DELETE("/manual-grading-inputs/:id", func(c *gin.Context) { DeleteManualGradingInput(c, db) })
	router.GET("/manual-grading-inputs/stock/:stockId", func(c *gin.Context) { GetManualGradingInputsByStock(c, db) })
} 