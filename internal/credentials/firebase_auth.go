package credentials

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App
var AuthCl *auth.Client

// this is not going to be used per se, exchanged for mongo
var StoreCl *firestore.Client

func init() {

	opt := option.WithCredentialsFile("firebase_credentials.json")
	ctx := context.Background()

	FirebaseApp, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return
	}

	AuthCl, err = FirebaseApp.Auth(ctx)
	if err != nil {
		return
	}

	StoreCl, err = FirebaseApp.Firestore(ctx)
	if err != nil {
		return
	}

}
