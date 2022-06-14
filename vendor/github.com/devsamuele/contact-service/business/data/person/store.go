package person

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/devsamuele/elit/field"
// 	"github.com/devsamuele/elit/filter"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// var (
// 	ErrNotFound = errors.New("not found")
// )

// type Store struct {
// 	db  *mongo.Database
// 	log *log.Logger
// }

// func NewStore(db *mongo.Database, log *log.Logger) Store {
// 	return Store{
// 		db:  db,
// 		log: log,
// 	}
// }

// // Field

// // func (s Store) QueryField(ctx context.Context, tenantID primitive.ObjectID) ([]field.Fielder, error) {

// // 	// !only store
// // 	fieldBuilder, err := field.NewBuilder(field.Config{
// // 		Resource: resource,
// // 		DB:       s.db,
// // 		Logger:   s.log,
// // 	})
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	return fieldBuilder.Store().Query(ctx, tenantID)

// // }

// // func (s Store) QueryFieldByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) (field.Fielder, error) {

// // 	// !only store
// // 	fieldBuilder, err := field.NewBuilder(field.Config{
// // 		Resource: resource,
// // 		DB:       s.db,
// // 		Logger:   s.log,
// // 	})
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	return fieldBuilder.Store().QueryByID(ctx, tenantID, fieldID)
// // }

// // func (s Store) QueryFieldBySectionID(ctx context.Context, tenantID, sectionID primitive.ObjectID) ([]field.Fielder, error) {

// // 	// !only store
// // 	fieldBuilder, err := field.NewBuilder(field.Config{
// // 		Resource: resource,
// // 		DB:       s.db,
// // 		Logger:   s.log,
// // 	})
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	return fieldBuilder.Store().QueryBySectionID(ctx, tenantID, sectionID)
// // }

// // func (s Store) QueryFieldByName(ctx context.Context, tenantID primitive.ObjectID, name string) (field.Fielder, error) {

// // 	// !only store
// // 	fieldBuilder, err := field.NewBuilder(field.Config{
// // 		Resource: resource,
// // 		DB:       s.db,
// // 		Logger:   s.log,
// // 	})
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	return fieldBuilder.Store().QueryBySectionID(ctx, tenantID, sectionID)
// // }

// // func (s Store) QueryFieldByLabel(ctx context.Context, tenantID primitive.ObjectID, label Translation) (FieldWithSection, error) {
// // 	return newFieldWithSection(s.db, s.log, personFieldCollection).queryByLabel(ctx, tenantID, label)
// // }

// // func (s Store) CheckFieldByLabel(ctx context.Context, tenantID primitive.ObjectID, label Translation) error {
// // 	return newFieldWithSection(s.db, s.log, personFieldCollection).checkByLabel(ctx, tenantID, label)
// // }

// // func (s Store) InsertField(ctx context.Context, nf FieldWithSection) error {
// // 	return newFieldWithSection(s.db, s.log, personFieldCollection).insert(ctx, nf)
// // }

// // func (s Store) BulkInsertField(ctx context.Context, nfs []FieldWithSection) error {
// // 	return newFieldWithSection(s.db, s.log, personFieldCollection).bulkInsert(ctx, nfs)
// // }

// // func (s Store) UpdateField(ctx context.Context, field FieldWithSection) error {
// // 	return newFieldWithSection(s.db, s.log, personFieldCollection).update(ctx, field)
// // }

// // func (s Store) DeleteFieldByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) error {
// // 	return newFieldWithSection(s.db, s.log, personFieldCollection).deleteByID(ctx, tenantID, fieldID)
// // }

// // func (s Store) DeleteFieldByIDs(ctx context.Context, tenantID primitive.ObjectID, fieldIDs []primitive.ObjectID) error {
// // 	return newFieldWithSection(s.db, s.log, personFieldCollection).deleteByIDs(ctx, tenantID, fieldIDs)
// // }

// // func (s Store) DeleteFieldBySectionID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {
// // 	return newFieldWithSection(s.db, s.log, personFieldCollection).deleteBySectionID(ctx, tenantID, sectionID)
// // }

