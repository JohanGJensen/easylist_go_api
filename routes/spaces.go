package routes

import (
	"context"
	"example/easylist-api/auth"
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
func InitSpacesRoutes() {
	spaces := router.Group("/spaces").Use(auth.Auth())
	{
		spaces.GET("/all", getSpaces)
		spaces.GET("/space/:id", getSpace)
		spaces.POST("/create", createSpace)
		spaces.POST("/update/:id", updateSpace)
		spaces.DELETE("/delete/all", deleteAllSpaces)
		spaces.DELETE("/delete/:id", deleteSpace)
	}
}

// GET:: call that gets all spaces documents in the spaces collection
func getSpaces(c *gin.Context) {
	cursor, err := spaces.Find(context.Background(), bson.M{})

	if err != nil {
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "A problem occured with loading the spaces.",
		})

		return
	}

	spaces := make([]Space, 0)
	if err = cursor.All(context.Background(), &spaces); err != nil {
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "A problem occured with decoding the spaces.",
		})

		return
	}

	c.IndentedJSON(http.StatusOK, spaces)
}

// GET:: return collection item - space
func getSpace(c *gin.Context) {
	id := c.Param("id")
	filter := bson.D{{Key: "id", Value: id}}

	result := make(bson.M)
	spaces.FindOne(context.Background(), filter).Decode(&result)

	if result != nil {
		c.IndentedJSON(http.StatusOK, result)
	} else {
		message := Message{Message: "could not find space"}
		c.IndentedJSON(http.StatusBadRequest, message)
	}
}

// POST:: create space
func createSpace(c *gin.Context) {
	// set uuid and items slice
	newSpace := Space{
		ID:    uuid.New().String(),
		Items: make([]Item, 0),
	}

	// Call Bind to bind the received data to
	// newSpace.
	if err := c.Bind(&newSpace); err != nil {
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "failed to create a new space.",
		})

		return
	}

	// insert space into mongodb
	_, err := spaces.InsertOne(context.Background(), newSpace)

	if err != nil {
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "could not insert new space into the database correctly.",
		})

		return
	}

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
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "could not bind data properly to space.",
		})

		return
	}
	// insert space into mongodb
	response, err := spaces.UpdateOne(
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
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "failed to correctly update space in the database",
		})

		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

// DELETE:: delete all spaces
func deleteAllSpaces(c *gin.Context) {
	_, err := spaces.DeleteMany(context.Background(), bson.M{})

	if err != nil {
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "failed to correctly delete spaces.",
		})

		return
	}

	c.IndentedJSON(http.StatusOK, Message{
		Message: "all spaces deleted.",
	})
}

// DELETE:: delete space
func deleteSpace(c *gin.Context) {
	id := c.Param("id")
	filter := bson.D{{Key: "id", Value: id}}

	_, err := spaces.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "failed to delete space.",
		})

		return
	}

	c.IndentedJSON(http.StatusOK, Message{
		Message: "space deleted.",
	})
}
