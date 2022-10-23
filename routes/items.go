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
func InitItemsRoutes() {
	router.POST("/items/create/:spaceid", createItem)
	router.POST("/items/update/:spaceid/:itemid", updateItem)
	router.DELETE("/items/delete/all/:spaceid", deleteAllItems)
	router.DELETE("/items/delete/:spaceid/:itemid", deleteItem)
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
	spaces.UpdateOne(
		context.Background(),
		filter,
		bson.D{{Key: "$push",
			Value: bson.D{
				{Key: "items", Value: newItem},
			},
		}},
		opts,
	)

	c.IndentedJSON(http.StatusOK, newItem)
}

func deleteAllItems(c *gin.Context) {
	opts := options.Update().SetUpsert(false)
	filter := bson.D{{Key: "id", Value: c.Param("spaceid")}}

	// insert space into mongodb
	response, err := spaces.UpdateOne(
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

func deleteItem(c *gin.Context) {
	spaceid := c.Param("spaceid")
	itemid := c.Param("itemid")

	filter := bson.D{{Key: "id", Value: spaceid}}

	upsert := false
	after := options.After
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var deletedItem = Item{
		ID: itemid,
	}
	// insert space into mongodb
	spaces.FindOneAndUpdate(
		context.Background(),
		filter,
		bson.D{{Key: "$pull",
			Value: bson.D{
				{Key: "items", Value: bson.D{{Key: "id", Value: itemid}}},
			},
		}},
		&opts,
	)

	c.IndentedJSON(http.StatusOK, deletedItem)
}

func updateItem(c *gin.Context) {
	spaceid := c.Param("spaceid")
	itemid := c.Param("itemid")

	opts := options.Update().SetUpsert(false)
	filter := bson.D{
		{Key: "id", Value: spaceid},
		{Key: "items.id", Value: itemid},
	}

	var updatedItem = Item{
		ID: itemid,
	}

	if err := c.Bind(&updatedItem); err != nil {
		return
	}

	// insert space into mongodb
	spaces.UpdateOne(
		context.Background(),
		filter,
		bson.D{{Key: "$set",
			Value: bson.D{
				{Key: "items.$.name", Value: updatedItem.Name},
				{Key: "items.$.complete", Value: updatedItem.Complete},
			},
		}},
		opts,
	)

	c.IndentedJSON(http.StatusOK, updatedItem)
}
