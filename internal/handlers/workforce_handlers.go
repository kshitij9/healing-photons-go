package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllWorkforce - Get all workforce records
func GetAllWorkforce(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
        SELECT id, name, aadhaar, address
        FROM workforce`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var workforceList []models.Workforce
	for rows.Next() {
		var workforce models.Workforce
		if err := rows.Scan(
			&workforce.ID,
			&workforce.Name,
			&workforce.Aadhaar,
			&workforce.Addresss,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		workforceList = append(workforceList, workforce)
	}
	c.JSON(http.StatusOK, workforceList)
}

// GetWorkforce - Get single workforce record
func GetWorkforce(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var workforce models.Workforce
	err := db.QueryRow(`
        SELECT id, name, aadhaar, address
        FROM workforce WHERE id = ?`, id).Scan(
		&workforce.ID,
		&workforce.Name,
		&workforce.Aadhaar,
		&workforce.Addresss,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workforce)
}

// CreateWorkforce - Create new workforce record
func CreateWorkforce(c *gin.Context, db *sql.DB) {
	var workforce models.Workforce
	if err := c.ShouldBindJSON(&workforce); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        INSERT INTO workforce (
            id, name, aadhaar, address
        )
        VALUES (?, ?, ?, ?)`,
		workforce.ID,
		workforce.Name,
		workforce.Aadhaar,
		workforce.Addresss,
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

	workforce.ID = string(lastID)
	c.JSON(http.StatusCreated, workforce)
}

// UpdateWorkforce - Update existing workforce record
func UpdateWorkforce(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var workforce models.Workforce
	if err := c.ShouldBindJSON(&workforce); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        UPDATE workforce
        SET name = ?,
            aadhaar = ?,
            address = ?
        WHERE id = ?`,
		workforce.Name,
		workforce.Aadhaar,
		workforce.Addresss,
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

// DeleteWorkforce - Delete workforce record
func DeleteWorkforce(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM workforce WHERE id = ?", id)
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

// SetupWorkforceRoutes - Setup all routes for workforce
func SetupWorkforceRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/workforce", func(c *gin.Context) { GetAllWorkforce(c, db) })
	router.GET("/workforce/:id", func(c *gin.Context) { GetWorkforce(c, db) })
	router.POST("/workforce", func(c *gin.Context) { CreateWorkforce(c, db) })
	router.PUT("/workforce/:id", func(c *gin.Context) { UpdateWorkforce(c, db) })
	router.DELETE("/workforce/:id", func(c *gin.Context) { DeleteWorkforce(c, db) })
} 