package main

import (
	"github.com/plant_disease_detection/internal/photos_api"
)

var args struct {
	Storage string `arg:"-s,--storage" help:"The storage where photos are going to be stored. default: 'local'"`
}

func main() {
	photos_api.Main()
}
