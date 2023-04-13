package photo_storage

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/plant_disease_detection/internal/auth"
	"github.com/plant_disease_detection/internal/handlers"
	"go.mongodb.org/mongo-driver/bson"
)

type LocalStorage struct{}

func (s *LocalStorage) SavePhoto(c *gin.Context) {

	form, _ := c.MultipartForm()
	files := form.File["uploads"]

	claims, _ := c.Get("user")
	user, _ := claims.(*auth.Claims)

	dir := "./uploads"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	var fp string
	for _, file := range files {
		log.Println(file.Filename)

		filepath := dir + "/" + uuid.New().String() + filepath.Ext(file.Filename)
		fp = filepath

		c.SaveUploadedFile(file, filepath)
		SaveInDb(filepath, user.Username)

		err := handlers.SendAI(c, filepath)
		if err != nil {
			log.Printf("failed to send the photoURL to sqs queue")
			log.Print(err)
		}
	}

	c.JSON(http.StatusOK, bson.M{"photo_url": fp})
}
