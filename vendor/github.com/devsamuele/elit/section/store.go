package section

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNotFound = errors.New("not found")
)

type Store struct {
	db         *mongo.Database
	log        *log.Logger
	collection string
}

func NewStore(db *mongo.Database, log *log.Logger, resource string) Store {
	return Store{
		db:         db,
		log:        log,
		collection: fmt.Sprintf("%s_section", resource),
	}
}

func (s Store) Query(ctx context.Context, tenantID primitive.ObjectID) ([]Section, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}

	cur, err := s.db.Collection(s.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	sections := make([]Section, 0)
	if err := cur.All(ctx, &sections); err != nil {
		return nil, fmt.Errorf("unable to decode sections: %w", err)
	}

	return sections, nil
}

func (s Store) QueryByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) (Section, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: sectionID}}

	var section Section
	res := s.db.Collection(s.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Section{}, ErrNotFound
		}
		return Section{}, err
	}

	if err := res.Decode(&section); err != nil {
		return Section{}, err
	}

	return section, nil
}

func (s Store) CheckByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: sectionID}}

	res := s.db.Collection(s.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s Store) QueryByName(ctx context.Context, tenantID primitive.ObjectID, name string) (Section, error) {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "name", Value: name}}

	var section Section
	res := s.db.Collection(s.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Section{}, ErrNotFound
		}
		return Section{}, err
	}

	if err := res.Decode(&section); err != nil {
		return Section{}, err
	}

	return section, nil
}

func (s Store) Insert(ctx context.Context, ns Section) error {

	_, err := s.db.Collection(s.collection).InsertOne(ctx, ns)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) BulkInsert(ctx context.Context, nss []Section) error {

	sections := make([]interface{}, 0)
	for _, ns := range nss {
		sections = append(sections, ns)
	}

	_, err := s.db.Collection(s.collection).InsertMany(ctx, sections)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) Update(ctx context.Context, section Section) error {

	filter := bson.D{{Key: "tenant_id", Value: section.TenantID}, {Key: "_id", Value: section.ID}}
	update := bson.D{{Key: "$set", Value: section}}

	_, err := s.db.Collection(s.collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return err
}

func (s Store) DeleteByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: sectionID}}
	_, err := s.db.Collection(s.collection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) DeleteByIDs(ctx context.Context, tenantID primitive.ObjectID, sectionIDs []primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: bson.D{{Key: "$in", Value: sectionIDs}}}}
	_, err := s.db.Collection(s.collection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) Delete(ctx context.Context, tenantID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	_, err := s.db.Collection(s.collection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
