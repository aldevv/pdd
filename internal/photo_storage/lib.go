package photo_storage

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/plant_disease_detection/internal/db"
	"go.mongodb.org/mongo-driver/bson"
)

type Storage interface {
	SavePhoto(c *gin.Context)
}

func GetStorage(storage_name interface{}) Storage {
	switch storage_name {
	case "local", nil, "":
		return &LocalStorage{}
	case "google", "google-cloud", "cloud", "gcloud", "g":
		return &GCloudStorage{}
	default:
		panic("Storage method not found")
	}
}

func SaveInDb(photo_url string, username string) {
	client := db.MongoCl.Client
	collection := client.Database("photos").Collection("user_photos")

	_, err := collection.InsertOne(context.Background(), bson.M{"username": username, "photo_url": photo_url})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted photo_url successfully")
}
