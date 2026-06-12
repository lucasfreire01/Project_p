package user

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *sql.DB) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var user User
		var passhash string
		if err := ctx.ShouldBindJSON(&user); err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Corrupted Data"})
			return
		}
		err := db.QueryRow(
			"SELECT Username, Email, Password FROM users WHERE email = ?", user.Email,
		).Scan(&user.Username, &user.Email, &passhash)
		
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusBadRequest, gin.H{"error" : "email or password invalid"})
			return
		}
		if err != nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		
		err = bcrypt.CompareHashAndPassword([]byte(passhash), []byte(user.Password))
		if err != nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect error"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Login successfully",
			"user": gin.H{
				"username": user.Username,
				"email": user.Email,
			},
		})

	}
}
