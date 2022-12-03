package photo_storage

import "github.com/gin-gonic/gin"

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
