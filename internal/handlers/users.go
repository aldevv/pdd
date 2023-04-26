package handlers

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/plant_disease_detection/internal/auth"
	"github.com/plant_disease_detection/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UrlDecode(input string) (string, error) {
	decoded, err := url.QueryUnescape(input)
	if err != nil {
		return "", err
	}
	return decoded, nil
}

func GetUser(ctx *gin.Context) {
	username, err := UrlDecode(ctx.Param("username"))
	if err != nil {
		ctx.Status(500)
		return
	}
	col := db.GetCollection("users")
	var user auth.User
	res := col.FindOne(ctx, bson.M{"username": username}, options.FindOne().SetProjection(bson.M{"password": 0}))
	if res.Err() != nil {
		ctx.Status(400)
		return
	}

	err = res.Decode(&user)
	if err != nil {
		ctx.Status(500)
		return
	}

	ctx.JSON(200, gin.H{"data": user})
}
