package routes

import (
	"context"
	"example/easylist-api/auth"
	"example/easylist-api/mongodb"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var users *mongo.Collection = mongodb.GetCollection("users")

// initialize all item routes
func InitUserRoutes() {
	users := router.Group("/users")
	{
		users.POST("/register", RegisterUser)
		users.POST("/login", LoginUser)
	}
}

// METHOD: POST
func RegisterUser(c *gin.Context) {
	username := c.PostForm("username")
	pwd := c.PostForm("password")

	if username == "" {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "no username provided!",
		})
	}

	user := FindUserInCollection(username)

	if user != (User{}) {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "user already exists! Please select another username",
		})
	}

	if pwd == "" {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "no password provided!",
		})
	}

	hash, _ := HashPassword(pwd)

	newUser := User{
		ID:       uuid.New().String(),
		Username: username,
		Password: hash,
	}

	// insert user into mongodb
	_, err := users.InsertOne(context.Background(), newUser)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "Something went wrong with reqistering the user profile.",
		})
	}

	JWT, err := auth.GenerateJWT(user.Username)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "There was an error with generating the JWT!",
		})
	}

	c.IndentedJSON(http.StatusOK, Message{
		Message: "successfully registered.",
		Token:   JWT,
	})
}

// METHOD: POST
func LoginUser(c *gin.Context) {
	username := c.PostForm("username")
	pwd := c.PostForm("password")

	if username == "" {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "no username provided!",
		})
	}

	if pwd == "" {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "no password provided!",
		})
	}

	user := FindUserInCollection(username)

	if user == (User{}) {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "user does not exist! Please make sure the username is correct",
		})
	}

	match := CheckPasswordHash(pwd, user.Password)

	if match {
		JWT, err := auth.GenerateJWT(user.Username)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, Message{
				Message: "There was an error with generating the JWT!",
			})
		}

		c.IndentedJSON(http.StatusOK, Message{
			Message: "successfully logged in.",
			Token:   JWT,
		})
	}
}

// UTILITY
// checks database collection for username
func FindUserInCollection(username string) User {
	filter := bson.D{{Key: "username", Value: username}}

	var result User
	users.FindOne(context.Background(), filter).Decode(&result)

	return result
}
