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

type fieldWithSection struct {
	db         *mongo.Database
	log        *log.Logger
	collection string
}

func newFieldWithSection(db *mongo.Database, log *log.Logger, collection string) fieldWithSection {
	return fieldWithSection{
		db:         db,
		log:        log,
		collection: collection,
	}
}

func (fws fieldWithSection) query(ctx context.Context, tenantID primitive.ObjectID) ([]model.FieldWithSection, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}

	cur, err := fws.db.Collection(fws.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, ctx)

	var fields []model.FieldWithSection
	if err := cur.All(ctx, &fields); err != nil {
		return nil, err
	}

	return fields, nil
}

func (fws fieldWithSection) queryByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) (model.FieldWithSection, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: fieldID}}

	var field model.FieldWithSection
	res := fws.db.Collection(fws.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.FieldWithSection{}, ErrNotFound
		}
	}

	if err := res.Decode(&field); err != nil {
		return model.FieldWithSection{}, err
	}

	return field, nil
}

func (fws fieldWithSection) queryBySectionID(ctx context.Context, tenantID, sectionID primitive.ObjectID) ([]model.FieldWithSection, error) {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "section_id", Value: sectionID}}

	cur, err := fws.db.Collection(fws.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, ctx)

	fields := make([]model.FieldWithSection, 0)
	if err := cur.All(ctx, &fields); err != nil {
		return nil, err
	}

	return fields, nil
}

func (fws fieldWithSection) queryByName(ctx context.Context, tenantID primitive.ObjectID, name string) (model.FieldWithSection, error) {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "name", Value: name}}

	var field model.FieldWithSection
	res := fws.db.Collection(fws.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.FieldWithSection{}, ErrNotFound
		}
	}

	if err := res.Decode(&field); err != nil {
		return model.FieldWithSection{}, err
	}

	return field, nil
}

func (fws fieldWithSection) queryByLabel(ctx context.Context, tenantID primitive.ObjectID, label model.Translation) (model.FieldWithSection, error) {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "label.it", Value: label.It}},
			bson.D{{Key: "label.en", Value: label.En}},
			bson.D{{Key: "label.de", Value: label.De}},
			bson.D{{Key: "label.es", Value: label.Es}},
			bson.D{{Key: "label.fr", Value: label.Fr}},
		}}}

	res := fws.db.Collection(fws.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.FieldWithSection{}, ErrNotFound
		}
		return model.FieldWithSection{}, err
	}

	var field model.FieldWithSection
	if err := res.Decode(&field); err != nil {
		return model.FieldWithSection{}, err
	}

	return field, nil
}

func (fws fieldWithSection) checkByLabel(ctx context.Context, tenantID primitive.ObjectID, label model.Translation) error {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "label.it", Value: label.It}},
			bson.D{{Key: "label.en", Value: label.En}},
			bson.D{{Key: "label.de", Value: label.De}},
			bson.D{{Key: "label.es", Value: label.Es}},
			bson.D{{Key: "label.fr", Value: label.Fr}},
		}}}

	count, err := fws.db.Collection(fws.collection).CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	return ErrNotFound
}

func (fws fieldWithSection) insert(ctx context.Context, nf model.FieldWithSection) error {

	_, err := fws.db.Collection(fws.collection).InsertOne(ctx, nf)
	if err != nil {
		return err
	}

	return nil
}

func (fws fieldWithSection) bulkInsert(ctx context.Context, nfs []model.FieldWithSection) error {

	fields := make([]interface{}, 0)

	for _, nf := range nfs {
		fields = append(fields, nf)
	}

	_, err := fws.db.Collection(fws.collection).InsertMany(ctx, fields)
	if err != nil {
		return err
	}

	return nil
}

func (fws fieldWithSection) update(ctx context.Context, field model.FieldWithSection) error {

	filter := bson.D{{Key: "tenant_id", Value: field.TenantID}, {Key: "_id", Value: field.ID}}
	update := bson.D{{Key: "$set", Value: field}}

	_, err := fws.db.Collection(fws.collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return err
}

func (fws fieldWithSection) deleteByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: fieldID}}
	_, err := fws.db.Collection(fws.collection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (fws fieldWithSection) deleteByIDs(ctx context.Context, tenantID primitive.ObjectID, fieldIDs []primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: bson.D{{Key: "$in", Value: fieldIDs}}}}
	_, err := fws.db.Collection(fws.collection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (fws fieldWithSection) deleteBySectionID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "section_id", Value: sectionID}}
	_, err := fws.db.Collection(fws.collection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (fws fieldWithSection) deleteBySectionIDs(ctx context.Context, tenantID primitive.ObjectID, sectionIDs []primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: bson.D{{Key: "$in", Value: sectionIDs}}}}
	_, err := fws.db.Collection(fws.collection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (fws fieldWithSection) delete(ctx context.Context, tenantID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	_, err := fws.db.Collection(fws.collection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