// // func (s Store) DeleteFieldBySectionIDs(ctx context.Context, tenantID primitive.ObjectID, sectionIDs []primitive.ObjectID) error {
// // 	return newFieldWithSection(s.db, s.log, personFieldCollection).deleteBySectionIDs(ctx, tenantID, sectionIDs)
// // }

// // func (s Store) DeleteField(ctx context.Context, tenantID primitive.ObjectID) error {
// // 	return newFieldWithSection(s.db, s.log, personFieldCollection).delete(ctx, tenantID)
// // }

// // Person

// func (s Store) Count(ctx context.Context, tenantID primitive.ObjectID) (int64, error) {

// 	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
// 	count, err := s.db.Collection(resource).CountDocuments(ctx, filter)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return count, err
// }

// // func (s Store) QueryResponseByID(ctx context.Context, tenantID, personID primitive.ObjectID) (PersonResponse, error) {

// // 	var pipeline mongo.Pipeline

// // 	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: personID}}}}
// // 	pipeline = append(pipeline, matchStage)

// // 	lookupOrganizationStage := bson.D{{Key: "$lookup", Value: bson.D{
// // 		{Key: "from", Value: "organization"},
// // 		{Key: "localField", Value: "organization_id"},
// // 		{Key: "foreignField", Value: "_id"},
// // 		{Key: "as", Value: "organization"},
// // 	}}}

// // 	unwindOrganizationStage := bson.D{{
// // 		Key: "$unwind",
// // 		Value: bson.D{{
// // 			Key:   "path",
// // 			Value: "$organization",
// // 		}, {
// // 			Key:   "preserveNullAndEmptyArrays",
// // 			Value: true,
// // 		}},
// // 	}}

// // 	lookupOfficeStage := bson.D{{Key: "$lookup", Value: bson.D{
// // 		{Key: "from", Value: "office"},
// // 		{Key: "localField", Value: "office_id"},
// // 		{Key: "foreignField", Value: "_id"},
// // 		{Key: "as", Value: "office"},
// // 	}}}

// // 	unwindOfficeStage := bson.D{{
// // 		Key: "$unwind",
// // 		Value: bson.D{{
// // 			Key:   "path",
// // 			Value: "$office",
// // 		}, {
// // 			Key:   "preserveNullAndEmptyArrays",
// // 			Value: true,
// // 		}},
// // 	}}

// // 	pipeline = append(pipeline, lookupOrganizationStage, unwindOrganizationStage, lookupOfficeStage, unwindOfficeStage)

// // 	cur, err := s.db.Collection(resource).Aggregate(ctx, pipeline)
// // 	if err != nil {
// // 		return PersonResponse{}, err
// // 	}

// // 	defer func(cur *mongo.Cursor, ctx context.Context) {
// // 		err := cur.Close(ctx)
// // 		if err != nil {

// // 		}
// // 	}(cur, ctx)

// // 	persons := make([]PersonResponse, 0)
// // 	if err := cur.All(ctx, &persons); err != nil {
// // 		return PersonResponse{}, err
// // 	}

// // 	if len(persons) == 0 {
// // 		return PersonResponse{}, ErrNotFound
// // 	}

// // 	return persons[0], nil
// // }

// func (s Store) CheckFK(ctx context.Context, tenantID primitive.ObjectID, collection string, ID primitive.ObjectID) error {

// 	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: ID}}

// 	res := s.db.Collection(collection).FindOne(ctx, filter)
// 	if err := res.Err(); err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return ErrNotFound
// 		}
// 		return err
// 	}

// 	return nil
// }

// func (s Store) QueryByID(ctx context.Context, tenantID, personID primitive.ObjectID) (Person, error) {

// 	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: personID}}

// 	res := s.db.Collection(resource).FindOne(ctx, filter)
// 	if err := res.Err(); err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return Person{}, ErrNotFound
// 		}
// 		return Person{}, err
// 	}

// 	var person Person
// 	if err := res.Decode(&person); err != nil {
// 		return Person{}, err
// 	}

// 	return person, nil
// }

// func (s Store) Insert(ctx context.Context, person Person) error {

// 	if _, err := s.db.Collection(resource).InsertOne(ctx, person); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s Store) BulkInsert(ctx context.Context, persons []Person) error {

