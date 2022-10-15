package main

import (
	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"
	photos_api "github.com/plant_disease_detection/internal/photos_api"
)

var args struct {
	Storage string `arg:"-s,--storage" help:"The storage where photos are going to be stored. default: 'local'"`
}

func main() {
	arg.MustParse(&args)
	storage := photos_api.GetStorage(args.Storage)

	router := gin.Default()
	router.MaxMultipartMemory = 100 << 20 // 8 MiB
	router.POST("/upload", storage.SavePhoto)
	router.Run(":8080")
}
