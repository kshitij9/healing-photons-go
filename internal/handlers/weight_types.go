package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllWeightTypes - Get all weight type records
func GetAllWeightTypes(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
        SELECT id, type 
        FROM weight_types`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var weightTypes []models.WeightTypes
	for rows.Next() {
		var weightType models.WeightTypes
		if err := rows.Scan(
			&weightType.ID,
			&weightType.Type,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		weightTypes = append(weightTypes, weightType)
	}
	c.JSON(http.StatusOK, weightTypes)
}

// GetWeightType - Get single weight type record
func GetWeightType(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var weightType models.WeightTypes
	err := db.QueryRow(`
        SELECT id, type 
        FROM weight_types WHERE id = ?`, id).Scan(
		&weightType.ID,
		&weightType.Type,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, weightType)
}

// CreateWeightType - Create new weight type record
func CreateWeightType(c *gin.Context, db *sql.DB) {
	var weightType models.WeightTypes
	if err := c.ShouldBindJSON(&weightType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        INSERT INTO weight_types (id, type)
        VALUES (?, ?)`,
		weightType.ID,
		weightType.Type,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	weightType.ID = string(lastID)
	c.JSON(http.StatusCreated, weightType)
}

// UpdateWeightType - Update existing weight type record
func UpdateWeightType(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var weightType models.WeightTypes
	if err := c.ShouldBindJSON(&weightType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        UPDATE weight_types
        SET type = ?
        WHERE id = ?`,
		weightType.Type,
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

// DeleteWeightType - Delete weight type record
func DeleteWeightType(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM weight_types WHERE id = ?", id)
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

// GetWeightTypesByUsage - Get weight types with their usage count
func GetWeightTypesByUsage(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
        SELECT wt.id, wt.type, COUNT(mg.id) as usage_count
        FROM weight_types wt
        LEFT JOIN machine_grading mg ON wt.id = mg.weight_type_id
        GROUP BY wt.id, wt.type
        ORDER BY usage_count DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var weightTypes []struct {
		models.WeightTypes
		UsageCount int `json:"usage_count"`
	}

	for rows.Next() {
		var wt struct {
			models.WeightTypes
			UsageCount int `json:"usage_count"`
		}
		if err := rows.Scan(
			&wt.ID,
			&wt.Type,
			&wt.UsageCount,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		weightTypes = append(weightTypes, wt)
	}

	if len(weightTypes) == 0 {
		c.JSON(http.StatusOK, []struct{}{}) // Return empty array instead of null
		return
	}

	c.JSON(http.StatusOK, weightTypes)
}

// SetupWeightTypeRoutes - Setup all routes for weight types
func SetupWeightTypeRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/weight-types", func(c *gin.Context) { GetAllWeightTypes(c, db) })
	router.GET("/weight-types/:id", func(c *gin.Context) { GetWeightType(c, db) })
	router.POST("/weight-types", func(c *gin.Context) { CreateWeightType(c, db) })
	router.PUT("/weight-types/:id", func(c *gin.Context) { UpdateWeightType(c, db) })
	router.DELETE("/weight-types/:id", func(c *gin.Context) { DeleteWeightType(c, db) })
	router.GET("/weight-types/usage", func(c *gin.Context) { GetWeightTypesByUsage(c, db) })
}
