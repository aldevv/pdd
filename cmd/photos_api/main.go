package main

import (
	"github.com/alexflint/go-arg"
	"github.com/plant_disease_detection/internal/photos_api"
)

var args struct {
	Storage string `arg:"-s,--storage" help:"The storage where photos are going to be stored. default: 'local'"`
}

func main() {
	arg.MustParse(&args)
	photos_api.Serve(args.Storage, ":8080")
}
