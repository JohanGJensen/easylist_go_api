package routes

import (
	"easylist/routing/health"
	"easylist/routing/items"
	"easylist/routing/spaces"
	"easylist/routing/users"
	util "easylist/utility"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine = gin.Default()

// type User struct {
// 	ID       string `bson:"id" json:"id" form:"id"`
// 	Username string `bson:"username" json:"username" form:"username"`
// 	Password string `bson:"password" json:"password" form:"password"`
// }

// type Space struct {
// 	ID    string `json:"id" form:"id"`
// 	Items []Item `json:"items" form:"items"`
// 	Name  string `json:"name" form:"name"`
// 	User  string `json:"user" form:"user"`
// }

// type Item struct {
// 	ID       string `json:"id" form:"id"`
// 	Name     string `json:"name" form:"name"`
// 	Complete bool   `json:"complete" form:"complete"`
// }

type Message struct {
	Message string `json:"message" binding:"required"`
	Token   string `json:"token,omitempty"`
}

func Init() {
	env, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// setting cors origin rules
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{env.Host}
	router.Use(cors.New(config))

	router.GET("/health", health.CheckHealth)

	// initialize all the routes
	users.InitUserRoutes(router)
	spaces.InitSpacesRoutes(router)
	items.InitItemsRoutes(router)

	router.Run()
}

// func HashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return string(bytes), err
// }

// func CheckPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }
