package credentials

import (
	"context"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type authApp struct {
	app *firebase.App
}

var AuthApp *authApp

func init() {

	opt := option.WithCredentialsFile("firebase_credentials.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return
	}
	AuthApp = &authApp{app: app}
}
