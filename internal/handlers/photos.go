package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/plant_disease_detection/internal/auth"
	"github.com/plant_disease_detection/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "github.com/plant_disease_detection/internal/db"
	// "go.mongodb.org/mongo-driver/bson"
)

type userPhoto struct {
	Id        string `json:"id" bson:"_id" bson:"_id"`
	Username  string `json:"username" bson:"username"`
	Email     string `json:"email" bson:"email"`
	Photo_url string `json:"photo_url" bson:"photo_url"`
	Sickness  string `json:"sickness,omitempty" bson:"sickness,omitempty"`
	Accuracy  string `json:"accuracy,omitempty" bson:"accuracy,omitempty"`
}

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

	var userPhotos []userPhoto
	if err = cursor.All(c, &userPhotos); err != nil {
		log.Printf(err.Error())
		return
	}
	c.JSON(200, gin.H{"data": userPhotos})
}

func GetPhoto(c *gin.Context) {
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

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println("could not convert id given to objectID")
		return
	}
	var record userPhoto
	err = collection.FindOne(c, bson.M{"_id": id, "username": user.Username}).Decode(&record)

	if err != nil {
		log.Printf("there was an error decoding the record for user with id %s", id)
		return
	}

	c.JSON(200, gin.H{"data": record})
}
