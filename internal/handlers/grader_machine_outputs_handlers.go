package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllGraderMachineOutputs - Get all grader machine output records
func GetAllGraderMachineOutputs(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
        SELECT id, type 
        FROM grader_machine_outputs`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var outputs []models.GraderMachineOutputs
	for rows.Next() {
		var output models.GraderMachineOutputs
		if err := rows.Scan(
			&output.ID,
			&output.Type,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		outputs = append(outputs, output)
	}
	c.JSON(http.StatusOK, outputs)
}

// GetGraderMachineOutput - Get single grader machine output record
func GetGraderMachineOutput(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var output models.GraderMachineOutputs
	err := db.QueryRow(`
        SELECT id, type 
        FROM grader_machine_outputs WHERE id = ?`, id).Scan(
		&output.ID,
		&output.Type,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, output)
}

// CreateGraderMachineOutput - Create new grader machine output record
func CreateGraderMachineOutput(c *gin.Context, db *sql.DB) {
	var output models.GraderMachineOutputs
	if err := c.ShouldBindJSON(&output); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(`
        INSERT INTO grader_machine_outputs (
            id, type
        )
        VALUES (?, ?)`,
		output.ID,
		output.Type,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, output)
}

// UpdateGraderMachineOutput - Update existing grader machine output record
func UpdateGraderMachineOutput(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var output models.GraderMachineOutputs
	if err := c.ShouldBindJSON(&output); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        UPDATE grader_machine_outputs
        SET type = ?
        WHERE id = ?`,
		output.Type,
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

// DeleteGraderMachineOutput - Delete grader machine output record
func DeleteGraderMachineOutput(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM grader_machine_outputs WHERE id = ?", id)
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

// SetupGraderMachineOutputRoutes - Setup all routes for grader machine outputs
func SetupGraderMachineOutputRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/grader-machine-outputs", func(c *gin.Context) { GetAllGraderMachineOutputs(c, db) })
	router.GET("/grader-machine-outputs/:id", func(c *gin.Context) { GetGraderMachineOutput(c, db) })
	router.POST("/grader-machine-outputs", func(c *gin.Context) { CreateGraderMachineOutput(c, db) })
	router.PUT("/grader-machine-outputs/:id", func(c *gin.Context) { UpdateGraderMachineOutput(c, db) })
	router.DELETE("/grader-machine-outputs/:id", func(c *gin.Context) { DeleteGraderMachineOutput(c, db) })
}
