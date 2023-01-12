package photos_api

import (
	"github.com/gin-gonic/gin"
	"github.com/plant_disease_detection/internal/auth"
	"github.com/plant_disease_detection/internal/db"
	"github.com/plant_disease_detection/internal/middleware"
	"github.com/plant_disease_detection/internal/photo_storage"
	"github.com/plant_disease_detection/internal/routes"
)

var addr string

func setAddr(addrs ...string) {
	if len(addrs) == 0 {
		addr = "localhost:8080"
	} else {
		addr = addrs[0]
	}
}

func Addr() string {
	return addr
}

func Serve(storage_name interface{}, addrs ...string) {
	db.ConnectDB()
	storage := photo_storage.GetStorage(storage_name)

	router := gin.Default()
	router.MaxMultipartMemory = 100 << 20 // 8 MiB

	router.POST("/create_user", auth.CreateUser)
	router.POST("/login", auth.Login)

	private := router.Group("/api", middleware.Protect)

	// TODO: move savephoto to routes
	private.POST("/upload", storage.SavePhoto)
	private.POST("/get_photos", routes.GetPhotos)
	// private.GET("/get_user", auth.GetUser)

	setAddr(addrs...)
	router.Run(Addr())
}
