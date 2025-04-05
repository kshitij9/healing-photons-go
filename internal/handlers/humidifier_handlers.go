package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllHumidifiers - Get all humidifier records
func GetAllHumidifiers(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
        SELECT id, stock_id, weight, created_at, updated_at 
        FROM humidifier`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var humidifiers []models.Humidifier
	for rows.Next() {
		var humidifier models.Humidifier
		if err := rows.Scan(
			&humidifier.ID,
			&humidifier.StockID,
			&humidifier.Weight,
			&humidifier.CreatedAt,
			&humidifier.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		humidifiers = append(humidifiers, humidifier)
	}
	c.JSON(http.StatusOK, humidifiers)
}

// GetHumidifier - Get single humidifier record
func GetHumidifier(c *gin.Context, db *sql.DB) {
	id := c.Param("stock_id")

	var humidifier models.Humidifier
	err := db.QueryRow(`
        SELECT id, stock_id, weight, created_at, updated_at 
        FROM humidifier WHERE stock_id = ?`, id).Scan(
		&humidifier.ID,
		&humidifier.StockID,
		&humidifier.Weight,
		&humidifier.CreatedAt,
		&humidifier.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, humidifier)
}

// CreateHumidifier - Create new humidifier record
func CreateHumidifier(c *gin.Context, db *sql.DB) {
	var humidifier models.Humidifier
	if err := c.ShouldBindJSON(&humidifier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRow(`
        INSERT INTO humidifier (
            id, stock_id, weight, created_at, updated_at
        )
        VALUES (?, ?, ?, NOW(), NOW())`,
		humidifier.ID,
		humidifier.StockID,
		humidifier.Weight,
	).Scan(&humidifier.CreatedAt, &humidifier.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, humidifier)
}

// UpdateHumidifier - Update existing humidifier record
func UpdateHumidifier(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var humidifier models.Humidifier
	if err := c.ShouldBindJSON(&humidifier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        UPDATE humidifier
        SET stock_id = ?,
            weight = ?,
            updated_at = NOW()
        WHERE id = ?`,
		humidifier.StockID,
		humidifier.Weight,
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

// DeleteHumidifier - Delete humidifier record
func DeleteHumidifier(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM humidifier WHERE id = ?", id)
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

// SetupHumidifierRoutes - Setup all routes for humidifier
func SetupHumidifierRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/humidifiers", func(c *gin.Context) { GetAllHumidifiers(c, db) })
	router.GET("/humidifiers/:id", func(c *gin.Context) { GetHumidifier(c, db) })
	router.POST("/humidifiers", func(c *gin.Context) { CreateHumidifier(c, db) })
	router.PUT("/humidifiers/:id", func(c *gin.Context) { UpdateHumidifier(c, db) })
	router.DELETE("/humidifiers/:id", func(c *gin.Context) { DeleteHumidifier(c, db) })
}
