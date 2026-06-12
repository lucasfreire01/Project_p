package client
import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)
type client struct{
	ID       int    `json:"id"`
	Username string `json:"username"`
	Tel    string `json:"tel"`
	Parents string `json:"parents"`
}

func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context){
		var user client 
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := db.Exec(
			"INSERT INTO clients (username, tel, parents) VALUES (?, ?, ?)",
			user.Username, user.Tel, user.Parents,
		)
		if err != nil{
			c.JSON(http.StatusConflict, gin.H{"error": "Client already exists"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Client created successfully"})

	}
	
	
}