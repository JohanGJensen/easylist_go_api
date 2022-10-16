package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SpaceRequest struct {
	Name string `json:"name" form:"name"`
	User string `json:"user" form:"user"`
}

// initialize all spaces routes
func InitSpacesRoutes(router *gin.Engine) {
	router.GET("/spaces/all", getSpaces)
	router.GET("/spaces/space/:id", getSpace)
	router.POST("/spaces/create", createSpace)
	router.POST("/spaces/update/:id", updateSpace)
	router.DELETE("/spaces/delete/all", deleteAllSpaces)
	router.DELETE("/spaces/delete/:id", deleteSpace)
}

// GET:: call that gets all spaces documents in the spaces collection
func getSpaces(c *gin.Context) {
	cursor, err := spacesCollection.Find(context.Background(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	var spaces []Space
	if err = cursor.All(context.Background(), &spaces); err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, spaces)
}

// GET:: return collection item - space
func getSpace(c *gin.Context) {
	id := c.Param("id")
	filter := bson.D{{Key: "id", Value: id}}

	var result bson.M
	spacesCollection.FindOne(context.Background(), filter).Decode(&result)

	if result != nil {
		c.IndentedJSON(http.StatusCreated, result)
	} else {
		var message = Message{Msg: "could not find space"}
		c.IndentedJSON(http.StatusBadRequest, message)
	}
}

// POST:: create space
func createSpace(c *gin.Context) {
	var newSpace Space

	// set uuid and items slice
	newSpace = Space{
		ID:    uuid.New().String(),
		Items: make([]Item, 0),
	}

	// Call Bind to bind the received data to
	// newSpace.
	if err := c.Bind(&newSpace); err != nil {
		return
	}

	// insert space into mongodb
	spacesCollection.InsertOne(context.Background(), newSpace)

	c.IndentedJSON(http.StatusCreated, newSpace)
}

// POST:: update space
func updateSpace(c *gin.Context) {
	opts := options.Update().SetUpsert(false)
	id := c.Param("id")
	filter := bson.D{{Key: "id", Value: id}}

	var update Space

	// Call Bind to bind the received data to
	// newSpace.
	if err := c.Bind(&update); err != nil {
		return
	}
	// insert space into mongodb
	response, err := spacesCollection.UpdateOne(
		context.Background(),
		filter,
		bson.D{{Key: "$set",
			Value: bson.D{
				{Key: "name", Value: update.Name},
				{Key: "user", Value: update.User},
			},
		}},
		opts,
	)

	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, response)
}

// DELETE:: delete all spaces
func deleteAllSpaces(c *gin.Context) {
	spacesCollection.DeleteMany(context.Background(), bson.M{})

	var message = bson.M{"msg": "all spaces deleted"}
	c.IndentedJSON(http.StatusOK, message)
}

// DELETE:: delete space
func deleteSpace(c *gin.Context) {
	id := c.Param("id")
	filter := bson.D{{Key: "id", Value: id}}
	spacesCollection.DeleteOne(context.Background(), filter)

	var message = Message{Msg: "space deleted"}
	c.IndentedJSON(http.StatusOK, message)
}
