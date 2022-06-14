package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session struct {
	db *mongo.Database
}

func NewSession(db *mongo.Database) Session {
	return Session{db: db}
}

func (s Session) Start(opts ...*options.SessionOptions) (mongo.Session, error) {
	return s.db.Client().StartSession(opts...)
}

type Cfg struct {
	URI     string
	Name    string
	Timeout time.Duration
}

func Open(cfg Cfg) (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		cancel()
		return nil, nil, nil, fmt.Errorf("opening db: %w", err)
	}

	if err := setUpIndexes(ctx, client, cfg.Name); err != nil {
		cancel()
		return nil, nil, nil, fmt.Errorf("opening db: %w", err)
	}

	return client, ctx, cancel, nil
}

// move to admin build tenant
func setUpIndexes(ctx context.Context, client *mongo.Client, dbName string) error {
	db := client.Database(dbName)

	// Organization
	_, err := db.Collection("organization").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "tenant_id", Value: 1}},
		Options: nil,
	})
	if err != nil {
		return fmt.Errorf("creating organization index: %w", err)
	}

	_, err = db.Collection("organization_field").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "tenant_id", Value: 1}},
		Options: nil,
	})
	if err != nil {
		return fmt.Errorf("creating organization_field index: %w", err)
	}

	_, err = db.Collection("organization_section").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "tenant_id", Value: 1}},
		Options: nil,
	})
	if err != nil {
		return fmt.Errorf("creating organization_section index: %w", err)
	}

	// Person
	_, err = db.Collection("person").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "tenant_id", Value: 1}},
			Options: nil},
		{Keys: bson.D{{Key: "office_id", Value: 1}},
			Options: nil},
		{Keys: bson.D{{Key: "organization_id", Value: 1}},
			Options: nil},
		{Keys: bson.D{{Key: "name", Value: 1}},
			Options: nil},
		{Keys: bson.D{{Key: "primary_email", Value: 1}},
			Options: nil},
		{Keys: bson.D{{Key: "others_email", Value: 1}},
			Options: nil},
		{Keys: bson.D{{Key: "phone", Value: 1}},
			Options: nil},
	})
	if err != nil {
		return fmt.Errorf("creating person index: %w", err)
	}

	_, err = db.Collection("person_field").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "tenant_id", Value: 1}},
		Options: nil,
	})
	if err != nil {
		return fmt.Errorf("creating person_field index: %w", err)
	}

	_, err = db.Collection("person_section").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "tenant_id", Value: 1}},
		Options: nil,
	})
	if err != nil {
		return fmt.Errorf("creating person_section index: %w", err)
	}

	// Office
	_, err = db.Collection("office").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "tenant_id", Value: 1}},
			Options: nil},
		{Keys: bson.D{{Key: "organization_id", Value: 1}},
			Options: nil},
	})
	if err != nil {
		return fmt.Errorf("creating office indexes: %w", err)
	}

	_, err = db.Collection("office_field").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "tenant_id", Value: 1}},
		Options: nil,
	})
	if err != nil {
		return fmt.Errorf("creating office_field index: %w", err)
	}

	return nil
}
