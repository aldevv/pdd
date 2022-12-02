package photos_api

import (
	"github.com/gin-gonic/gin"
	"github.com/plant_disease_detection/internal/auth"
	"github.com/plant_disease_detection/internal/credentials"
	"github.com/plant_disease_detection/internal/storage"
)

var addr string

func setAddr(addrs ...string) {
	addr = "localhost:8080"
	if len(addrs) > 0 {
		addr = addrs[0]
	}
}

func Addr() string {
	return addr
}

func Serve(storage_name interface{}, addrs ...string) {
	credentials.ConnectDB()

	router := gin.Default()

	router.MaxMultipartMemory = 100 << 20 // 8 MiB

	storage := storage.GetStorage(storage_name)
	router.POST("/upload", storage.SavePhoto)
	router.POST("/create_user", auth.CreateUser)
	// router.POST("/create_user", auth.CreateUser)
	// router.GET("/get_user", auth.GetUser)
	// router.GET("/get_users", auth.GetUsers)

	setAddr(addrs...)
	router.Run(Addr())
}
