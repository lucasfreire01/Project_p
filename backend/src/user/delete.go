package user

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}
		result, err := db.Exec("DELETE FROM users WHERE id = ?", id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete user"})
			return
		}
		rowAffect, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify deletion"})
			return
		}
		if rowAffect == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
	}
}
