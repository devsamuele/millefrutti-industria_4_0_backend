package resource

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrLabelExists = errors.New("label already exists")
)

// type Storer[T any] interface {
// 	Count(ctx context.Context, tenantID primitive.ObjectID) (int64, error)
// 	CheckByID(ctx context.Context, tenantID, resourceID primitive.ObjectID) error
// 	QueryByID(ctx context.Context, tenantID, resourceID primitive.ObjectID) (Resource, error)
// 	QueryResponseByID(ctx context.Context, tenantID, resourceID primitive.ObjectID, fields []field.Fielder) (T, error)

// 	Insert(ctx context.Context, resource Resource) error
// 	BulkInsert(ctx context.Context, resources []Resource) error
// 	AddFieldValue(ctx context.Context, tenantID, fieldID primitive.ObjectID, value interface{}, now time.Time) error

// 	Update(ctx context.Context, resource Resource) error
// 	UpdateFieldValue(ctx context.Context, tenantID primitive.ObjectID, fieldID primitive.ObjectID, value interface{}, now time.Time) error

// 	DeleteByID(ctx context.Context, tenantID, resourceID primitive.ObjectID) error
// 	DeleteByIDs(ctx context.Context, tenantID primitive.ObjectID, resourceIDs []primitive.ObjectID) error
// 	DeleteByFieldID(ctx context.Context, tenantID primitive.ObjectID, fieldID primitive.ObjectID, fieldValue interface{}) error

// 	UnsetFieldValues(ctx context.Context, tenantID primitive.ObjectID, fieldIDs []primitive.ObjectID, now time.Time) error
// 	UnsetAllFieldValues(ctx context.Context, tenantID primitive.ObjectID, now time.Time) error

// 	PullChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error
// 	RemoveChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error

// 	Search(ctx context.Context, tenantID primitive.ObjectID, fields []field.Fielder, filter bson.D, sort filter.Sort, pageNumber, itemsPerPage int) ([]T, error)
// }

type Store struct {
	db       *mongo.Database
	log      *log.Logger
	resource string
}

func NewStore(db *mongo.Database, log *log.Logger, resource string) Store {
	return Store{
		db:       db,
		log:      log,
		resource: resource,
	}
}

func (s Store) Count(ctx context.Context, tenantID primitive.ObjectID) (int64, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	count, err := s.db.Collection(s.resource).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, err
}

func (s Store) CheckByID(ctx context.Context, tenantID, resourceID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: resourceID}}

	res := s.db.Collection(s.resource).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (s Store) QueryByID(ctx context.Context, tenantID, resourceID primitive.ObjectID) (Resource, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: resourceID}}

	res := s.db.Collection(s.resource).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Resource{}, ErrNotFound
		}
		return Resource{}, err
	}

	var r Resource
	if err := res.Decode(&r); err != nil {
		return Resource{}, err
	}

	return r, nil
}

// func (s Store) QueryResponseByID(ctx context.Context, tenantID, resourceID primitive.ObjectID, r Resource) (Response, error) {
// 	var fkPipeline mongo.Pipeline
// 	for _, fkResource := range r.fkResourceMap {

// 		if fkResource.LocalForeignField {
// 			fields, err := fkResource.fieldStore.Query(ctx, tenantID)
// 			if err != nil {
// 				return nil, err
// 			}
// 			f, found := slice.Find(fields, func(f Fielder) bool {
// 				return f.GetBase().Type == FieldTypeForeignKey && f.GetBase().Name == r.name
// 			})
// 			if found {
// 				lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
// 					{Key: "from", Value: fkResource.Name},
// 					{Key: "localField", Value: "_id"},
// 					{Key: "foreignField", Value: fmt.Sprintf("field_values.%v", f.GetBase().ID.Hex())},
// 					{Key: "as", Value: fkResource.Name},
// 				}}}
// 				fkPipeline = append(fkPipeline, lookupStage)
// 			}
// 		} else {
// 			fields, err := r.fieldStore.Query(ctx, tenantID)
// 			if err != nil {
// 				return nil, err
// 			}
// 			f, found := slice.Find(fields, func(f Fielder) bool {
// 				return f.GetBase().Type == FieldTypeForeignKey && f.GetBase().Name == fkResource.name
// 			})

