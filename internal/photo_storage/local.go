package photo_storage

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/plant_disease_detection/internal/auth"
)

type LocalStorage struct{}

func (s *LocalStorage) SavePhoto(c *gin.Context) {

	form, _ := c.MultipartForm()
	files := form.File["uploads"]

	claims, _ := c.Get("user")
	user, _ := claims.(*auth.Claims)

	for _, file := range files {
		log.Println(file.Filename)

		filepath := "./uploads/" + uuid.New().String() + filepath.Ext(file.Filename)
		c.SaveUploadedFile(file, filepath)
		SaveInDb(filepath, user.Username)
	}

	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}
