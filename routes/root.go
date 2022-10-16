package routes

import (
	"example/web-service-gin/mongodb"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var spacesCollection *mongo.Collection = mongodb.GetCollection(mongodb.DB, "spacesCollection")

type Space struct {
	ID    string `json:"id" form:"id"`
	Items []Item `json:"items" form:"items"`
	Name  string `json:"name" form:"name"`
	User  string `json:"user" form:"user"`
}

type Item struct {
	ID       string `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Complete bool   `json:"complete" form:"complete"`
}

type Message struct {
	Msg string `json:"msg"`
}

func Init() {
	router := gin.Default()
	router.GET("/health", CheckHealth)

	// initialize all the routes
	InitSpacesRoutes(router)
	InitItemsRoutes(router)

	router.Run(os.Getenv("HOST"))
}
