package database

import (
	"context"
	"fmt"
	e "github.com/atomAltera/youcaster/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type requestsStore struct {
	col *mongo.Collection
}

func newRequestsStore(ctx context.Context, db *mongo.Database) (*requestsStore, error) {
	col := db.Collection("requests")

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.M{"id": 1},
		},
		{
			Keys: bson.M{"created_at": 1},
		},
		{
			Keys: bson.D{{"status", 1}, {"created_at", 1}},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("creating indexes: %w", err)
	}

	return &requestsStore{col: col}, nil
}

func (s *requestsStore) Create(ctx context.Context, doc e.Request) error {
	_, err := s.col.InsertOne(ctx, doc)
	return err
}

func (s *requestsStore) Update(ctx context.Context, doc e.Request) error {
	filter := bson.M{"id": doc.ID}
	update := bson.M{"$set": doc}
	_, err := s.col.UpdateOne(ctx, filter, update)
	return err
}

func (s *requestsStore) List(ctx context.Context, ss []e.RequestStatus) ([]e.Request, error) {
	filter := bson.M{"status": bson.M{"$in": ss}}
	opts := options.Find().SetSort(bson.M{"created_at": 1})

	cur, err := s.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var docs []e.Request
	if err = cur.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}
