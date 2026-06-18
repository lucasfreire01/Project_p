package client

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Client struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Tel      string `json:"tel"`
	Parents  string `json:"parents"`
}

func onlyNumbers(value string) bool {
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var client Client
		if err := c.ShouldBindJSON(&client); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if client.Username == "" || client.Tel == "" || client.Parents == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if len(client.Username) < 3 || len(client.Tel) < 9 || len(client.Parents) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "characters not enough"})
			return
		}
		if !onlyNumbers(client.Tel) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "number(s) invalid"})
			return
		}
		_, err := db.Exec(
			"INSERT INTO clients (username, tel, parents) VALUES (?, ?, ?)",
			client.Username, client.Tel, client.Parents,
		)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Client already exists"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Client created successfully"})

	}

}
