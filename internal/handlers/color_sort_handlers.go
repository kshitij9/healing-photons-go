package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllColorSorts - Get all color sort records
func GetAllColorSorts(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
        SELECT id, stock_id, peel_id, acc_wholes, acc_k, acc_lwp, acc_swp, 
               acc_bb, acc_bbnp, acc_husk, created_at, updated_at 
        FROM color_sort`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var colorSorts []models.ColorSort
	for rows.Next() {
		var colorSort models.ColorSort
		if err := rows.Scan(
			&colorSort.ID,
			&colorSort.StockID,
			&colorSort.PeelId,
			&colorSort.AccWholes,
			&colorSort.AccK,
			&colorSort.AccLwp,
			&colorSort.AccSwp,
			&colorSort.AccBb,
			&colorSort.AccBbnp,
			&colorSort.AccHusk,
			&colorSort.CreatedAt,
			&colorSort.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		colorSorts = append(colorSorts, colorSort)
	}
	c.JSON(http.StatusOK, colorSorts)
}

// GetColorSort - Get single color sort record
func GetColorSort(c *gin.Context, db *sql.DB) {
	id := c.Param("stock_id")

	var colorSort models.ColorSort
	err := db.QueryRow(`
        SELECT id, stock_id, peel_id, acc_wholes, acc_k, acc_lwp, acc_swp, 
               acc_bb, acc_bbnp, acc_husk, created_at, updated_at 
        FROM color_sort WHERE stock_id = $1`, id).Scan(
		&colorSort.ID,
		&colorSort.StockID,
		&colorSort.PeelId,
		&colorSort.AccWholes,
		&colorSort.AccK,
		&colorSort.AccLwp,
		&colorSort.AccSwp,
		&colorSort.AccBb,
		&colorSort.AccBbnp,
		&colorSort.AccHusk,
		&colorSort.CreatedAt,
		&colorSort.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, colorSort)
}

// CreateColorSort - Create new color sort record
func CreateColorSort(c *gin.Context, db *sql.DB) {
	var colorSort models.ColorSort
	if err := c.ShouldBindJSON(&colorSort); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(`
        INSERT INTO color_sort (
            id, stock_id, peel_id, acc_wholes, acc_k, acc_lwp, acc_swp, 
            acc_bb, acc_bbnp, acc_husk, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
        ON CONFLICT (id) 
        DO UPDATE SET 
            stock_id = EXCLUDED.stock_id,
            peel_id = EXCLUDED.peel_id,
            acc_wholes = EXCLUDED.acc_wholes,
            acc_k = EXCLUDED.acc_k,
            acc_lwp = EXCLUDED.acc_lwp,
            acc_swp = EXCLUDED.acc_swp,
            acc_bb = EXCLUDED.acc_bb,
            acc_bbnp = EXCLUDED.acc_bbnp,
            acc_husk = EXCLUDED.acc_husk,
            updated_at = NOW()`,
		colorSort.ID,
		colorSort.StockID,
		colorSort.PeelId,
		colorSort.AccWholes,
		colorSort.AccK,
		colorSort.AccLwp,
		colorSort.AccSwp,
		colorSort.AccBb,
		colorSort.AccBbnp,
		colorSort.AccHusk,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, colorSort)
}

// UpdateColorSort - Update existing color sort record
func UpdateColorSort(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var colorSort models.ColorSort
	if err := c.ShouldBindJSON(&colorSort); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        UPDATE color_sort
        SET stock_id = $1,
            peel_id = $2,
            acc_wholes = $3,
            acc_k = $4,
            acc_lwp = $5,
            acc_swp = $6,
            acc_bb = $7,
            acc_bbnp = $8,
            acc_husk = $9,
            updated_at = NOW()
        WHERE id = $10`,
		colorSort.StockID,
		colorSort.PeelId,
		colorSort.AccWholes,
		colorSort.AccK,
		colorSort.AccLwp,
		colorSort.AccSwp,
		colorSort.AccBb,
		colorSort.AccBbnp,
		colorSort.AccHusk,
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

// DeleteColorSort - Delete color sort record
func DeleteColorSort(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM color_sort WHERE id = $1", id)
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

// SetupColorSortRoutes - Setup all routes for color sort
func SetupColorSortRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/color-sorts", func(c *gin.Context) { GetAllColorSorts(c, db) })
	router.GET("/color-sorts/:id", func(c *gin.Context) { GetColorSort(c, db) })
	router.POST("/color-sorts", func(c *gin.Context) { CreateColorSort(c, db) })
	router.PUT("/color-sorts/:id", func(c *gin.Context) { UpdateColorSort(c, db) })
	router.DELETE("/color-sorts/:id", func(c *gin.Context) { DeleteColorSort(c, db) })
}
