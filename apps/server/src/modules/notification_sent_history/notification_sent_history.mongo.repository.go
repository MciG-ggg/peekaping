package notification_sent_history

import (
	"context"
	"errors"
	"peekaping/src/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoModel struct {
	ID        primitive.ObjectID `bson:"_id"`
	Type      string             `bson:"type"`
	MonitorID primitive.ObjectID `bson:"monitor_id"`
	Days      int                `bson:"days"`
	SentAt    time.Time          `bson:"sent_at"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func toDomainModel(mm *mongoModel) *Model {
	return &Model{
		ID:        mm.ID.Hex(),
		Type:      mm.Type,
		MonitorID: mm.MonitorID.Hex(),
		Days:      mm.Days,
		SentAt:    mm.SentAt,
		CreatedAt: mm.CreatedAt,
		UpdatedAt: mm.UpdatedAt,
	}
}

type MongoRepositoryImpl struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func NewMongoRepository(client *mongo.Client, cfg *config.Config) Repository {
	db := client.Database(cfg.DBName)
	collection := db.Collection("notification_sent_history")

	// Create indexes for efficient queries
	ctx := context.Background()

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "monitor_id", Value: 1},
			{Key: "type", Value: 1},
		},
	})
	if err != nil {
		panic("Failed to create index on notification_sent_history collection:" + err.Error())
	}

	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "type", Value: 1},
			{Key: "monitor_id", Value: 1},
			{Key: "days", Value: 1},
		},
	})
	if err != nil {
		panic("Failed to create index on notification_sent_history collection:" + err.Error())
	}

	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "sent_at", Value: 1}},
	})
	if err != nil {
		panic("Failed to create index on notification_sent_history collection:" + err.Error())
	}

	return &MongoRepositoryImpl{client, db, collection}
}

func (r *MongoRepositoryImpl) Create(ctx context.Context, entity *CreateDto) (*Model, error) {
	monitorObjectID, err := primitive.ObjectIDFromHex(entity.MonitorID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	mm := &mongoModel{
		ID:        primitive.NewObjectID(),
		Type:      entity.Type,
		MonitorID: monitorObjectID,
		Days:      entity.Days,
		SentAt:    now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = r.collection.InsertOne(ctx, mm)
	if err != nil {
		return nil, err
	}

	return toDomainModel(mm), nil
}

func (r *MongoRepositoryImpl) FindByTypeMonitorAndDays(ctx context.Context, notificationType string, monitorID string, days int) (*Model, error) {
	monitorObjectID, err := primitive.ObjectIDFromHex(monitorID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"type":       notificationType,
		"monitor_id": monitorObjectID,
		"days":       bson.M{"$lte": days},
	}

	opts := options.FindOne().SetSort(bson.D{{Key: "days", Value: -1}})

	var mm mongoModel
	err = r.collection.FindOne(ctx, filter, opts).Decode(&mm)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainModel(&mm), nil
}

func (r *MongoRepositoryImpl) DeleteByMonitorID(ctx context.Context, monitorID string) error {
	monitorObjectID, err := primitive.ObjectIDFromHex(monitorID)
	if err != nil {
		return err
	}

	filter := bson.M{"monitor_id": monitorObjectID}
	_, err = r.collection.DeleteMany(ctx, filter)
	return err
}