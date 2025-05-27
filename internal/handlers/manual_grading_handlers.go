package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllManualGradings - Get all manual grading records
func GetAllManualGradings(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
		SELECT id, grader_machine_outputs_id, stock_id, category_id, 
			size_id, piece_id, weight, worker_id, created_at, updated_at
		FROM manual_grading
		ORDER BY created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var gradings []models.ManualGrading
	for rows.Next() {
		var grading models.ManualGrading
		if err := rows.Scan(
			&grading.ID,
			&grading.GraderMachineOutputsID,
			&grading.StockID,
			&grading.CategoryID,
			&grading.SizeID,
			&grading.PieceID,
			&grading.Weight,
			&grading.WorkerID,
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

// GetManualGrading - Get single manual grading record
func GetManualGrading(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var grading models.ManualGrading
	err := db.QueryRow(`
		SELECT id, grader_machine_outputs_id, stock_id, category_id, 
			size_id, piece_id, weight, worker_id, created_at, updated_at
		FROM manual_grading WHERE id = ?`, id).Scan(
		&grading.ID,
		&grading.GraderMachineOutputsID,
		&grading.StockID,
		&grading.CategoryID,
		&grading.SizeID,
		&grading.PieceID,
		&grading.Weight,
		&grading.WorkerID,
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

// CreateManualGrading - Create new manual grading record
func CreateManualGrading(c *gin.Context, db *sql.DB) {
	var grading models.ManualGrading
	if err := c.ShouldBindJSON(&grading); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the record
	_, err := db.Exec(`
		INSERT INTO manual_grading (
			id, grader_machine_outputs_id, stock_id, category_id, 
			size_id, piece_id, weight, worker_id, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		grading.GraderMachineOutputsID,
		grading.StockID,
		grading.CategoryID,
		grading.SizeID,
		grading.PieceID,
		grading.Weight,
		grading.WorkerID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch the created record to get timestamps
	err = db.QueryRow(`
        SELECT id, grader_machine_outputs_id, stock_id, category_id, 
			size_id, piece_id, weight, worker_id, created_at, updated_at
        FROM manual_grading WHERE id = ?`, grading.ID).Scan(
		&grading.ID,
		&grading.GraderMachineOutputsID,
		&grading.StockID,
		&grading.CategoryID,
		&grading.SizeID,
		&grading.PieceID,
		&grading.Weight,
		&grading.WorkerID,
		&grading.CreatedAt,
		&input.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// UpdateManualGrading - Update existing manual grading record
func UpdateManualGrading(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var grading models.ManualGrading
	if err := c.ShouldBindJSON(&grading); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
		UPDATE manual_grading
		SET grader_machine_outputs_id = ?,
			stock_id = ?,
			category_id = ?,
			size_id = ?,
			piece_id = ?,
			weight = ?,
			worker_id = ?,
			updated_at = NOW()
		WHERE id = ?`,
		grading.GraderMachineOutputsID,
		grading.StockID,
		grading.CategoryID,
		grading.SizeID,
		grading.PieceID,
		grading.Weight,
		grading.WorkerID,
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

// DeleteManualGrading - Delete manual grading record
func DeleteManualGrading(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM manual_grading WHERE id = ?", id)
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

// GetManualGradingsByStock - Get manual grading records for a specific stock ID
func GetManualGradingsByStock(c *gin.Context, db *sql.DB) {
	stockID := c.Param("stockId")

	rows, err := db.Query(`
		SELECT id, grader_machine_outputs_id, stock_id, category_id, 
			size_id, piece_id, weight, worker_id, created_at, updated_at
		FROM manual_grading 
		WHERE stock_id = ?
		ORDER BY created_at DESC`, stockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var gradings []models.ManualGrading
	for rows.Next() {
		var grading models.ManualGrading
		if err := rows.Scan(
			&grading.ID,
			&grading.GraderMachineOutputsID,
			&grading.StockID,
			&grading.CategoryID,
			&grading.SizeID,
			&grading.PieceID,
			&grading.Weight,
			&grading.WorkerID,
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

// SetupManualGradingRoutes sets up all the routes for manual grading
func SetupManualGradingRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/manual-grading", func(c *gin.Context) { GetAllManualGradings(c, db) })
	router.GET("/manual-grading/:id", func(c *gin.Context) { GetManualGrading(c, db) })
	router.POST("/manual-grading", func(c *gin.Context) { CreateManualGrading(c, db) })
	router.PUT("/manual-grading/:id", func(c *gin.Context) { UpdateManualGrading(c, db) })
	router.DELETE("/manual-grading/:id", func(c *gin.Context) { DeleteManualGrading(c, db) })
	router.GET("/manual-grading/stock/:stockId", func(c *gin.Context) { GetManualGradingsByStock(c, db) })
} 