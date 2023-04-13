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
	"github.com/plant_disease_detection/internal/handlers"
	"go.mongodb.org/mongo-driver/bson"
)

type GCloudStorage struct{}

func (s *GCloudStorage) _savePhoto(c *gin.Context) error {
	if credentials.GClient == nil {
		return fmt.Errorf("no credentials client defined")
	}
	form, _ := c.MultipartForm()
	files := form.File["uploads"]

	claims, _ := c.Get("user")
	user, _ := claims.(*auth.Claims)

	var fp string
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
		err = credentials.GClient.UploadFile(opened_file, filename)
		SaveInDb(filename, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return err
		}
		err = handlers.SendAI(c, filename)
		if err != nil {
			log.Printf("failed to send the photoURL to sqs queue")
			log.Print(err)
		}
	}
	c.JSON(http.StatusOK, bson.M{"photo_url": fp})
	return nil
}

// must receive the user ID
func (s *GCloudStorage) SavePhoto(c *gin.Context) {
	s._savePhoto(c)
}
