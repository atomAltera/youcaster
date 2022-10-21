package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	client *mongo.Client
	db     *mongo.Database

	Requests *requestsStore
}

func New(uri string) (*Client, error) {
	dbName, err := getDBName(uri)
	if err != nil {
		return nil, fmt.Errorf("extracting store name from uri: %w", err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("creating mongo client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}

	db := client.Database(dbName)

	requests, err := newRequestsStore(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("init requests collection: %w", err)
	}

	return &Client{
		client: client,
		db:     db,

		Requests: requests,
	}, nil
}

func getDBName(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	name := strings.TrimPrefix(u.Path, "/")
	return name, nil
}
