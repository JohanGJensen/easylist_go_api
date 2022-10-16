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

// initialize all item routes
func InitItemsRoutes(router *gin.Engine) {
	router.POST("/items/create/:spaceid", createItem)
	router.DELETE("/items/delete/all/:spaceid", deleteAllItems)
}

func createItem(c *gin.Context) {
	opts := options.Update().SetUpsert(false)
	filter := bson.D{{Key: "id", Value: c.Param("spaceid")}}

	var newItem = Item{
		ID:       uuid.New().String(),
		Complete: false,
	}

	// Call Bind to bind the received data to
	// newItem.
	if err := c.Bind(&newItem); err != nil {
		return
	}

	// insert space into mongodb
	response, err := spacesCollection.UpdateOne(
		context.Background(),
		filter,
		bson.D{{Key: "$push",
			Value: bson.D{
				{Key: "items", Value: newItem},
			},
		}},
		opts,
	)

	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, response)
}

func deleteAllItems(c *gin.Context) {
	opts := options.Update().SetUpsert(false)
	filter := bson.D{{Key: "id", Value: c.Param("spaceid")}}

	// insert space into mongodb
	response, err := spacesCollection.UpdateOne(
		context.Background(),
		filter,
		bson.D{{Key: "$set",
			Value: bson.D{
				{Key: "items", Value: make([]Item, 0)},
			},
		}},
		opts,
	)

	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, response)
}
