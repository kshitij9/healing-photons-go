package handlers

import (
	"database/sql"
	"healing_photons/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllPieces - Get all pieces records
func GetAllPieces(c *gin.Context, db *sql.DB) {
	rows, err := db.Query(`
        SELECT piece_id, piece_code, description
        FROM pieces`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var piecesList []models.Pieces
	for rows.Next() {
		var piece models.Pieces
		if err := rows.Scan(
			&piece.PieceID,
			&piece.PieceCode,
			&piece.Description,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		piecesList = append(piecesList, piece)
	}
	c.JSON(http.StatusOK, piecesList)
}

// GetPiece - Get single piece record
func GetPiece(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var piece models.Pieces
	err := db.QueryRow(`
        SELECT piece_id, piece_code, description
        FROM pieces WHERE piece_id = ?`, id).Scan(
		&piece.PieceID,
		&piece.PieceCode,
		&piece.Description,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, piece)
}

// CreatePiece - Create new piece record
func CreatePiece(c *gin.Context, db *sql.DB) {
	var piece models.Pieces
	if err := c.ShouldBindJSON(&piece); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        INSERT INTO pieces (
            piece_id, piece_code, description
        )
        VALUES (?, ?, ?)`,
		piece.PieceID,
		piece.PieceCode,
		piece.Description,
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

	piece.PieceID = int(lastID)
	c.JSON(http.StatusCreated, piece)
}

// UpdatePiece - Update existing piece record
func UpdatePiece(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var piece models.Pieces
	if err := c.ShouldBindJSON(&piece); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(`
        UPDATE pieces
        SET piece_code = ?,
            description = ?
        WHERE piece_id = ?`,
		piece.PieceCode,
		piece.Description,
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

// DeletePiece - Delete piece record
func DeletePiece(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM pieces WHERE piece_id = ?", id)
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

// SetupPiecesRoutes - Setup all routes for pieces
func SetupPiecesRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/pieces", func(c *gin.Context) { GetAllPieces(c, db) })
	router.GET("/pieces/:id", func(c *gin.Context) { GetPiece(c, db) })
	router.POST("/pieces", func(c *gin.Context) { CreatePiece(c, db) })
	router.PUT("/pieces/:id", func(c *gin.Context) { UpdatePiece(c, db) })
	router.DELETE("/pieces/:id", func(c *gin.Context) { DeletePiece(c, db) })
} 