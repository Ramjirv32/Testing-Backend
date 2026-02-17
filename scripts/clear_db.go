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
	// Initialize context
	ctx := context.Background()

	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Firebase
	fmt.Println("🚀 Initializing Firebase Services...")
	config.InitFirebase(cfg)

	client := config.FirestoreClient
	if client == nil {
		log.Fatal("❌ Firestore client is nil")
	}

	fmt.Println("🧹 Starting database cleanup...")

	// Get all collections
	iter := client.Collections(ctx)
	count := 0
	for {
		coll, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("❌ Error listing collections: %v", err)
		}

		fmt.Printf("🗑️  Deleting collection: %s\n", coll.ID)
		if err := deleteCollection(ctx, client, coll, 500); err != nil {
			log.Printf("❌ Failed to delete collection %s: %v", coll.ID, err)
		} else {
			fmt.Printf("✅ Deleted collection: %s\n", coll.ID)
			count++
		}
	}

	fmt.Printf("\n✨ Database cleanup complete! Deleted %d collections.\n", count)
}

// deleteCollection deletes a collection and its documents in batches.
func deleteCollection(ctx context.Context, client *firestore.Client, coll *firestore.CollectionRef, batchSize int) error {
	for {
		// Iterate through the documents in the collection in batches
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

		// If no documents were found, the collection is empty
		if numDeleted == 0 {
			return nil
		}

		// Commit the batch
		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}

		// If we deleted fewer than batchSize, we are done
		if numDeleted < batchSize {
			return nil
		}
	}
}