// 			if found {
// 				lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
// 					{Key: "from", Value: r.Name},
// 					{Key: "localField", Value: fmt.Sprintf("field_values.%v", f.GetBase().ID.Hex())},
// 					{Key: "foreignField", Value: "_id"},
// 					{Key: "as", Value: r.Name},
// 				}}}
// 				fkPipeline = append(fkPipeline, lookupStage)
// 			}
// 		}

// 		if !fkResource.RelOneToMany {
// 			unwindStage := bson.D{{
// 				Key: "$unwind",
// 				Value: bson.D{{
// 					Key:   "path",
// 					Value: fmt.Sprintf("$%v", fkResource.name),
// 				}, {
// 					Key:   "preserveNullAndEmptyArrays",
// 					Value: true,
// 				}},
// 			}}
// 			fkPipeline = append(fkPipeline, unwindStage)
// 		}
// 	}

// 	matchTenantStage := bson.D{{Key: "$match", Value: bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: resourceID}}}}

// 	var pipeline mongo.Pipeline
// 	pipeline = append(pipeline, matchTenantStage)
// 	pipeline = append(pipeline, fkPipeline...)

// 	cur, err := s.db.Collection(s.resource).Aggregate(ctx, pipeline)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cur.Close(ctx)

// 	var resources []Response
// 	for cur.Next(ctx) {
// 		var rawBson bson.Raw
// 		if err := cur.Decode(&rawBson); err != nil {
// 			return nil, err
// 		}
// 		resources = append(resources, Response(rawBson))
// 	}

// 	if len(resources) == 0 {
// 		return nil, ErrNotFound
// 	}

// 	return resources[0], nil
// }

func (s Store) DeleteByFieldID(ctx context.Context, tenantID primitive.ObjectID, fieldID primitive.ObjectID, fieldValue interface{}) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: fieldValue}}
	_, err := s.db.Collection(s.resource).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (s Store) Insert(ctx context.Context, resource Resource) error {

	if _, err := s.db.Collection(s.resource).InsertOne(ctx, resource); err != nil {
		return err
	}

	return nil
}

func (s Store) BulkInsert(ctx context.Context, resources []Resource) error {

	_resources := make([]interface{}, 0)
	for _, resource := range resources {
		_resources = append(_resources, resource)
	}

	if _, err := s.db.Collection(s.resource).InsertMany(ctx, _resources); err != nil {
		return err
	}

	return nil
}

func (s Store) Update(ctx context.Context, resource Resource) error {

	filter := bson.D{{Key: "tenant_id", Value: resource.TenantID}, {Key: "_id", Value: resource.ID}}
	update := bson.D{{Key: "$set", Value: resource}}

	_, err := s.db.Collection(s.resource).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) AddFieldValue(ctx context.Context, tenantID, fieldID primitive.ObjectID, value interface{}, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: value}, {Key: "updated", Value: now}}}}

	_, err := s.db.Collection(s.resource).UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) UpdateFieldValue(ctx context.Context, tenantID primitive.ObjectID, fieldID primitive.ObjectID, value interface{}, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: value}, {Key: "updated", Value: now}}}}

	if _, err := s.db.Collection(s.resource).UpdateMany(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (s Store) DeleteByIDs(ctx context.Context, tenantID primitive.ObjectID, resourceIDs []primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: bson.D{{Key: "$in", Value: resourceIDs}}}}
	_, err := s.db.Collection(s.resource).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (s Store) DeleteByID(ctx context.Context, tenantID, resourceID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: resourceID}}
	_, err := s.db.Collection(s.resource).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (s Store) UnsetFieldValues(ctx context.Context, tenantID primitive.ObjectID, fieldIDs []primitive.ObjectID, now time.Time) error {

	models := make([]mongo.WriteModel, 0)
	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	for _, fieldID := range fieldIDs {
		update := bson.D{{Key: "$unset", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: ""}, {Key: "updated", Value: now}}}}
		models = append(models, mongo.NewUpdateManyModel().SetFilter(filter).SetUpdate(update))
	}

	if _, err := s.db.Collection(s.resource).BulkWrite(ctx, models); err != nil {
		return err
	}

	return nil
}

