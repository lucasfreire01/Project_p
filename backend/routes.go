package backend

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	_ "modernc.org/sqlite"
	"project_p/backend/src/client"
	"project_p/backend/src/middleware"
	"project_p/backend/src/user"
)

func Routes() *gin.Engine {

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddle())
	router.Use(middleware.RequestLogger())
	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'employee'
	);
`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS clients(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		tel TEXT NOT NULL UNIQUE,
		parents TEXT NOT NULL
	);
`)
	if err != nil {
		log.Fatal(err)
	}

	router.POST("/register", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), user.Register(db))
	router.POST("/login", user.Login(db))
	router.POST("/register_client", middleware.AuthMiddleware(), client.Register(db))
	router.DELETE("/client/:id", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), client.DeleteClient(db))
	router.PUT("/client/:id", middleware.AuthMiddleware(), client.UpdateClient(db))
	router.GET("/client/:id", middleware.AuthMiddleware(), client.ListClient(db))
	router.PUT("/update/:id", middleware.AuthMiddleware(), user.UpdateUser(db))
	router.Run(":8080")
	return router
}
