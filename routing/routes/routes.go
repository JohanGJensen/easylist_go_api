package routes

import (
	"example/easylist-api/routing/health"
	"example/easylist-api/routing/items"
	"example/easylist-api/routing/spaces"
	"example/easylist-api/routing/users"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine = gin.Default()

type Message struct {
	Message string `json:"message" binding:"required"`
	Token   string `json:"token,omitempty"`
}

func Init() {
	// setting cors origin rules
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("HOST")}
	config.AllowHeaders = []string{"Origin", "Authorization"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	router.GET("/health", health.CheckHealth)

	// initialize all the routes
	users.InitUserRoutes(router)
	spaces.InitSpacesRoutes(router)
	items.InitItemsRoutes(router)

	router.Run()
}
