package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/plant_disease_detection/internal/db"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/bson"
)

type UserPassRecovery struct {
	Email string `json:"email" binding:"required"`
}

var client *sendgrid.Client = sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

func PasswordRecovery(c *gin.Context) {
	var requestBody UserPassRecovery
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists in the database
	client := db.MongoCl.Client
	collection := client.Database("photos").Collection("users")
	var user User
	err := collection.FindOne(c, bson.M{"email": requestBody.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	tokenString := create_jwt_token(c, user.Username, user.Email)
	if tokenString == "" {
		c.JSON(500, gin.H{"error": "error creating token"})
		return
	}

	resetLink := "http://example.com/reset-password?token=" + tokenString
	fmt.Println("username: ", user.Username)
	fmt.Println("email: ", user.Email)
	fmt.Println("token: ", tokenString)
	err = SendResetEmail(user.Email, resetLink)
	fmt.Println("error: ", err)
	if err != nil {
		c.JSON(500, gin.H{"error": "error sending email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password recovery email sent"})
}

func SendResetEmail(email string, link string) error {
	from := mail.NewEmail("Plant Disease Detection", os.Getenv("PDD_EMAIL"))
	subject := "Password Reset"
	to := mail.NewEmail("Plant", email)

	plainTextContent := fmt.Sprintf("Click on this link to reset your password: %s", link)
	htmlContent := fmt.Sprintf("<p>Click on this <a href=\"%s\">link</a> to reset your password</p>", link)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	res, err := client.Send(message)
	log.Println(res)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func PasswordReset(c *gin.Context) {
	// Get the user ID from the JWT token
	user_claims, exists := c.Get("user")
	if !exists {
		log.Printf("user does not exist in the db, which means he has no pictures stored")
		return
	}
	// Get the new password from the request body
	var requestBody struct {
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the new password
	hashedPassword := HashPassword(requestBody.NewPassword)

	// Update the user's password in the database
	client := db.MongoCl.Client
	collection := client.Database("photos").Collection("users")

	user := user_claims.(*Claims)
	filter := bson.M{"username": user.Username}

	update := bson.M{"$set": bson.M{"password": hashedPassword}}
	result, err := collection.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}
	if result.ModifiedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	c.Status(http.StatusOK)
}
