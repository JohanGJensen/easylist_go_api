package items

import (
	"context"
	"example/easylist-api/auth"
	"example/easylist-api/routing/spaces"
	"example/easylist-api/structs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message structs.Message

type Space spaces.Space
type Item spaces.Item

// initialize all item routes
func InitItemsRoutes(e *gin.Engine) {
	items := e.Group("/items").Use(auth.Auth())
	{
		items.POST("/create/:spaceid", createItem)
		items.POST("/update/:spaceid/:itemid", updateItem)
		items.DELETE("/delete/all/:spaceid", deleteAllItems)
		items.DELETE("/delete/:spaceid/:itemid", deleteItem)
	}
}

func createItem(c *gin.Context) {
	opts := options.Update().SetUpsert(false)
	filter := bson.D{{Key: "id", Value: c.Param("spaceid")}}

	newItem := Item{
		ID:       uuid.New().String(),
		Complete: false,
	}

	// Call Bind to bind the received data to
	// newItem.
	if err := c.Bind(&newItem); err != nil {
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "failed to bind item request parameters",
		})

		return
	}

	// insert space into mongodb
	_, err := spaces.SCollection.UpdateOne(
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
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "failed to create item on space",
		})

		return
	}

	c.IndentedJSON(http.StatusOK, newItem)
}

func deleteAllItems(c *gin.Context) {
	opts := options.Update().SetUpsert(false)
	filter := bson.D{{Key: "id", Value: c.Param("spaceid")}}

	// insert space into mongodb
	response, err := spaces.SCollection.UpdateOne(
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
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "failed to delete items on space",
		})

		return
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

	// insert space into mongodb
	err := spaces.FindOneAndUpdate(
		context.Background(),
		filter,
		bson.D{{Key: "$pull",
			Value: bson.D{
				{Key: "items", Value: bson.D{{Key: "id", Value: itemid}}},
			},
		}},
		&opts,
	)

	if err != nil {
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "failed to delete item on space",
		})

		return
	}

	c.IndentedJSON(http.StatusOK, Message{
		Message: "item deleted",
	})
}

func updateItem(c *gin.Context) {
	spaceid := c.Param("spaceid")
	itemid := c.Param("itemid")

	opts := options.Update().SetUpsert(false)
	filter := bson.D{
		{Key: "id", Value: spaceid},
		{Key: "items.id", Value: itemid},
	}

	updatedItem := Item{
		ID: itemid,
	}

	if err := c.Bind(&updatedItem); err != nil {
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "failed to bind request params onto item.",
		})

		return
	}

	// insert space into mongodb
	_, err := spaces.SCollection.UpdateOne(
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

	if err != nil {
		log.Print(err)

		c.IndentedJSON(http.StatusBadRequest, Message{
			Message: "failed to update item.",
		})

		return
	}

	c.IndentedJSON(http.StatusOK, updatedItem)
}
