package backend


import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"project_p/backend/src/user"
	_ "modernc.org/sqlite"
	"log"

)

func Routes() *gin.Engine {

router := gin.Default()
db, err := sql.Open("sqlite", "./database.db")
if err != nil {
	log.Fatal(err)
}
_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
`)
if err != nil {
	log.Fatal(err)
}

router.POST("/register", user.Register(db))
router.POST("/login", user.Login(db))
router.Run(":8080")
return router
}