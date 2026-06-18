package user

import (
	"database/sql"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		hashUpper := false
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if user.Username == "" || user.Email == "" || user.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if len(user.Username) < 3 || len(user.Password) < 12 || len(user.Email) < 15 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "characters not enough"})
			return
		}
		if !strings.Contains(user.Email, "@") || !strings.Contains(user.Email, ".") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}
		for _, char := range user.Password {
			if unicode.IsUpper(char) {
				hashUpper = true
				break
			}
		}
		if !hashUpper {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one uppercase letter"})
			return
		}
		if user.Role == "" {
			user.Role = "employee"
		}
		if user.Role != "admin" && user.Role != "employee" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
			return

		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		_, err = db.Exec(
			"INSERT INTO users (username, email, password,role) VALUES (?, ?, ?, ?)",
			user.Username, user.Email, string(hashedPassword), user.Role,
		)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})

	}

}
