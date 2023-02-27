package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/plant_disease_detection/internal/db"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Create a struct to read the username and password from the request body
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Create the JWT key used to create the signature
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func HashPassword(password string) string {
	pw := []byte(password)
	// NOTE: the hashed password contains the salt, so password's with same text, won't be stored the same
	result, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(result)
}

func ComparePassword(hashPassword string, password string) error {
	pw := []byte(password)
	hw := []byte(hashPassword)
	err := bcrypt.CompareHashAndPassword(hw, pw)
	return err
}

func create_jwt_token(c *gin.Context, username string) string {
	claims := &Claims{
		Username:         username,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		c.String(http.StatusInternalServerError, fmt.Sprintf("There was an error creating the JWT token"))
		return ""
	}
	return tokenString
}

func CreateUser(c *gin.Context) {
	// return's a JWT token
	client := db.MongoCl.Client
	collection := client.Database("photos").Collection("users")
	var user User
	if err := c.BindJSON(&user); err != nil {
		log.Println(err.Error())
		return
	}

	username := user.Username
	password := user.Password
	email := user.Email
	record := bson.M{"username": username, "password": HashPassword(password), "email": email}
	_, err := collection.InsertOne(c, record)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{})
		return
	}
	tokenString := create_jwt_token(c, user.Username)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Login(c *gin.Context) {
	var requestBody Credentials
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println(err.Error())
		return
	}

	client := db.MongoCl.Client
	collection := client.Database("photos").Collection("users")

	var credentials Credentials
	collection.FindOne(c, bson.M{"username": requestBody.Username}).Decode(&credentials)
	err := ComparePassword(credentials.Password, requestBody.Password)
	if err != nil {
		fmt.Println("wrong password")
		return

	}
	tokenString := create_jwt_token(c, credentials.Username)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
