package auction

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getTestMongoURI() string {
	uri := os.Getenv("MONGODB_URL")
	if uri == "" {
		// default to localhost without auth
		return "mongodb://localhost:27017"
	}
	// if the compose file uses the host 'mongodb' (DNS name inside docker network),
	// replace it with localhost so tests run from the host can connect to the container.
	// handle cases like: mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
	uri = strings.Replace(uri, "@mongodb:", "@localhost:", 1)
	uri = strings.Replace(uri, "://mongodb:", "://localhost:", 1)
	uri = strings.ReplaceAll(uri, "mongodb:27017", "localhost:27017")
	return uri
}

func TestAutoCloseAuctions(t *testing.T) {
	// set a short auction interval for fast test
	os.Setenv("AUCTION_INTERVAL", "2s")

	ctx := context.Background()
	mongoURI := getTestMongoURI()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v (uri=%s)", err, mongoURI)
	}
	defer client.Disconnect(ctx)

	db := client.Database("auctions_test")
	// ensure clean collection
	db.Collection("auctions").Drop(ctx)

	repo := NewAuctionRepository(db)

	auction := &auction_entity.Auction{
		Id:          "test-auction-auto-close",
		ProductName: "Test Product",
		Category:    "Test",
		Description: "This is a test auction",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now().Add(-5 * time.Second), // already expired
	}

	if err := repo.CreateAuction(ctx, auction); err != nil {
		t.Fatalf("failed to create auction: %v", err)
	}

	// wait for the auto-close routine to run and update the auction status
	deadline := time.Now().Add(6 * time.Second)
	for time.Now().Before(deadline) {
		var stored AuctionEntityMongo
		err := db.Collection("auctions").FindOne(ctx, bson.M{"_id": auction.Id}).Decode(&stored)
		if err != nil {
			// try again until deadline
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if stored.Status == auction_entity.Completed {
			// success
			db.Collection("auctions").Drop(ctx)
			return
		}

		time.Sleep(500 * time.Millisecond)
	}

	db.Collection("auctions").Drop(ctx)
	t.Fatalf("auction was not auto-closed within expected time")
}
