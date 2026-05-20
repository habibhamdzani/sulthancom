package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sulthancom/pos-backend/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.MongoDB)
	log.Printf("Connected to MongoDB: %s", cfg.MongoDB)

	return &Database{Client: client, DB: db}, nil
}

func (d *Database) Collection(name string) *mongo.Collection {
	return d.DB.Collection(name)
}

func (d *Database) Close() error {
	return d.Client.Disconnect(context.Background())
}
