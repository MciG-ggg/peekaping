package heartbeat

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
	MonitorID primitive.ObjectID `bson:"monitor_id"`
	Status    MonitorStatus      `bson:"status"`
	Msg       string             `bson:"msg"`
	Ping      int                `bson:"ping"`
	Duration  int                `bson:"duration"`
	DownCount int                `bson:"down_count"`
	Retries   int                `bson:"retries"`
	Important bool               `bson:"important"`
	Time      time.Time          `bson:"time"`
	EndTime   time.Time          `bson:"end_time"`
	Notified  bool               `bson:"notified"`
}

type RepositoryImpl struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func toDomainModel(mm *mongoModel) *Model {
	return &Model{
		ID:        mm.ID.Hex(),
		MonitorID: mm.MonitorID.Hex(),
		Status:    mm.Status,
		Msg:       mm.Msg,
		Ping:      mm.Ping,
		Duration:  mm.Duration,
		DownCount: mm.DownCount,
		Retries:   mm.Retries,
		Important: mm.Important,
		Time:      mm.Time,
		EndTime:   mm.EndTime,
		Notified:  mm.Notified,
	}
}

func NewRepository(client *mongo.Client, cfg *config.Config) Repository {
	db := client.Database(cfg.DBName)
	collection := db.Collection("heartbeat")

	ctx := context.Background()

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "monitor_id", Value: 1}, {Key: "time", Value: -1}},
	})
	if err != nil {
		// Log or handle the error as needed, but don't panic to avoid breaking app startup
		// println("Failed to create index on heartbeat collection:", err.Error())
		panic("Failed to create index on heartbeat collection:" + err.Error())
	}

	return &RepositoryImpl{client, db, collection}
}

func (r *RepositoryImpl) Create(ctx context.Context, entity *Model) (*Model, error) {
	monitorID, err := primitive.ObjectIDFromHex(entity.MonitorID)
	if err != nil {
		return nil, err
	}

	mm := &mongoModel{
		ID:        primitive.NewObjectID(),
		MonitorID: monitorID,
		Status:    entity.Status,
		Msg:       entity.Msg,
		Ping:      entity.Ping,
		Duration:  entity.Duration,
		DownCount: entity.DownCount,
		Retries:   entity.Retries,
		Important: entity.Important,
		Time:      entity.Time,
		EndTime:   entity.EndTime,
		Notified:  entity.Notified,
	}

	_, err = r.collection.InsertOne(ctx, mm)
	if err != nil {
		return nil, err
	}

	return toDomainModel(mm), nil
}

func (r *RepositoryImpl) FindByID(ctx context.Context, id string) (*Model, error) {
	var mm mongoModel

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(&mm)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainModel(&mm), nil
}

func (r *RepositoryImpl) FindAll(ctx context.Context, page int, limit int) ([]*Model, error) {
	var entities []*Model

	// Calculate the number of documents to skip
	skip := int64((page) * limit)
	limit64 := int64(limit)

	// Define options for pagination
	options := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit64,
	}

	cursor, err := r.collection.Find(ctx, bson.M{}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var monitor Model
		if err := cursor.Decode(&monitor); err != nil {
			return nil, err
		}
		entities = append(entities, &monitor)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return entities, nil
}

// func (r *RepositoryImpl) Update(entity *Model) error {
// 	filter := bson.M{"_id": entity.ID}
// 	update := bson.M{"$set": entity}
// 	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
// 	return err
// }

func (r *RepositoryImpl) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *RepositoryImpl) FindActive(ctx context.Context) ([]*Model, error) {
	var entities []*Model

	options := &options.FindOptions{}

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
		entities = append(entities, toDomainModel(&mm))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return entities, nil
}

