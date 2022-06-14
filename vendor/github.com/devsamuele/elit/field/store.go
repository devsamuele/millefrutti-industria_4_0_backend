package field

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/devsamuele/elit/translation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrLabelExists = errors.New("label already exists")
)

type Store struct {
	db         *mongo.Database
	log        *log.Logger
	collection string
}

func NewStore(db *mongo.Database, log *log.Logger, resource string) Store {
	return Store{db, log, fmt.Sprintf("%v_field", resource)}
}

func (s Store) Query(ctx context.Context, tenantID primitive.ObjectID) ([]Fielder, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}

	cur, err := s.db.Collection(s.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var fields []Fielder
	for cur.Next(ctx) {
		var rawBson bson.Raw
		if err := cur.Decode(&rawBson); err != nil {
			return nil, err
		}

		field, err := bsonDecoder(rawBson)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}

	return fields, nil
}

func (s Store) QueryByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) (Fielder, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: fieldID}}

	res := s.db.Collection(s.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
	}

	var rawBson bson.Raw
	if err := res.Decode(&rawBson); err != nil {
		return nil, err
	}

	field, err := bsonDecoder(rawBson)
	if err != nil {
		return nil, err
	}

	return field, nil
}

func (s Store) QueryBySectionID(ctx context.Context, tenantID, sectionID primitive.ObjectID) ([]Fielder, error) {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "section_id", Value: sectionID}}

	cur, err := s.db.Collection(s.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var fields []Fielder
	for cur.Next(ctx) {
		var rawBson bson.Raw
		if err := cur.Decode(&rawBson); err != nil {
			return nil, err
		}

		field, err := bsonDecoder(rawBson)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}

	return fields, nil
}

func (s Store) QueryByName(ctx context.Context, tenantID primitive.ObjectID, name string) (Fielder, error) {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "name", Value: name}}

	var field Fielder
	res := s.db.Collection(s.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
	}

	var rawBson bson.Raw
	if err := res.Decode(&rawBson); err != nil {
		return nil, err
	}

	field, err := bsonDecoder(rawBson)
	if err != nil {
		return nil, err
	}
	return field, nil
}

func (s Store) CheckByLabel(ctx context.Context, tenantID primitive.ObjectID, label translation.Translation) error {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "label.it", Value: label.It}},
			bson.D{{Key: "label.en", Value: label.En}},
			bson.D{{Key: "label.de", Value: label.De}},
			bson.D{{Key: "label.es", Value: label.Es}},
			bson.D{{Key: "label.fr", Value: label.Fr}},
		}}}

	count, err := s.db.Collection(s.collection).CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrLabelExists
	}

	return nil
}

func (s Store) QueryByLabel(ctx context.Context, tenantID primitive.ObjectID, label translation.Translation) (Fielder, error) {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "label.it", Value: label.It}},
			bson.D{{Key: "label.en", Value: label.En}},
			bson.D{{Key: "label.de", Value: label.De}},
			bson.D{{Key: "label.es", Value: label.Es}},
			bson.D{{Key: "label.fr", Value: label.Fr}},
		}}}

	res := s.db.Collection(s.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var rawBson bson.Raw
	if err := res.Decode(&rawBson); err != nil {
		return nil, err
	}

	field, err := bsonDecoder(rawBson)
	if err != nil {
		return nil, err
	}

	return field, nil
}

func (s Store) Insert(ctx context.Context, nf Fielder) error {

	_, err := s.db.Collection(s.collection).InsertOne(ctx, nf)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) BulkInsert(ctx context.Context, nfs []Fielder) error {

	fields := make([]interface{}, 0)

	for _, nf := range nfs {
		fields = append(fields, nf)
	}

	_, err := s.db.Collection(s.collection).InsertMany(ctx, fields)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) Update(ctx context.Context, tenantID, fieldID primitive.ObjectID, field Fielder) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: fieldID}}
	update := bson.D{{Key: "$set", Value: field}}

	_, err := s.db.Collection(s.collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return err
}

func (s Store) DeleteByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: fieldID}}
	_, err := s.db.Collection(s.collection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) DeleteByIDs(ctx context.Context, tenantID primitive.ObjectID, fieldIDs []primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: bson.D{{Key: "$in", Value: fieldIDs}}}}
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

func (s Store) DeleteBySectionID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "section_id", Value: sectionID}}
	_, err := s.db.Collection(s.collection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) DeleteBySectionIDs(ctx context.Context, tenantID primitive.ObjectID, sectionIDs []primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: bson.D{{Key: "$in", Value: sectionIDs}}}}
	_, err := s.db.Collection(s.collection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
