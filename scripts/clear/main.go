package main

import (
	"backend/config"
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func main() {
	
	ctx := context.Background()

	cfg := config.LoadConfig()

	fmt.Println(" Initializing Firebase Services...")
	config.InitFirebase(cfg)

	client := config.FirestoreClient
	if client == nil {
		log.Fatal(" Firestore client is nil")
	}

	fmt.Println(" Starting database cleanup...")

	iter := client.Collections(ctx)
	count := 0
	for {
		coll, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf(" Error listing collections: %v", err)
		}

		fmt.Printf("  Deleting collection: %s\n", coll.ID)
		if err := deleteCollection(ctx, client, coll, 500); err != nil {
			log.Printf(" Failed to delete collection %s: %v", coll.ID, err)
		} else {
			fmt.Printf(" Deleted collection: %s\n", coll.ID)
			count++
		}
	}

	fmt.Printf("\n Database cleanup complete! Deleted %d collections.\n", count)
}

func deleteCollection(ctx context.Context, client *firestore.Client, coll *firestore.CollectionRef, batchSize int) error {
	for {
		
		iter := coll.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}

		if numDeleted < batchSize {
			return nil
		}
	}
}