// FindByMonitorIDAndTimeRange fetches heartbeat data and aggregates it into one-minute intervals
func (r *RepositoryImpl) FindByMonitorIDAndTimeRange(
	ctx context.Context,
	monitorID string,
	startTime,
	endTime time.Time,
) ([]*ChartPoint, error) {
	objectID, err := primitive.ObjectIDFromHex(monitorID)
	if err != nil {
		return nil, err
	}

	// Round startTime down and endTime up to minute boundaries
	startTime = startTime.Truncate(time.Minute)
	endTime = endTime.Truncate(time.Minute).Add(time.Minute)

	// MongoDB aggregation pipeline
	pipeline := bson.A{
		// Match documents for the monitor and time range
		bson.M{
			"$match": bson.M{
				"monitor_id": objectID,
				"time": bson.M{
					"$gte": startTime,
					"$lte": endTime,
				},
			},
		},
		// Group by minute
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					// Truncate time to minute
					"$dateTrunc": bson.M{
						"date": "$time",
						"unit": "minute",
					},
				},
				"up": bson.M{
					"$sum": bson.M{
						"$cond": bson.M{
							"if":   bson.M{"$eq": []interface{}{"$status", 1}},
							"then": 1,
							"else": 0,
						},
					},
				},
				"down": bson.M{
					"$sum": bson.M{
						"$cond": bson.M{
							"if":   bson.M{"$eq": []interface{}{"$status", 0}},
							"then": 1,
							"else": 0,
						},
					},
				},
				"avgPing": bson.M{
					"$avg": "$ping",
				},
				"minPing": bson.M{
					"$min": "$ping",
				},
				"maxPing": bson.M{
					"$max": "$ping",
				},
			},
		},
		// Sort by timestamp
		bson.M{
			"$sort": bson.M{
				"_id": 1,
			},
		},
	}

	// Execute aggregation
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Temporary struct to decode aggregation results
	type aggResult struct {
		ID      time.Time `bson:"_id"`
		Up      int       `bson:"up"`
		Down    int       `bson:"down"`
		AvgPing float64   `bson:"avgPing"`
		MinPing int       `bson:"minPing"`
		MaxPing int       `bson:"maxPing"`
	}

	var results []aggResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// Create a map of timestamps to ChartPoints for quick lookup
	chartMap := make(map[int64]*ChartPoint)
	for _, res := range results {
		timestamp := res.ID.Unix() * 1000 // Convert to epoch milliseconds
		chartMap[timestamp] = &ChartPoint{
			Up:        res.Up,
			Down:      res.Down,
			AvgPing:   res.AvgPing,
			MinPing:   res.MinPing,
			MaxPing:   res.MaxPing,
			Timestamp: timestamp,
		}
	}

	// Generate ChartPoints for every minute in the range
	var chartPoints []*ChartPoint
	currentTime := startTime
	for currentTime.Before(endTime) {
		timestamp := currentTime.Unix() * 1000
		if point, exists := chartMap[timestamp]; exists {
			chartPoints = append(chartPoints, point)
		} else {
			// No data for this minute, append empty ChartPoint
			chartPoints = append(chartPoints, &ChartPoint{
				Up:        0,
				Down:      0,
				AvgPing:   0,
				MinPing:   0,
				MaxPing:   0,
				Timestamp: timestamp,
			})
		}
		currentTime = currentTime.Add(time.Minute)
	}

	return chartPoints, nil
}

func (r *RepositoryImpl) FindLastNByMonitorID(ctx context.Context, monitorID string, n int) ([]*Model, error) {
	objectID, err := primitive.ObjectIDFromHex(monitorID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"monitor_id": objectID}
	options := &options.FindOptions{
		Sort:  bson.M{"time": -1},
		Limit: int64Ptr(int64(n)),
	}
	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*Model
	for cursor.Next(ctx) {
		var hb Model
		if err := cursor.Decode(&hb); err != nil {
			return nil, err
		}
		results = append(results, &hb)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// Reverse the results before returning
	for i, j := 0, len(results)-1; i < j; i, j = i+1, j-1 {
		results[i], results[j] = results[j], results[i]
	}

	return results, nil
}

func int64Ptr(i int64) *int64 {
	return &i
}

func (r *RepositoryImpl) FindUptimeStatsByMonitorID(ctx context.Context, monitorID string, periods map[string]time.Duration, now time.Time) (map[string]float64, error) {
	objectID, err := primitive.ObjectIDFromHex(monitorID)
	if err != nil {
		return nil, err
	}

	facets := bson.M{}
	for name, duration := range periods {
		start := now.Add(-duration)
		facets[name] = bson.A{
			bson.M{"$match": bson.M{
				"monitor_id": objectID,
				"time":       bson.M{"$gte": start, "$lte": now},
			}},
			bson.M{"$group": bson.M{
				"_id":  nil,
				"up":   bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$status", 1}}, 1, 0}}},
				"down": bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$status", 0}}, 1, 0}}},
			}},
		}
	}
	pipeline := bson.A{
		bson.M{"$facet": facets},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var facetResult []map[string][]struct {
		Up   int `bson:"up"`
		Down int `bson:"down"`
	}
	if err := cursor.All(ctx, &facetResult); err != nil {
		return nil, err
	}
	if len(facetResult) == 0 {
		return nil, nil
	}
	result := make(map[string]float64)
	for name := range periods {
		arr := facetResult[0][name]
		if len(arr) == 0 {
			result[name] = 0
			continue
		}
		up := arr[0].Up
		down := arr[0].Down
		total := up + down
		if total > 0 {
			result[name] = float64(up) / float64(total) * 100
		} else {
			result[name] = 0
		}
	}
	return result, nil
}

func (r *RepositoryImpl) DeleteOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	filter := bson.M{"time": bson.M{"$lt": cutoff}}
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (r *RepositoryImpl) FindByMonitorIDPaginated(ctx context.Context, monitorID string, limit, page int, important *bool, reverse bool) ([]*Model, error) {
	objectID, err := primitive.ObjectIDFromHex(monitorID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"monitor_id": objectID}
	if important != nil {
		filter["important"] = *important
	}

	skip := int64(page * limit)
	limit64 := int64(limit)
	options := &options.FindOptions{
		Sort:  bson.M{"time": -1}, // Always sort descending
		Limit: &limit64,
		Skip:  &skip,
	}
	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*mongoModel
	for cursor.Next(ctx) {
		var mm mongoModel
		if err := cursor.Decode(&mm); err != nil {
			return nil, err
		}
		results = append(results, &mm)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	models := make([]*Model, len(results))
	for i, mm := range results {
		models[i] = toDomainModel(mm)
	}

	if reverse && len(models) > 1 {
		for i, j := 0, len(models)-1; i < j; i, j = i+1, j-1 {
			models[i], models[j] = models[j], models[i]
		}
	}

	return models, nil
}
