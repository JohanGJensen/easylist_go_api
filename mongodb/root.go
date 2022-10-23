package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURI string = ""

func SetURI(login string, pwd string, cluster string) {
	mongoURI = fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", login, pwd, cluster)
}

func ConnectDB() *mongo.Client {
	godotenv.Load(".env")

	login := os.Getenv("LOGIN")
	password := os.Getenv("PASSWORD")
	cluster := os.Getenv("CLUSTER")

	SetURI(login, password, cluster)

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

// getting database collections
func GetCollection(collectionName string) *mongo.Collection {
	collection := DB.Database("listDB").Collection(collectionName)
	return collection
}
