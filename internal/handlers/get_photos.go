package handlers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/plant_disease_detection/internal/auth"
	"github.com/plant_disease_detection/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	// "github.com/plant_disease_detection/internal/db"
	// "go.mongodb.org/mongo-driver/bson"
)

func GetPhotos(c *gin.Context) {
	client := db.MongoCl.Client
	collection := client.Database("photos").Collection("user_photos")

	user_claims, exists := c.Get("user")

	if !exists {
		log.Printf("user does not exist in the db, which means he has no pictures stored")
		return
	}

	user, ok := user_claims.(*auth.Claims)
	if !ok {
		log.Printf("the user in the context does not have the correct Claims shape")
		return
	}

	cursor, err := collection.Find(c, bson.M{"username": user.Username})

	if err != nil {
		log.Printf("there was an error looking for the records with the user's username: " + user.Username)
		return
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Printf(err.Error())
			return
		}
		fmt.Println(doc["username"])
		fmt.Println(doc["photo_url"])
	}

	if err := cursor.Err(); err != nil {
		log.Printf(err.Error())
		return
	}
}