func (s Store) UnsetAllFieldValues(ctx context.Context, tenantID primitive.ObjectID, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "field_values", Value: bson.D{}}, {Key: "updated", Value: now}}}}

	if _, err := s.db.Collection(s.resource).UpdateMany(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (s Store) PullChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error {

	// make bulk
	for _, choiceID := range choiceIDs {
		filter := bson.D{{Key: "tenant_id", Value: tenantID}}
		update := bson.D{
			{Key: "$pull", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: choiceID}}},
			{Key: "$set", Value: bson.D{{Key: "updated", Value: now}}}}
		if _, err := s.db.Collection(s.resource).UpdateMany(ctx, filter, update); err != nil {
			return err
		}
	}
	return nil
}

func (s Store) RemoveChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error {

	// make bulk
	for _, choiceID := range choiceIDs {
		filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: choiceID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: nil}, {Key: "updated", Value: now}}}}
		if _, err := s.db.Collection(s.resource).UpdateMany(ctx, filter, update); err != nil {
			return err
		}
	}

	return nil
}

// func (s Store) Search(ctx context.Context, tenantID primitive.ObjectID, r Resource, filter bson.D, sort FilterSort, pageNumber, itemsPerPage int) ([]Response, error) {
// 	var fkPipeline mongo.Pipeline
// 	for _, fkResource := range r.fkResourceMap {

// 		if fkResource.LocalForeignField {
// 			fields, err := fkResource.fieldStore.Query(ctx, tenantID)
// 			if err != nil {
// 				return nil, err
// 			}
// 			f, found := slice.Find(fields, func(f Fielder) bool {
// 				return f.GetBase().Type == FieldTypeForeignKey && f.GetBase().Name == r.name
// 			})
// 			if found {
// 				lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
// 					{Key: "from", Value: fkResource.Name},
// 					{Key: "localField", Value: "_id"},
// 					{Key: "foreignField", Value: fmt.Sprintf("field_values.%v", f.GetBase().ID.Hex())},
// 					{Key: "as", Value: fkResource.Name},
// 				}}}
// 				fkPipeline = append(fkPipeline, lookupStage)
// 			}
// 		} else {
// 			fields, err := r.fieldStore.Query(ctx, tenantID)
// 			if err != nil {
// 				return nil, err
// 			}
// 			f, found := slice.Find(fields, func(f Fielder) bool {
// 				return f.GetBase().Type == FieldTypeForeignKey && f.GetBase().Name == fkResource.name
// 			})

// 			if found {
// 				lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
// 					{Key: "from", Value: fkResource.Name},
// 					{Key: "localField", Value: fmt.Sprintf("field_values.%v", f.GetBase().ID.Hex())},
// 					{Key: "foreignField", Value: "_id"},
// 					{Key: "as", Value: fkResource.Name},
// 				}}}
// 				fkPipeline = append(fkPipeline, lookupStage)
// 			}
// 		}

// 		if !fkResource.RelOneToMany {
// 			unwindStage := bson.D{{
// 				Key: "$unwind",
// 				Value: bson.D{{
// 					Key:   "path",
// 					Value: fmt.Sprintf("$%v", fkResource.name),
// 				}, {
// 					Key:   "preserveNullAndEmptyArrays",
// 					Value: true,
// 				}},
// 			}}
// 			fkPipeline = append(fkPipeline, unwindStage)
// 		}
// 	}

// 	if sort.Field == "" {
// 		sort.Field = "name"
// 	}

// 	if sort.Order == 0 {
// 		sort.Order = 1
// 	}

// 	matchTenantStage := bson.D{{Key: "$match", Value: bson.D{{Key: "tenant_id", Value: tenantID}}}}
// 	matchFilterStage := bson.D{{Key: "$match", Value: filter}}
// 	limitStage := bson.D{{Key: "$limit", Value: itemsPerPage}}
// 	skipStage := bson.D{{Key: "$skip", Value: (pageNumber - 1) * itemsPerPage}}
// 	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: sort.Field, Value: sort.Order}}}}

// 	var pipeline mongo.Pipeline
// 	pipeline = append(pipeline, matchTenantStage)
// 	pipeline = append(pipeline, fkPipeline...)
// 	pipeline = append(pipeline, matchFilterStage, sortStage, skipStage, limitStage)

// 	cur, err := s.db.Collection(s.resource).Aggregate(ctx, pipeline)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cur.Close(ctx)

// 	var resources []Response
// 	for cur.Next(ctx) {
// 		var rawBson bson.Raw
// 		if err := cur.Decode(&rawBson); err != nil {
// 			return nil, err
// 		}
// 		resources = append(resources, Response(rawBson))
// 	}

// 	return resources, nil
// }
