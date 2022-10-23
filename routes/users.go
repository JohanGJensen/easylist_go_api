package routes

import (
	"context"
	"example/web-service-gin/mongodb"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var users *mongo.Collection = mongodb.GetCollection("users")

// initialize all item routes
func InitUserRoutes() {
	router.POST("users/register", RegisterUser)
	router.POST("users/login", LoginUser)
}

// METHOD: POST
// creates new user
func RegisterUser(c *gin.Context) {
	username := c.PostForm("username")
	pwd := c.PostForm("password")

	if username == "" {
		fmt.Println("no username provided!")
		return
	}

	user := FindUserInCollection(username)

	if user != (User{}) {
		fmt.Println("user already exists! Please select another username")
		return
	}

	if pwd == "" {
		fmt.Println("no password provided!")
		return
	}

	hash, _ := HashPassword(pwd)

	var newUser = User{
		ID:       uuid.New().String(),
		Username: username,
		Password: hash,
	}

	// insert user into mongodb
	users.InsertOne(context.Background(), newUser)

	c.IndentedJSON(http.StatusOK, newUser)
}

// METHOD: POST
// attempts to login user
func LoginUser(c *gin.Context) {
	username := c.PostForm("username")
	pwd := c.PostForm("password")

	if username == "" {
		fmt.Println("no username provided!")
		return
	}

	if pwd == "" {
		fmt.Println("no password provided!")
		return
	}

	user := FindUserInCollection(username)

	if user == (User{}) {
		fmt.Println("user does not exist! Please make sure the username is correct")
		return
	}

	match := CheckPasswordHash(pwd, user.Password)

	fmt.Println(match)
	fmt.Println(user.Password)
}

// UTILITY
// checks database collection for username
func FindUserInCollection(username string) User {
	filter := bson.D{{Key: "username", Value: username}}

	var result User
	users.FindOne(context.Background(), filter).Decode(&result)

	return result
}
