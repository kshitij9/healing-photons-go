package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllGradingCategories - Get all grading categories
func GetAllGradingCategories(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
		SELECT category_id, category_code, description
		FROM grading_categories
		ORDER BY category_code`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var categories []models.GradingCategory
	for rows.Next() {
		var category models.GradingCategory
		if err := rows.Scan(
			&category.CategoryID,
			&category.CategoryCode,
			&category.Description,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		categories = append(categories, category)
	}
	c.JSON(http.StatusOK, categories)
}

// GetGradingCategory - Get single grading category
func GetGradingCategory(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var category models.GradingCategory
	err := db.QueryRow(`
		SELECT category_id, category_code, description
		FROM grading_categories 
		WHERE category_id = ?`, id).Scan(
		&category.CategoryID,
		&category.CategoryCode,
		&category.Description,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// CreateGradingCategory - Create new grading category
func CreateGradingCategory(c *gin.Context, db *sql.DB) {
	var category models.GradingCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the record
	_, err := db.Exec(`
		INSERT INTO grading_categories (
			category_id, category_code, description
		)
		VALUES (?, ?, ?)`,
		category.CategoryID,
		category.CategoryCode,
		category.Description,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch the created record to get timestamps
	err = db.QueryRow(`
		SELECT category_id, category_code, description
		FROM grading_categories WHERE category_id = ?`, category.CategoryID).Scan(
		&category.CategoryID,
		&category.CategoryCode,
		&category.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// UpdateGradingCategory - Update existing grading category
func UpdateGradingCategory(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var category models.GradingCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
		UPDATE grading_categories
		SET category_code = ?,
			description = ?
		WHERE category_id = ?`,
		category.CategoryCode,
		category.Description,
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

// DeleteGradingCategory - Delete grading category
func DeleteGradingCategory(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM grading_categories WHERE category_id = ?", id)
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

// SetupGradingCategoryRoutes sets up all the routes for grading categories
func SetupGradingCategoryRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/grading-categories", func(c *gin.Context) { GetAllGradingCategories(c, db) })
	router.GET("/grading-categories/:id", func(c *gin.Context) { GetGradingCategory(c, db) })
	router.POST("/grading-categories", func(c *gin.Context) { CreateGradingCategory(c, db) })
	router.PUT("/grading-categories/:id", func(c *gin.Context) { UpdateGradingCategory(c, db) })
	router.DELETE("/grading-categories/:id", func(c *gin.Context) { DeleteGradingCategory(c, db) })
} 