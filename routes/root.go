package routes

import (
	"example/web-service-gin/mongodb"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var router *gin.Engine = gin.Default()
var spaces *mongo.Collection = mongodb.GetCollection("spaces")

type User struct {
	ID       string `bson:"id" json:"id" form:"id"`
	Username string `bson:"username" json:"username" form:"username"`
	Password string `bson:"password" json:"password" form:"password"`
}

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
	// router := gin.Default()

	// setting cors origin rules
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("HOST")}
	router.Use(cors.New(config))

	router.GET("/health", CheckHealth)

	// initialize all the routes
	InitUserRoutes()
	InitSpacesRoutes()
	InitItemsRoutes()

	router.Run()
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
