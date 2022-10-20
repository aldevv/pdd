package credentials

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectDB() {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	url_conn := "mongodb://miusuario:pass@localhost:27017/photos?authSource=admin"
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url_conn))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("photos").Collection("user_photos")
	res, err := collection.InsertOne(context.Background(), bson.M{"hello": "world"})
	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Println("==========")
	fmt.Println(id)
}
