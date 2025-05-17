package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllSizeVariations - Get all size variations records
func GetAllSizeVariations(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
        SELECT size_id, size_value
        FROM size_variations`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var variationsList []models.SizeVariations
	for rows.Next() {
		var variation models.SizeVariations
		if err := rows.Scan(
			&variation.SizeID,
			&variation.SizeValue,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		variationsList = append(variationsList, variation)
	}
	c.JSON(http.StatusOK, variationsList)
}

// GetSizeVariation - Get single size variation record
func GetSizeVariation(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var variation models.SizeVariations
	err := db.QueryRow(`
        SELECT size_id, size_value
        FROM size_variations WHERE size_id = ?`, id).Scan(
		&variation.SizeID,
		&variation.SizeValue,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, variation)
}

// CreateSizeVariation - Create new size variation record
func CreateSizeVariation(c *gin.Context, db *sql.DB) {
	var variation models.SizeVariations
	if err := c.ShouldBindJSON(&variation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        INSERT INTO size_variations (
            size_id, size_value
        )
        VALUES (?, ?)`,
		variation.SizeID,
		variation.SizeValue,
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

	variation.SizeID = int(lastID)
	c.JSON(http.StatusCreated, variation)
}

// UpdateSizeVariation - Update existing size variation record
func UpdateSizeVariation(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var variation models.SizeVariations
	if err := c.ShouldBindJSON(&variation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        UPDATE size_variations
        SET size_value = ?
        WHERE size_id = ?`,
		variation.SizeValue,
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

// DeleteSizeVariation - Delete size variation record
func DeleteSizeVariation(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM size_variations WHERE size_id = ?", id)
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

// SetupSizeVariationsRoutes - Setup all routes for size variations
func SetupSizeVariationsRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/size-variations", func(c *gin.Context) { GetAllSizeVariations(c, db) })
	router.GET("/size-variations/:id", func(c *gin.Context) { GetSizeVariation(c, db) })
	router.POST("/size-variations", func(c *gin.Context) { CreateSizeVariation(c, db) })
	router.PUT("/size-variations/:id", func(c *gin.Context) { UpdateSizeVariation(c, db) })
	router.DELETE("/size-variations/:id", func(c *gin.Context) { DeleteSizeVariation(c, db) })
} 