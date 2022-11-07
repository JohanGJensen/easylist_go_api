package routes

import (
	"context"
	"example/easylist-api/auth"
	"example/easylist-api/mongodb"
	"example/easylist-api/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var users *mongo.Collection = mongodb.GetCollection("users")

type UserRequest struct {
	Username string `bson:"username" json:"username" form:"username" binding:"required,min=3,max=16"`
	Password string `bson:"password" json:"password" form:"password" binding:"required,min=3"`
}

// initialize all item routes
func InitUserRoutes() {
	users := router.Group("/users")
	{
		users.POST("/register", registerUser)
		users.POST("/login", loginUser)
	}
}

// METHOD: POST
func registerUser(c *gin.Context) {
	body := UserRequest{}
	// handle validation errors
	if er := c.ShouldBind(&body); er != nil {
		validation.Validate(c, er)
		return
	}

	if body.Username == "" {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "no username provided!",
		})
	}

	user := FindUserInCollection(body.Username)

	if user != (User{}) {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "user already exists! Please select another username",
		})
	}

	if body.Password == "" {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "no password provided!",
		})
	}

	hash, _ := HashPassword(body.Password)

	newUser := User{
		ID:       uuid.New().String(),
		Username: body.Username,
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
func loginUser(c *gin.Context) {
	body := UserRequest{}
	// handle validation errors
	if er := c.ShouldBind(&body); er != nil {
		validation.Validate(c, er)
		return
	}

	if body.Username == "" {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "no username provided!",
		})
	}

	if body.Password == "" {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "no password provided!",
		})
	}

	user := FindUserInCollection(body.Username)

	if user == (User{}) {
		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "user does not exist! Please make sure the username is correct",
		})
	}

	match := CheckPasswordHash(body.Password, user.Password)

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
