package db

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

func GetCollection(col string) *mongo.Collection {
	client := MongoCl.Client
	return client.Database("photos").Collection(col)
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

func createUserIndex(db *mongo.Database, ctx context.Context) {
	model := mongo.IndexModel{
		Keys: bson.M{
			"username": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	db.Collection("users").Indexes().CreateOne(ctx, model)

	model = mongo.IndexModel{
		Keys: bson.M{
			"email": 1,
		},
		Options: options.Index().SetUnique(true),
	}
	db.Collection("users").Indexes().CreateOne(ctx, model)
}

func ConnectDB() {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	url_conn := os.Getenv("MONGO_URL")
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url_conn))
	if err != nil {
		log.Fatal(err)
	}
	MongoCl = &MongoClient{Client: client}

	err = MongoCl.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	db := MongoCl.Client.Database("photos")
	createUserIndex(db, ctx)
}
