package credentials

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App
var AuthCl *auth.Client

func init() {

	opt := option.WithCredentialsFile("firebase_credentials.json")
	FirebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return
	}
	AuthCl, err = FirebaseApp.Auth(context.Background())
	if err != nil {
		return
	}
}
