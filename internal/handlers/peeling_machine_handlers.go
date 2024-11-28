package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllPeelingMachineData - Get all peeling machine records
func GetAllPeelingMachineData(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
        SELECT id, humidifier_id, wholes, k, lwp, swp, bb, bbnp, 
               husk, stock_id, created_at, updated_at 
        FROM peeling_machine`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var machines []models.PeelingMachine
	for rows.Next() {
		var machine models.PeelingMachine
		if err := rows.Scan(
			&machine.ID,
			&machine.HumidifierID,
			&machine.Wholes,
			&machine.K,
			&machine.Lwp,
			&machine.Swp,
			&machine.Bb,
			&machine.Bbnp,
			&machine.Husk,
			&machine.StockID,
			&machine.CreatedAt,
			&machine.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		machines = append(machines, machine)
	}
	c.JSON(http.StatusOK, machines)
}

// GetPeelingMachine - Get single peeling machine record
func GetPeelingMachine(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var machine models.PeelingMachine
	err := db.QueryRow(`
        SELECT id, humidifier_id, wholes, k, lwp, swp, bb, bbnp, 
               husk, stock_id, created_at, updated_at 
        FROM peeling_machine WHERE id = $1`, id).Scan(
		&machine.ID,
		&machine.HumidifierID,
		&machine.Wholes,
		&machine.K,
		&machine.Lwp,
		&machine.Swp,
		&machine.Bb,
		&machine.Bbnp,
		&machine.Husk,
		&machine.StockID,
		&machine.CreatedAt,
		&machine.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, machine)
}

// CreatePeelingMachine - Create new peeling machine record
func CreatePeelingMachine(c *gin.Context, db *sql.DB) {
	var machine models.PeelingMachine
	if err := c.ShouldBindJSON(&machine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRow(`
        INSERT INTO peeling_machine (
            humidifier_id, wholes, k, lwp, swp, bb, bbnp, husk, stock_id, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
        RETURNING id, created_at, updated_at`,
		machine.HumidifierID,
		machine.Wholes,
		machine.K,
		machine.Lwp,
		machine.Swp,
		machine.Bb,
		machine.Bbnp,
		machine.Husk,
		machine.StockID,
	).Scan(&machine.ID, &machine.CreatedAt, &machine.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, machine)
}

// UpdatePeelingMachine - Update existing peeling machine record
func UpdatePeelingMachine(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var machine models.PeelingMachine
	if err := c.ShouldBindJSON(&machine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        UPDATE peeling_machine 
        SET humidifier_id = $1,
            wholes = $2,
            k = $3,
            lwp = $4,
            swp = $5,
            bb = $6,
            bbnp = $7,
            husk = $8,
            stock_id = $9,
            updated_at = NOW()
        WHERE id = $10`,
		machine.HumidifierID,
		machine.Wholes,
		machine.K,
		machine.Lwp,
		machine.Swp,
		machine.Bb,
		machine.Bbnp,
		machine.Husk,
		machine.StockID,
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

// DeletePeelingMachine - Delete peeling machine record
func DeletePeelingMachine(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM peeling_machine WHERE id = $1", id)
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

// SetupPeelingMachineRoutes - Setup all routes for peeling machine
func SetupPeelingMachineRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/peeling-machines", func(c *gin.Context) { GetAllPeelingMachineData(c, db) })
	router.GET("/peeling-machines/:id", func(c *gin.Context) { GetPeelingMachine(c, db) })
	router.POST("/peeling-machines", func(c *gin.Context) { CreatePeelingMachine(c, db) })
	router.PUT("/peeling-machines/:id", func(c *gin.Context) { UpdatePeelingMachine(c, db) })
	router.DELETE("/peeling-machines/:id", func(c *gin.Context) { DeletePeelingMachine(c, db) })
}
