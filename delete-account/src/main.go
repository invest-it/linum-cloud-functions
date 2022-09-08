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
	_, err := client.Collection("balance").Doc(ref).Delete(ctx)
	log.Printf("there was an error deleting an document: %v\n", err)
	// TODO: Handle those errors more specifically and don't abort the entire process
	return err
}

func deleteBalanceDataDocuments(docRefs interface{}, ctx context.Context) int {
	errorCount := 0
	refs, ok := docRefs.([]string)
	if ok {
		for _, ref := range refs {
			err := deleteBalanceDataDocument(ref, ctx)
			if err != nil {
				// TODO: Log errors
				errorCount++
			}
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

func DeleteUserData(ctx context.Context, e AuthEvent) error {
	log.Printf("UserId: %s\n", e.UID)
	// Currently this is
	docToUser, err := client.Collection("balance").Doc("documentToUser").Get(ctx)
	docRefs, ok := docToUser.Data()["_uid"]
	if ok {
		errorCount := deleteBalanceDataDocuments(docRefs, ctx)
		if errorCount > 0 {
			log.Fatal("there where errors deleting the users documents\n")
		}
	}
	if err != nil {
		log.Fatalf("there was an error while fetching the doc 'documentToUser': %v\n", err)
	}
	
	deleteAccountSettings(ctx, e.UID)
	// _, err = client.Collection("balance").Doc("documentToUser").Delete(ctx)
	// if err != nil {
	//	log.Fatalf("there was an error while the deleting the doc 'documentToUser': %v\n", err)
	// }

	return nil
}
