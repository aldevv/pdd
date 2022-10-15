package photos_api

import "github.com/gin-gonic/gin"

type Storage interface {
	SavePhoto(c *gin.Context)
}

func GetStorage(storage_name interface{}) Storage {
	switch storage_name {
	case "local", nil, "":
		return &LocalStorage{}
	default:
		panic("Storage method not found")
	}
}
