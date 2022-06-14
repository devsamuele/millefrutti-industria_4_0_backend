package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/devsamuele/contact-service/business/data/contact/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type section struct {
	db         *mongo.Database
	log        *log.Logger
	collection string
}

func newSection(db *mongo.Database, log *log.Logger, collection string) section {
	return section{
		db:         db,
		log:        log,
		collection: collection,
	}
}

func (s section) query(ctx context.Context, tenantID primitive.ObjectID) ([]model.Section, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}

	cur, err := s.db.Collection(s.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, ctx)

	sections := make([]model.Section, 0)
	if err := cur.All(ctx, &sections); err != nil {
		return nil, fmt.Errorf("unable to decode sections: %w", err)
	}

	return sections, nil
}

func (s section) queryByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) (model.Section, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: sectionID}}

	var section model.Section
	res := s.db.Collection(s.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Section{}, ErrNotFound
		}
		return model.Section{}, err
	}

	if err := res.Decode(&section); err != nil {
		return model.Section{}, err
	}

	return section, nil
}

func (s section) checkByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: sectionID}}

	count, err := s.db.Collection(s.collection).CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}
	return ErrNotFound
}

func (s section) queryByName(ctx context.Context, tenantID primitive.ObjectID, name string) (model.Section, error) {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "name", Value: name}}

	var section model.Section
	res := s.db.Collection(s.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Section{}, ErrNotFound
		}
		return model.Section{}, err
	}

	if err := res.Decode(&section); err != nil {
		return model.Section{}, err
	}

	return section, nil
}

func (s section) queryByTitle(ctx context.Context, tenantID primitive.ObjectID, title model.Translation) (model.Section, error) {

	filter := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "title.it", Value: title.It}},
			bson.D{{Key: "title.en", Value: title.En}},
			bson.D{{Key: "title.de", Value: title.De}},
			bson.D{{Key: "title.es", Value: title.Es}},
			bson.D{{Key: "title.fr", Value: title.Fr}},
		}}}

	res := s.db.Collection(s.collection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Section{}, ErrNotFound
		}
		return model.Section{}, err
	}

	var section model.Section
	if err := res.Decode(section); err != nil {
		return model.Section{}, err
	}

	return section, nil
}

func (s section) insert(ctx context.Context, ns model.Section) error {

	_, err := s.db.Collection(s.collection).InsertOne(ctx, ns)
	if err != nil {
		return err
	}

	return nil
}

func (s section) bulkInsert(ctx context.Context, nss []model.Section) error {

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

func (s section) update(ctx context.Context, section model.Section) error {

	filter := bson.D{{Key: "tenant_id", Value: section.TenantID}, {Key: "_id", Value: section.ID}}
	update := bson.D{{Key: "$set", Value: section}}

	_, err := s.db.Collection(s.collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return err
}

func (s section) deleteByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: sectionID}}
	_, err := s.db.Collection(s.collection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (s section) deleteByIDs(ctx context.Context, tenantID primitive.ObjectID, sectionIDs []primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: bson.D{{Key: "$in", Value: sectionIDs}}}}
	_, err := s.db.Collection(s.collection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (s section) delete(ctx context.Context, tenantID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	_, err := s.db.Collection(s.collection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
