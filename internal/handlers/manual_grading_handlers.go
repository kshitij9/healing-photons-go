package handlers

import (
	"database/sql"
	"encoding/json"
	"healing_photons/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SetupManualGradingRoutes sets up all the routes for manual grading
func SetupManualGradingRoutes(router *gin.Engine, db *sql.DB) {
	manualGrading := router.Group("/api/manual-grading")
	{
		manualGrading.POST("", createManualGrading(db))
		manualGrading.GET("", getAllManualGradings(db))
		manualGrading.GET("/:id", getManualGradingByID(db))
		manualGrading.PUT("/:id", updateManualGrading(db))
		manualGrading.DELETE("/:id", deleteManualGrading(db))
	}
}

// createManualGrading handles the creation of a new manual grading record
func createManualGrading(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var manualGrading models.ManualGrading
		if err := c.ShouldBindJSON(&manualGrading); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate UUID for new record
		manualGrading.ID = uuid.New().String()
		manualGrading.CreatedAt = time.Now()
		manualGrading.UpdatedAt = time.Now()

		// Insert into database
		query := `
			INSERT INTO manual_grading (
				id, grader_machine_outputs_id, stock_id, category_id, 
				size_id, piece_id, weight, worker_id, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err := db.Exec(query,
			manualGrading.ID,
			manualGrading.GraderMachineOutputsID,
			manualGrading.StockID,
			manualGrading.CategoryID,
			manualGrading.SizeID,
			manualGrading.PieceID,
			manualGrading.Weight,
			manualGrading.WorkerID,
			manualGrading.CreatedAt,
			manualGrading.UpdatedAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, manualGrading)
	}
}

// getAllManualGradings retrieves all manual grading records
func getAllManualGradings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT id, grader_machine_outputs_id, stock_id, category_id, 
				size_id, piece_id, weight, worker_id, created_at, updated_at
			FROM manual_grading
			ORDER BY created_at DESC
		`
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var manualGradings []models.ManualGrading
		for rows.Next() {
			var mg models.ManualGrading
			err := rows.Scan(
				&mg.ID,
				&mg.GraderMachineOutputsID,
				&mg.StockID,
				&mg.CategoryID,
				&mg.SizeID,
				&mg.PieceID,
				&mg.Weight,
				&mg.WorkerID,
				&mg.CreatedAt,
				&mg.UpdatedAt,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			manualGradings = append(manualGradings, mg)
		}

		c.JSON(http.StatusOK, manualGradings)
	}
}

// getManualGradingByID retrieves a specific manual grading record by ID
func getManualGradingByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		query := `
			SELECT id, grader_machine_outputs_id, stock_id, category_id, 
				size_id, piece_id, weight, worker_id, created_at, updated_at
			FROM manual_grading
			WHERE id = ?
		`
		var mg models.ManualGrading
		err := db.QueryRow(query, id).Scan(
			&mg.ID,
			&mg.GraderMachineOutputsID,
			&mg.StockID,
			&mg.CategoryID,
			&mg.SizeID,
			&mg.PieceID,
			&mg.Weight,
			&mg.WorkerID,
			&mg.CreatedAt,
			&mg.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Manual grading not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, mg)
	}
}

// updateManualGrading updates an existing manual grading record
func updateManualGrading(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var manualGrading models.ManualGrading
		if err := c.ShouldBindJSON(&manualGrading); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		manualGrading.ID = id
		manualGrading.UpdatedAt = time.Now()

		query := `
			UPDATE manual_grading
			SET grader_machine_outputs_id = ?, stock_id = ?, category_id = ?,
				size_id = ?, piece_id = ?, weight = ?, worker_id = ?, updated_at = ?
			WHERE id = ?
		`
		result, err := db.Exec(query,
			manualGrading.GraderMachineOutputsID,
			manualGrading.StockID,
			manualGrading.CategoryID,
			manualGrading.SizeID,
			manualGrading.PieceID,
			manualGrading.Weight,
			manualGrading.WorkerID,
			manualGrading.UpdatedAt,
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
			c.JSON(http.StatusNotFound, gin.H{"error": "Manual grading not found"})
			return
		}

		c.JSON(http.StatusOK, manualGrading)
	}
}

// deleteManualGrading deletes a manual grading record
func deleteManualGrading(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		query := "DELETE FROM manual_grading WHERE id = ?"
		result, err := db.Exec(query, id)
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
			c.JSON(http.StatusNotFound, gin.H{"error": "Manual grading not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Manual grading deleted successfully"})
	}
} 