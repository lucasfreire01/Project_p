package user

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		targetID, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		emailValue, exist := c.Get("email")
		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		currentEmail, ok := emailValue.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token email"})
			return
		}

		roleValue, exist := c.Get("role")
		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		currentRole, ok := roleValue.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token role"})
			return
		}

		var user User
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

		hasUp := false
		for _, char := range user.Password {
			if unicode.IsUpper(char) {
				hasUp = true
				break
			}
		}
		if !hasUp {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one uppercase letter"})
			return
		}

		if currentRole != "admin" {
			var currentUserID int
			err := db.QueryRow(
				"SELECT id FROM users WHERE email = ?",
				currentEmail,
			).Scan(&currentUserID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
				return
			}
			if currentUserID != targetID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
				return
			}
			user.Role = "employee"
		}

		if currentRole == "admin" {
			if user.Role == "" {
				user.Role = "employee"
			}
			if user.Role != "admin" && user.Role != "employee" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
				return
			}
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		result, err := db.Exec(
			"UPDATE users SET username = ?, email = ?, password = ?, role = ? WHERE id = ?",
			user.Username, user.Email, string(hashPassword), user.Role, targetID,
		)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "failed to update user"})
			return
		}

		rowAfect, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify update"})
			return
		}
		if rowAfect == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user updated"})
	}
}
