package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RoleMiddleware(requestRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != requestRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}

}
