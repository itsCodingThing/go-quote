package utils

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func InitFirebase(ctx context.Context) *firestore.Client {
	firebaseApp, firebaseInitErr := firebase.NewApp(ctx, nil, option.WithCredentialsFile("./database.json"))
	if firebaseInitErr != nil {
		log.Fatal(firebaseInitErr)
	}

	client, firstoreErr := firebaseApp.Firestore(ctx)
	if firstoreErr != nil {
		log.Fatal(firstoreErr)
	}

	return client
}