// 	_persons := make([]interface{}, 0)
// 	for _, person := range persons {
// 		_persons = append(_persons, person)
// 	}

// 	if _, err := s.db.Collection(resource).InsertMany(ctx, _persons); err != nil {
// 		return err
// 	}

// 	return nil

// }

// func (s Store) Update(ctx context.Context, person Person) error {

// 	filter := bson.D{{Key: "tenant_id", Value: person.TenantID}, {Key: "_id", Value: person.ID}}
// 	update := bson.D{{Key: "$set", Value: person}}

// 	_, err := s.db.Collection(resource).UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s Store) AddFieldValue(ctx context.Context, tenantID, fieldID primitive.ObjectID, value interface{}, now time.Time) error {

// 	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
// 	update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: value}, {Key: "updated", Value: now}}}}

// 	_, err := s.db.Collection(resource).UpdateMany(ctx, filter, update)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s Store) UpdateFieldValue(ctx context.Context, tenantID primitive.ObjectID, fieldID primitive.ObjectID, value interface{}, now time.Time) error {

// 	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
// 	update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: value}, {Key: "updated", Value: now}}}}

// 	if _, err := s.db.Collection(resource).UpdateMany(ctx, filter, update); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s Store) DeleteByIDs(ctx context.Context, tenantID primitive.ObjectID, personIDs []primitive.ObjectID) error {

// 	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: bson.D{{Key: "$in", Value: personIDs}}}}
// 	_, err := s.db.Collection(resource).DeleteMany(ctx, filter)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s Store) DeleteByID(ctx context.Context, tenantID, personID primitive.ObjectID) error {

// 	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: personID}}
// 	_, err := s.db.Collection(resource).DeleteOne(ctx, filter)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s Store) UnsetFieldValues(ctx context.Context, tenantID primitive.ObjectID, fieldIDs []primitive.ObjectID, now time.Time) error {

// 	models := make([]mongo.WriteModel, 0)
// 	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
// 	for _, fieldID := range fieldIDs {
// 		update := bson.D{{Key: "$unset", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: ""}, {Key: "updated", Value: now}}}}
// 		models = append(models, mongo.NewUpdateManyModel().SetFilter(filter).SetUpdate(update))
// 	}

// 	if _, err := s.db.Collection(resource).BulkWrite(ctx, models); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s Store) UnsetAllFieldValues(ctx context.Context, tenantID primitive.ObjectID, now time.Time) error {

// 	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
// 	update := bson.D{{Key: "$set", Value: bson.D{{Key: "field_values", Value: bson.D{}}, {Key: "updated", Value: now}}}}

// 	if _, err := s.db.Collection(resource).UpdateMany(ctx, filter, update); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s Store) PullChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error {

// 	// make bulk
// 	for _, choiceID := range choiceIDs {
// 		filter := bson.D{{Key: "tenant_id", Value: tenantID}}
// 		update := bson.D{
// 			{Key: "$pull", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: choiceID}}},
// 			{Key: "$set", Value: bson.D{{Key: "updated", Value: now}}}}
// 		if _, err := s.db.Collection(resource).UpdateMany(ctx, filter, update); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (s Store) RemoveChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error {

// 	// make bulk
// 	for _, choiceID := range choiceIDs {
// 		filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: choiceID}}
// 		update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: nil}, {Key: "updated", Value: now}}}}
// 		if _, err := s.db.Collection(resource).UpdateMany(ctx, filter, update); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s Store) Search(ctx context.Context, tenantID primitive.ObjectID, filter bson.D, sort filter.Sort, pageNumber, itemsPerPage int) ([]PersonResponse, error) {

// 	fieldStore := field.NewStore(s.db, s.log, resource)
// 	organizationIDField, err := fieldStore.QueryByName(ctx, tenantID, "organization_id")
// 	if err != nil {
// 		return nil, err
// 	}

// 	officeIDField, err := fieldStore.QueryByName(ctx, tenantID, "office_id")
// 	if err != nil {
// 		return nil, err
// 	}

// 	var pipeline mongo.Pipeline

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

