package client

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdateClient(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
			return
		}

		var client Client
		if err := c.ShouldBindJSON(&client); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if client.Username == "" || client.Parents == "" || client.Tel == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "field(s) empty"})
			return
		}
		if len(client.Username) < 3 || len(client.Parents) < 3 || len(client.Tel) < 9 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "characters not enought"})
			return
		}
		if !onlyNumbers(client.Tel) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tel invalid(just numbers)"})
			return
		}

		result, err := db.Exec("UPDATE clients SET username = ?, tel =?, parents =? WHERE id =?",
			client.Username, client.Tel, client.Parents, id)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "filed to update client"})
			return
		}
		rowsEfect, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "filed to verify update"})
			return
		}
		if rowsEfect == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Update successfully"})
	}
}
