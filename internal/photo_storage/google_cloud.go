package photo_storage

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/plant_disease_detection/internal/auth"
	"github.com/plant_disease_detection/internal/credentials"
	"github.com/plant_disease_detection/internal/db"
	"github.com/plant_disease_detection/internal/handlers"
	"go.mongodb.org/mongo-driver/bson"
)

type GCloudStorage struct{}

func (s *GCloudStorage) _savePhoto(c *gin.Context) error {
	fmt.Println("trace _savePhoto google storage")
	if credentials.GClient == nil {
		return fmt.Errorf("no credentials client defined")
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println("error in multipart form")
		return err
	}
	fmt.Println("before claims")
	files := form.File["uploads"]

	claims, _ := c.Get("user")
	user, _ := claims.(*auth.Claims)

	var fp string
	fmt.Println("before for")
	for _, file := range files {

		opened_file, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return err
		}

		filename := uuid.New().String() + filepath.Ext(file.Filename)
		fp = filename
		fmt.Println("trace _savePhoto uploading file to google storage")
		err = credentials.GClient.UploadFile(opened_file, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return err
		}

		fmt.Println("trace _savePhoto saving in db")
		SaveInDb(filename, user.Username)
		err = handlers.SendAI(c, filename)
		if err != nil {
			log.Printf("failed to send the photoURL to sqs queue")
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return err
		}
	}
	c.JSON(http.StatusOK, bson.M{"photo_url": fp})
	return nil
}

// must receive the user ID
func (s *GCloudStorage) SavePhoto(c *gin.Context) {
	err := s._savePhoto(c)
	if err != nil {
		fmt.Println("error:")
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
	}
}

func (s *GCloudStorage) DeletePhoto(c *gin.Context) {
	client := db.MongoCl.Client
	collection := client.Database("photos").Collection("user_photos")

	user_claims, exists := c.Get("user")

	if !exists {
		log.Printf("user does not exist in the db, which means he has no pictures stored")
		return
	}

	_, ok := user_claims.(*auth.Claims)
	if !ok {
		log.Printf("the user in the context does not have the correct Claims shape")
		return
	}

	id := c.Param("id")
	_, err := collection.DeleteOne(c, bson.M{"photo_url": id})

	if err != nil {
		log.Printf("there was an error deleting the record for user with id %s", id)
		return
	}

	credentials.GClient.RemoveFile(id)

	c.Status(200)
}
