package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection      *mongo.Collection
	auctionInterval time.Duration
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repo := &AuctionRepository{
		auctionInterval: getAuctionInterval(),
		Collection:      database.Collection("auctions"),
	}

	repo.startAutoCloseRoutine(context.Background())

	return repo
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}

	return duration
}

// startAutoCloseRoutine periodically finds active auctions that passed their end time
// and sets their status to Completed.
func (ar *AuctionRepository) startAutoCloseRoutine(ctx context.Context) {
	ticker := time.NewTicker(time.Second) // check every second
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// cutoff timestamp (unix seconds) for auctions that expired
				cutoff := time.Now().Add(-ar.auctionInterval).Unix()

				filter := bson.M{
					"status":    auction_entity.Active,
					"timestamp": bson.M{"$lte": cutoff},
				}
				update := bson.M{
					"$set": bson.M{"status": auction_entity.Completed},
				}

				res, err := ar.Collection.UpdateMany(ctx, filter, update)
				if err != nil {
					logger.Error("Error trying to auto-close auctions", err)
					continue
				}
				if res.ModifiedCount > 0 {
					logger.Info("Auto-closed auctions" /* optional fields */)
				}
			}
		}
	}()
}
