package main

import (
	"github.com/alexflint/go-arg"
	api "github.com/plant_disease_detection/internal/photos_api"
)

var args struct {
	Storage string `arg:"-s,--storage" help:"The storage where photos are going to be stored. default: 'google'"`
}

func main() {
	arg.MustParse(&args)
	s := api.NewServer(api.WithStorage(args.Storage), api.WithAddress(":8080"))
	s.Serve()
}
