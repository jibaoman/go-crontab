package main


import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

func main() {
	var (
		client *mongo.Client
		err error
		database *mongo.Database
		collection *mongo.Collection
	)
	if client,err = mongo.Connect(context.TODO(),"mongodb://172.29.3.100:27017",clientopt.ConnectTimeout(5 * time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	database = client.Database("my_db")

	collection = database.Collection("my_collection")

	collection = collection

}

