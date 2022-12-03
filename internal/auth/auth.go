package auth

import (
	"context"
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
	Email    string `json:"email"`
	Password string `json:"password"`
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

func CreateUser(c *gin.Context) {
	// return's a JWT token
	client := db.MongoCl.Client
	collection := client.Database("photos").Collection("users")
	var user User
	if err := c.BindJSON(&user); err != nil {
		log.Fatal(err.Error())
	}
	username := user.Username
	password := user.Password
	email := user.Email
	res, err := collection.InsertOne(context.Background(), bson.M{"username": username, "password": HashPassword(password), "email": email})
	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Println("==========")
	fmt.Println(id)
	fmt.Println(res)

	claims := &Claims{
		Username:         user.Username,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		c.String(http.StatusInternalServerError, fmt.Sprintf("There was an error creating the JWT token"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Create the Signin handler
// func Signin(c *gin.Context) {
// 	// TODO: to signin query password from DB
//
// 	var creds Credentials
// 	// Get the JSON body and decode into credentials
// 	// TODO: remove creds struct
// 	err := json.NewDecoder(c.Request.Body).Decode(&creds)
// 	if err != nil {
// 		// If the structure of the body is wrong, return an HTTP error
// 		c.String(http.StatusPartialContent, fmt.Sprintf("no request body was sent..."))
// 		return
// 	}
//
// 	// Get the expected password from our in memory map
// 	expectedPassword, ok := users[creds.Username]
//
// 	// If a password exists for the given user
// 	// AND, if it is the same as the password we received, the we can move ahead
// 	// if NOT, then we return an "Unauthorized" status
// 	if !ok || expectedPassword != creds.Password {
// 		c.String(http.StatusUnauthorized, fmt.Sprintf("Password invalid"))
// 		return
// 	}
//
// 	// Declare the expiration time of the token
// 	// here, we have kept it as 5 minutes
// 	expirationTime := time.Now().Add(5 * time.Minute)
// 	// Create the JWT claims, which includes the username and expiry time
// 	claims := &Claims{
// 		Username:         creds.Username,
// 		RegisteredClaims: jwt.RegisteredClaims{},
// 	}
//
// 	// Declare the token with the algorithm used for signing, and the claims
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	// Create the JWT string
// 	tokenString, err := token.SignedString(jwtKey)
// 	if err != nil {
// 		// If there is an error in creating the JWT return an internal server error
// 		c.String(http.StatusInternalServerError, fmt.Sprintf("There was an error creating the JWT token"))
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"code":    http.StatusOK,
// 		"message": string(tokenString), // cast it to string before showing
// 	})
// }
