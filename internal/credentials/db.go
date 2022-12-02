package credentials

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient struct {
	Client *mongo.Client
}

func (c *MongoClient) InsertUserPhoto(uid string, filename string) {
	col := c.Client.Database("photos").Collection("user_photos")

	payload := bson.M{"uid": uid, "filename": filename}
	res, err := col.InsertOne(context.Background(), payload)
	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Println("==========")
	fmt.Println(id)

}

var MongoCl *MongoClient

func ConnectDB() {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	url_conn := os.Getenv("MONGO_URL")
	fmt.Print(url_conn)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url_conn))
	if err != nil {
		log.Fatal(err)
	}
	MongoCl = &MongoClient{Client: client}

	// // NOTE: needed?
	// defer client.Disconnect(ctx)
	err = MongoCl.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	collection := MongoCl.Client.Database("photos").Collection("user_photos")
	res, err := collection.InsertOne(context.Background(), bson.M{"bye": "moon"})
	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Println("==========")
	fmt.Println(id)
}
