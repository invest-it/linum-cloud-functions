package deleteaccount

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"log"
	"time"
)

var client *firestore.Client

// AuthEvent is the payload of a Firestore Auth event.
type AuthEvent struct {
	Email    string `json:"email"`
	Metadata struct {
		CreatedAt time.Time `json:"createdAt"`
	} `json:"metadata"`
	UID string `json:"uid"`
}

type DocumentToUser map[string][]string

func init() {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}
	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}
}

func deleteBalanceDataDocument(ref string, ctx context.Context) error {
	log.Printf("Now deleting document with ReferenceId: %s\n", ref)
	_, err := client.Collection("balance").Doc(ref).Delete(ctx)
	if err != nil {
		log.Printf("there was an error deleting an document: %v\n", err)
	}
	return err
}

func deleteBalanceDataDocuments(docRefs []string, ctx context.Context) int {
	errorCount := 0
	log.Printf("Found %d documents for user\n", len(docRefs)) // TODO: Remove in prod
	for _, ref := range docRefs {
		err := deleteBalanceDataDocument(ref, ctx)
		if err != nil {
			// TODO: Log errors
			errorCount++
		}
	}
	return errorCount
}

func deleteAccountSettings(ctx context.Context, uid string) {
	_, err := client.Collection("account_settings").Doc(uid).Delete(ctx)
	if err != nil {
		log.Fatalf("there was an error while deleting the account settings: %v\n", err)
	}
}

func deleteUserIdEntry(ctx context.Context, uid string) {
	update := []firestore.Update{
		{
			Path:  uid,
			Value: firestore.Delete,
		},
	}
	_, err := client.Collection("balance").Doc("documentToUser").Update(ctx, update)
	if err != nil {
		log.Fatalf("there was an error when deleting the userId entry: %v\n", err)
	}
}

func DeleteUserData(ctx context.Context, e AuthEvent) error {
	log.Printf("UserId: %s\n", e.UID)

	docToUserSnapshot, err := client.Collection("balance").Doc("documentToUser").Get(ctx)
	if err != nil {
		log.Fatalf("there was an error while fetching the doc 'documentToUser': %v\n", err)
	}

	var documentToUser DocumentToUser
	err = docToUserSnapshot.DataTo(&documentToUser)
	if err != nil {
		log.Fatalf("there was an error parsing the data: %v\n", err)
	}

	docRefs, ok := documentToUser[e.UID]
	if ok {
		log.Printf("Now deleting balance data for UserId: %s\n", e.UID)
		log.Printf("DocRefsInterface: %v\n", docRefs)
		errorCount := deleteBalanceDataDocuments(docRefs, ctx)
		if errorCount > 0 {
			log.Fatal("there where errors deleting the users documents\n")
		}
	} else {
		log.Println("No entry was found for this UserId")
	}

	log.Println("Now deleting UserId entry in 'documentToUser'")
	deleteUserIdEntry(ctx, e.UID)

	log.Printf("Now deleting account-settings for UserId: %s\n", e.UID)
	deleteAccountSettings(ctx, e.UID)

	return nil
}

// TODO: Try to avoid log.fatal whenever possible