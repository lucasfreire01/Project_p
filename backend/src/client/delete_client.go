package client
import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteClient(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client id"})
			return
		}
		result, err := db.Exec("DELETE FROM clients WHERE id = ?", id)
		if err !=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete client"})
			return
		}
		rowAffect, err := result.RowsAffected()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "filed to verify deletion"})
			return
		}
		if rowAffect == 0{
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "client deleted"})
	}
}