package user

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"project_p/backend/src/middleware"
)

func Login(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user User
		var passhash string
		var role string
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Corrupted Data"})
			return
		}
		if user.Email == "" || user.Password == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if !strings.Contains(user.Email, "@") || !strings.Contains(user.Email, ".") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}
		err := db.QueryRow(
			"SELECT Username, Email, Password, role FROM users WHERE email = ?", user.Email,
		).Scan(&user.Username, &user.Email, &passhash, &role)

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "email or password invalid"})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(passhash), []byte(user.Password))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect password"})
			return
		}
		token, err := middleware.GenerateToken(user.Email, role)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Login successfully",
			"token":   token,
			"user": gin.H{
				"username": user.Username,
				"email":    user.Email,
			},
		})

	}
}
