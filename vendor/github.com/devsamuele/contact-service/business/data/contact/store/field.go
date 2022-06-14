package store

import (
	"context"
	"errors"
	"github.com/devsamuele/contact-service/business/data/contact/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type field struct {
	db         *mongo.Database
	log        *log.Logger
	collection string
}

func newField(db *mongo.Database, log *log.Logger, collection string) field {
	return field{
		db:         db,
		log:        log,
		collection: collection,
	}
}

func (f field) query(ctx context.Context, tenantID primitive.ObjectID) ([]model.Field, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}

	cur, err := f.db.Collection(f.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, ctx)

	var fields []model.Field
	if err := cur.All(ctx, &fields); err != nil {
		return nil, err
	}

	return fields, nil
}

func (f field) queryByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) (model.Field, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: fieldID}}

	var field model.Field
	res := f.db.Collection(f.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Field{}, ErrNotFound
		}
	}

	if err := res.Decode(&field); err != nil {
		return model.Field{}, err
	}

	return field, nil
}

func (f field) queryByName(ctx context.Context, tenantID primitive.ObjectID, name string) (model.Field, error) {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "name", Value: name}}

	var field model.Field
	res := f.db.Collection(f.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Field{}, ErrNotFound
		}
	}

	if err := res.Decode(&field); err != nil {
		return model.Field{}, err
	}

	return field, nil
}

func (f field) checkByLabel(ctx context.Context, tenantID primitive.ObjectID, label model.Translation) error {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "label.it", Value: label.It}},
			bson.D{{Key: "label.en", Value: label.En}},
			bson.D{{Key: "label.de", Value: label.De}},
			bson.D{{Key: "label.es", Value: label.Es}},
			bson.D{{Key: "label.fr", Value: label.Fr}},
		}}}

	count, err := f.db.Collection(f.collection).CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	return ErrNotFound
}

func (f field) queryByLabel(ctx context.Context, tenantID primitive.ObjectID, label model.Translation) (model.Field, error) {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "label.it", Value: label.It}},
			bson.D{{Key: "label.en", Value: label.En}},
			bson.D{{Key: "label.de", Value: label.De}},
			bson.D{{Key: "label.es", Value: label.Es}},
			bson.D{{Key: "label.fr", Value: label.Fr}},
		}}}

	res := f.db.Collection(f.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Field{}, ErrNotFound
		}
		return model.Field{}, err
	}

	var field model.Field
	if err := res.Decode(&field); err != nil {
		return model.Field{}, err
	}

	return field, nil
}

func (f field) insert(ctx context.Context, nf model.Field) error {

	_, err := f.db.Collection(f.collection).InsertOne(ctx, nf)
	if err != nil {
		return err
	}

	return nil
}

func (f field) bulkInsert(ctx context.Context, nfs []model.Field) error {

	fields := make([]interface{}, 0)

	for _, nf := range nfs {
		fields = append(fields, nf)
	}

	_, err := f.db.Collection(f.collection).InsertMany(ctx, fields)
	if err != nil {
		return err
	}

	return nil
}

func (f field) update(ctx context.Context, field model.Field) error {

	filter := bson.D{{Key: "tenant_id", Value: field.TenantID}, {Key: "_id", Value: field.ID}}
	update := bson.D{{Key: "$set", Value: field}}

	_, err := f.db.Collection(f.collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return err
}

func (f field) deleteByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: fieldID}}
	_, err := f.db.Collection(f.collection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (f field) deleteByIDs(ctx context.Context, tenantID primitive.ObjectID, fieldIDs []primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: bson.D{{Key: "$in", Value: fieldIDs}}}}
	_, err := f.db.Collection(f.collection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (f field) delete(ctx context.Context, tenantID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	_, err := f.db.Collection(f.collection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
