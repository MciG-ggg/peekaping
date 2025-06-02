package monitor

import (
	"context"
	"errors"
	"fmt"
	"peekaping/src/config"
	"peekaping/src/modules/heartbeat"
	"peekaping/src/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoModel struct {
	ID             primitive.ObjectID      `bson:"_id"`
	Type           string                  `bson:"type"`
	Name           string                  `bson:"name"`
	Interval       int                     `bson:"interval"`
	Timeout        int                     `bson:"timeout"`
	MaxRetries     int                     `bson:"max_retries"`
	RetryInterval  int                     `bson:"retry_interval"`
	ResendInterval int                     `bson:"resend_interval"`
	Active         bool                    `bson:"active"`
	Status         heartbeat.MonitorStatus `bson:"status"`
	CreatedAt      time.Time               `bson:"created_at"`
	UpdatedAt      time.Time               `bson:"updated_at"`
	Config         string                  `bson:"config"`
	ProxyId        primitive.ObjectID      `bson:"proxy_id"`
}

type mongoUpdateModel struct {
	Type           *string                  `bson:"type,omitempty"`
	Name           *string                  `bson:"name,omitempty"`
	Interval       *int                     `bson:"interval,omitempty"`
	Timeout        *int                     `bson:"timeout,omitempty"`
	MaxRetries     *int                     `bson:"max_retries,omitempty"`
	RetryInterval  *int                     `bson:"retry_interval,omitempty"`
	ResendInterval *int                     `bson:"resend_interval,omitempty"`
	Active         *bool                    `bson:"active,omitempty"`
	Status         *heartbeat.MonitorStatus `bson:"status,omitempty"`
	CreatedAt      *time.Time               `bson:"created_at,omitempty"`
	UpdatedAt      *time.Time               `bson:"updated_at,omitempty"`
	Config         *string                  `bson:"config,omitempty"`
	ProxyId        *primitive.ObjectID      `bson:"proxy_id,omitempty"`
}

// func toMongoModel(m *Model) (*mongoModel, error) {
// 	var oid primitive.ObjectID
// 	var err error
// 	if m.ID != "" {
// 		oid, err = primitive.ObjectIDFromHex(m.ID)
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		oid = primitive.NewObjectID()
// 	}
// 	return &mongoModel{
// 		ID:             oid,
// 		Type:           m.Type,
// 		Name:           m.Name,
// 		Url:            m.Url,
// 		Interval:       m.Interval,
// 		Timeout:        m.Timeout,
// 		MaxRetries:     m.MaxRetries,
// 		RetryInterval:  m.RetryInterval,
// 		ResendInterval: m.ResendInterval,
// 		Active:         m.Active,
// 		Status:         int(m.Status),
// 		CreatedAt:      m.CreatedAt,
// 		UpdatedAt:      m.UpdatedAt,
// 	}, nil
// }

func toDomainModel(mm *mongoModel) *Model {
	return &Model{
		ID:             mm.ID.Hex(),
		Type:           mm.Type,
		Name:           mm.Name,
		Interval:       mm.Interval,
		Timeout:        mm.Timeout,
		MaxRetries:     mm.MaxRetries,
		RetryInterval:  mm.RetryInterval,
		ResendInterval: mm.ResendInterval,
		Active:         mm.Active,
		Status:         mm.Status,
		CreatedAt:      mm.CreatedAt,
		UpdatedAt:      mm.UpdatedAt,
		Config:         mm.Config,
		ProxyId:        mm.ProxyId.Hex(),
	}
}

type MonitorRepositoryImpl struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func NewMonitorRepository(client *mongo.Client, cfg *config.Config) MonitorRepository {
	db := client.Database(cfg.DBName)
	collection := db.Collection("monitor")
	ctx := context.Background()

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "active", Value: 1},
			{Key: "status", Value: 1},
			{Key: "created_at", Value: -1},
		},
	})
	if err != nil {
		panic("Failed to create index on monitor collection:" + err.Error())
	}

	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: 1}},
	})
	if err != nil {
		panic("Failed to create index on monitor collection:" + err.Error())
	}

	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "url", Value: 1}},
	})
	if err != nil {
		panic("Failed to create index on monitor collection:" + err.Error())
	}

	return &MonitorRepositoryImpl{client, db, collection}
}

func (r *MonitorRepositoryImpl) Create(ctx context.Context, monitor *Model) (*Model, error) {
	proxyObjectID, err := primitive.ObjectIDFromHex(monitor.ProxyId)
	if err != nil {
		return nil, err
	}

	mm := &mongoModel{
		ID:             primitive.NewObjectID(),
		Type:           monitor.Type,
		Name:           monitor.Name,
		Interval:       monitor.Interval,
		Timeout:        monitor.Timeout,
		MaxRetries:     monitor.MaxRetries,
		RetryInterval:  monitor.RetryInterval,
		ResendInterval: monitor.ResendInterval,
		Active:         monitor.Active,
		Status:         0, // Default or set as needed
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Config:         monitor.Config,
		ProxyId:        proxyObjectID,
	}

	_, err = r.collection.InsertOne(ctx, mm)
	if err != nil {
		return nil, err
	}

	return toDomainModel(mm), nil
}