// 	lookupOrganizationStage := bson.D{{Key: "$lookup", Value: bson.D{
// 		{Key: "from", Value: "organization"},
// 		{Key: "localField", Value: fmt.Sprintf("field_values.%v", organizationIDField.GetBase().ID.Hex())},
// 		{Key: "foreignField", Value: "_id"},
// 		{Key: "as", Value: "organization"},
// 	}}}

// 	unwindOrganizationStage := bson.D{{
// 		Key: "$unwind",
// 		Value: bson.D{{
// 			Key:   "path",
// 			Value: "$organization",
// 		}, {
// 			Key:   "preserveNullAndEmptyArrays",
// 			Value: true,
// 		}},
// 	}}

// 	lookupOfficeStage := bson.D{{Key: "$lookup", Value: bson.D{
// 		{Key: "from", Value: "office"},
// 		{Key: "localField", Value: fmt.Sprintf("field_values.%v", officeIDField.GetBase().ID.Hex())},
// 		{Key: "foreignField", Value: "_id"},
// 		{Key: "as", Value: "office"},
// 	}}}

// 	unwindOfficeStage := bson.D{{
// 		Key: "$unwind",
// 		Value: bson.D{{
// 			Key:   "path",
// 			Value: "$office",
// 		}, {
// 			Key:   "preserveNullAndEmptyArrays",
// 			Value: true,
// 		}},
// 	}}

// 	pipeline = append(pipeline, matchTenantStage, lookupOrganizationStage, unwindOrganizationStage,
// 		lookupOfficeStage, unwindOfficeStage, matchFilterStage, sortStage, skipStage, limitStage)

// 	cur, err := s.db.Collection(resource).Aggregate(ctx, pipeline)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cur.Close(ctx)

// 	persons := make([]PersonResponse, 0)
// 	if err := cur.All(ctx, &persons); err != nil {
// 		return nil, err
// 	}

// 	return persons, nil
// }

// // Section

// // func (s Store) QuerySection(ctx context.Context, tenantID primitive.ObjectID) ([]Section, error) {
// // 	return newSection(s.db, s.log, personSectionCollection).query(ctx, tenantID)
// // }

// // func (s Store) QuerySectionByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) (Section, error) {
// // 	return newSection(s.db, s.log, personSectionCollection).queryByID(ctx, tenantID, sectionID)
// // }

// // func (s Store) CheckSectionByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {
// // 	return newSection(s.db, s.log, personSectionCollection).checkByID(ctx, tenantID, sectionID)
// // }

// // func (s Store) QuerySectionByName(ctx context.Context, tenantID primitive.ObjectID, name string) (Section, error) {
// // 	return newSection(s.db, s.log, personSectionCollection).queryByName(ctx, tenantID, name)
// // }

// // func (s Store) QuerySectionByTitle(ctx context.Context, tenantID primitive.ObjectID, title Translation) (Section, error) {
// // 	return newSection(s.db, s.log, personSectionCollection).queryByTitle(ctx, tenantID, title)
// // }

// // func (s Store) InsertSection(ctx context.Context, ns Section) error {
// // 	return newSection(s.db, s.log, personSectionCollection).insert(ctx, ns)
// // }

// // func (s Store) BulkInsertSection(ctx context.Context, nss []Section) error {
// // 	return newSection(s.db, s.log, personSectionCollection).bulkInsert(ctx, nss)
// // }

// // func (s Store) UpdateSection(ctx context.Context, us Section) error {
// // 	return newSection(s.db, s.log, personSectionCollection).update(ctx, us)
// // }

// // func (s Store) DeleteSectionByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {
// // 	return newSection(s.db, s.log, personSectionCollection).deleteByID(ctx, tenantID, sectionID)
// // }

// // func (s Store) DeleteSectionByIDs(ctx context.Context, tenantID primitive.ObjectID, sectionIDs []primitive.ObjectID) error {
// // 	return newSection(s.db, s.log, personSectionCollection).deleteByIDs(ctx, tenantID, sectionIDs)
// // }

// // func (s Store) DeleteSection(ctx context.Context, tenantID primitive.ObjectID) error {
// // 	return newSection(s.db, s.log, personSectionCollection).delete(ctx, tenantID)
// // }
