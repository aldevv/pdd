package auth

import (
	"log"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/plant_disease_detection/internal/credentials"
	"google.golang.org/api/iterator"
)

var client *auth.Client = credentials.AuthCl

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func CreateUser(ctx *gin.Context) {
	var user User
	if err := ctx.BindJSON(&user); err != nil {
		return
	}
	log.Println(user.Email)

	params := (&auth.UserToCreate{}).
		Email(user.Email).
		EmailVerified(false).
		PhoneNumber(user.Phone).
		Password(user.Password).
		DisplayName(user.Name).
		PhotoURL("http://www.example.com/12345678/photo.png").
		Disabled(false)

	u, err := client.CreateUser(ctx, params)
	if err != nil {
		// log.Fatalf("error creating user: %v\n", err)
		log.Printf("error creating user: %v\n", err)
		return
	}
	log.Printf("Successfully created user: %v\n", u.DisplayName)
}

func GetUsers(ctx *gin.Context) {
	users := client.Users(ctx, "")
	for {
		user, err := users.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("error listing users: %s\n", err)
		}
		log.Printf("read user user: %v\n", user.DisplayName)
	}
}
func GetUser(ctx *gin.Context) {
	user, err := client.GetUser(ctx, "Ilwx07UPG6dJA8TOsr4qSh1lipy2")
	if err != nil {
		// log.Fatalf("error creating user: %v\n", err)
		log.Printf("error creating user: %v\n", err)
		return
	}

	log.Printf("%v", user.DisplayName)
}
