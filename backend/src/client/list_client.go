package client

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListClient(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
			return
		}

		var client Client
		err = db.QueryRow(
			"SELECT id, username, tel, parents FROM clients WHERE id = ?",
			id,
		).Scan(&client.ID, &client.Username, &client.Tel, &client.Parents)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "filed to get client"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "client listed"})
	}
}
