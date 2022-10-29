package credentials

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
}

func (c *MongoClient) InsertUserPhoto(uid string, filename string) {
	col := c.client.Database("photos").Collection("user_photos")

	payload := bson.M{"uid": uid, "filename": filename}
	res, err := col.InsertOne(context.Background(), payload)
	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Println("==========")
	fmt.Println(id)

}

var MongoCl *MongoClient = &MongoClient{}

func ConnectDB() {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	url_conn := "mongodb://miusuario:pass@localhost:27017/photos?authSource=admin"
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url_conn))
	if err != nil {
		log.Fatal(err)
	}
	MongoCl.client = client
	// // NOTE: needed?
	// // defer client.Disconnect(ctx)
	// err = client.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// collection := client.Database("photos").Collection("user_photos")
	// res, err := collection.InsertOne(context.Background(), bson.M{"hello": "world"})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// id := res.InsertedID
	// fmt.Println("==========")
	// fmt.Println(id)
}
