package photos_api

import (
	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"
	"github.com/plant_disease_detection/internal/auth"
	"github.com/plant_disease_detection/internal/storage"
)

var args struct {
	Storage string `arg:"-s,--storage" help:"The storage where photos are going to be stored. default: 'local'"`
}

func Main() {
	arg.MustParse(&args)
	storage := storage.GetStorage(args.Storage)

	router := gin.Default()
	router.MaxMultipartMemory = 100 << 20 // 8 MiB
	router.POST("/upload", storage.SavePhoto)
	router.POST("/create_user", auth.CreateUser)
	router.Run(":8080")
}
