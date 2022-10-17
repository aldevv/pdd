package auth

import (
	"log"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/plant_disease_detection/internal/credentials"
)

func CreateUser(ctx *gin.Context) {
	params := (&auth.UserToCreate{}).
		Email("user@example.com").
		EmailVerified(false).
		PhoneNumber("+15555550100").
		Password("secretPassword").
		DisplayName("John Doe").
		PhotoURL("http://www.example.com/12345678/photo.png").
		Disabled(false)

	client := credentials.AuthCl
	u, err := client.CreateUser(ctx, params)
	if err != nil {
		// log.Fatalf("error creating user: %v\n", err)
		log.Printf("error creating user: %v\n", err)
		return
	}
	log.Printf("Successfully created user: %v\n", u)
}
