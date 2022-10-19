package storage

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plant_disease_detection/internal/credentials"
)

type GCloudStorage struct{}

func (s *GCloudStorage) SavePhoto(c *gin.Context) {
	if credentials.GClient == nil {
		return
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
			return
		}
		err = credentials.GClient.UploadFile(opened_file, file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("%d files uploaded!: ", len(files)),
	})
}