func (r *MonitorRepositoryImpl) FindByID(ctx context.Context, id string) (*Model, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var mm mongoModel
	err = r.collection.FindOne(ctx, filter).Decode(&mm)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainModel(&mm), nil
}

func (r *MonitorRepositoryImpl) FindAll(
	ctx context.Context,
	page int,
	limit int,
	q string,
	active *bool,
	status *int,
) ([]*Model, error) {
	var monitors []*Model

	// Calculate the number of documents to skip
	skip := int64((page) * limit)
	limit64 := int64(limit)

	// Define options for pagination
	options := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit64,
		Sort:  bson.D{{Key: "created_at", Value: -1}},
	}

	filter := bson.M{}
	if q != "" {
		filter["$or"] = bson.A{
			bson.M{"name": bson.M{"$regex": q, "$options": "i"}},
			bson.M{"url": bson.M{"$regex": q, "$options": "i"}},
		}
	}
	if active != nil {
		filter["active"] = *active
	}
	if status != nil {
		filter["status"] = *status
	}

	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var mm mongoModel
		if err := cursor.Decode(&mm); err != nil {
			return nil, err
		}
		monitors = append(monitors, toDomainModel(&mm))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return monitors, nil
}

func (r *MonitorRepositoryImpl) UpdateFull(ctx context.Context, id string, monitor *Model) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	proxyObjectID, err := primitive.ObjectIDFromHex(monitor.ProxyId)
	if err != nil {
		return err
	}

	mm := &mongoModel{
		ID:             objectID,
		Type:           monitor.Type,
		Name:           monitor.Name,
		Interval:       monitor.Interval,
		Timeout:        monitor.Timeout,
		MaxRetries:     monitor.MaxRetries,
		RetryInterval:  monitor.RetryInterval,
		ResendInterval: monitor.ResendInterval,
		Active:         monitor.Active,
		Status:         0, // Default or set as needed
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Config:         monitor.Config,
		ProxyId:        proxyObjectID,
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": mm}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MonitorRepositoryImpl) UpdatePartial(ctx context.Context, id string, monitor *UpdateModel) error {
	var proxyObjectID *primitive.ObjectID
	if monitor.ProxyId != nil {
		objectID, err := primitive.ObjectIDFromHex(*monitor.ProxyId)
		if err != nil {
			return err
		}
		proxyObjectID = &objectID
	}

	mu := &mongoUpdateModel{
		Type:           monitor.Type,
		Name:           monitor.Name,
		Interval:       monitor.Interval,
		Timeout:        monitor.Timeout,
		MaxRetries:     monitor.MaxRetries,
		RetryInterval:  monitor.RetryInterval,
		ResendInterval: monitor.ResendInterval,
		Active:         monitor.Active,
		Status:         monitor.Status,
		CreatedAt:      monitor.CreatedAt,
		UpdatedAt:      monitor.UpdatedAt,
		Config:         monitor.Config,
		ProxyId:        proxyObjectID,
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	set, err := utils.ToBsonSet(mu)
	if err != nil {
		return err
	}
	fmt.Println("set", set)

	if len(set) == 0 {
		return errors.New("nothing to update")
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": set}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete removes a monitor from the MongoDB collection by its ID.
func (r *MonitorRepositoryImpl) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}

// FindActive retrieves all active monitors from the MongoDB collection.
func (r *MonitorRepositoryImpl) FindActive(ctx context.Context) ([]*Model, error) {
	var monitors []*Model

	// Define options for pagination
	options := &options.FindOptions{}

	// Filter for active monitors
	filter := bson.M{"active": true}

	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var mm mongoModel
		if err := cursor.Decode(&mm); err != nil {
			return nil, err
		}
		monitors = append(monitors, toDomainModel(&mm))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return monitors, nil
}

// RemoveProxyReference sets proxy_id to an empty string for all monitors with the given proxyId.
func (r *MonitorRepositoryImpl) RemoveProxyReference(ctx context.Context, proxyId string) error {
	filter := bson.M{"proxy_id": proxyId}
	update := bson.M{"$set": bson.M{"proxy_id": ""}}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// FindByProxyId returns all monitors using the given proxyId
func (r *MonitorRepositoryImpl) FindByProxyId(ctx context.Context, proxyId string) ([]*Model, error) {
	var monitors []*Model

	objectID, err := primitive.ObjectIDFromHex(proxyId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"proxy_id": objectID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var mm mongoModel
		if err := cursor.Decode(&mm); err != nil {
			return nil, err
		}
		monitors = append(monitors, toDomainModel(&mm))
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return monitors, nil
}
