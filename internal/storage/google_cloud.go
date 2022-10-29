package storage

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/plant_disease_detection/internal/credentials"
)

type GCloudStorage struct{}

func (s *GCloudStorage) _savePhoto(c *gin.Context, filename string) error {
	if credentials.GClient == nil {
		return fmt.Errorf("no credentials client defined")
	}
	form, _ := c.MultipartForm()
	files := form.File["uploads"]

	for _, file := range files {
		log.Println(file.Filename)

		opened_file, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return err
		}
		// err = credentials.GClient.UploadFile(opened_file, file.Filename)
		err = credentials.GClient.UploadFile(opened_file, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return err
		}
	}
	return nil
}

// must receive the user ID
func (s *GCloudStorage) SavePhoto(c *gin.Context) {
	filename := uuid.NewString()
	s._savePhoto(c, filename)
	credentials.MongoCl.InsertUserPhoto("16", filename)
	c.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("files uploaded!"),
	})
}
