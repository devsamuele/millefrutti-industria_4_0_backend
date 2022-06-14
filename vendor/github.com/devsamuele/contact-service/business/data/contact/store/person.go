package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/devsamuele/contact-service/business/data/contact/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Person struct {
	db  *mongo.Database
	log *log.Logger
}

func NewPerson(db *mongo.Database, log *log.Logger) Person {
	return Person{
		db:  db,
		log: log,
	}
}

// Field

func (p Person) QueryField(ctx context.Context, tenantID primitive.ObjectID) ([]model.FieldWithSection, error) {
	return newFieldWithSection(p.db, p.log, personFieldCollection).query(ctx, tenantID)
}

func (p Person) QueryFieldByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) (model.FieldWithSection, error) {
	return newFieldWithSection(p.db, p.log, personFieldCollection).queryByID(ctx, tenantID, fieldID)
}

func (p Person) QueryFieldBySectionID(ctx context.Context, tenantID, sectionID primitive.ObjectID) ([]model.FieldWithSection, error) {
	return newFieldWithSection(p.db, p.log, personFieldCollection).queryBySectionID(ctx, tenantID, sectionID)
}

func (p Person) QueryFieldByName(ctx context.Context, tenantID primitive.ObjectID, name string) (model.FieldWithSection, error) {
	return newFieldWithSection(p.db, p.log, personFieldCollection).queryByName(ctx, tenantID, name)
}

func (p Person) QueryFieldByLabel(ctx context.Context, tenantID primitive.ObjectID, label model.Translation) (model.FieldWithSection, error) {
	return newFieldWithSection(p.db, p.log, personFieldCollection).queryByLabel(ctx, tenantID, label)
}

func (p Person) CheckFieldByLabel(ctx context.Context, tenantID primitive.ObjectID, label model.Translation) error {
	return newFieldWithSection(p.db, p.log, personFieldCollection).checkByLabel(ctx, tenantID, label)
}

func (p Person) InsertField(ctx context.Context, nf model.FieldWithSection) error {
	return newFieldWithSection(p.db, p.log, personFieldCollection).insert(ctx, nf)
}

func (p Person) BulkInsertField(ctx context.Context, nfs []model.FieldWithSection) error {
	return newFieldWithSection(p.db, p.log, personFieldCollection).bulkInsert(ctx, nfs)
}

func (p Person) UpdateField(ctx context.Context, field model.FieldWithSection) error {
	return newFieldWithSection(p.db, p.log, personFieldCollection).update(ctx, field)
}

func (p Person) DeleteFieldByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) error {
	return newFieldWithSection(p.db, p.log, personFieldCollection).deleteByID(ctx, tenantID, fieldID)
}

func (p Person) DeleteFieldByIDs(ctx context.Context, tenantID primitive.ObjectID, fieldIDs []primitive.ObjectID) error {
	return newFieldWithSection(p.db, p.log, personFieldCollection).deleteByIDs(ctx, tenantID, fieldIDs)
}

func (p Person) DeleteFieldBySectionID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {
	return newFieldWithSection(p.db, p.log, personFieldCollection).deleteBySectionID(ctx, tenantID, sectionID)
}

func (p Person) DeleteFieldBySectionIDs(ctx context.Context, tenantID primitive.ObjectID, sectionIDs []primitive.ObjectID) error {
	return newFieldWithSection(p.db, p.log, personFieldCollection).deleteBySectionIDs(ctx, tenantID, sectionIDs)
}

func (p Person) DeleteField(ctx context.Context, tenantID primitive.ObjectID) error {
	return newFieldWithSection(p.db, p.log, personFieldCollection).delete(ctx, tenantID)
}

// Person

func (p Person) Count(ctx context.Context, tenantID primitive.ObjectID) (int64, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	count, err := p.db.Collection(personCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, err
}

func (p Person) QueryResponseByID(ctx context.Context, tenantID, personID primitive.ObjectID) (model.PersonResponse, error) {

	var pipeline mongo.Pipeline

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: personID}}}}
	pipeline = append(pipeline, matchStage)

	lookupOrganizationStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "organization"},
		{Key: "localField", Value: "organization_id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "organization"},
	}}}

	unwindOrganizationStage := bson.D{{
		Key: "$unwind",
		Value: bson.D{{
			Key:   "path",
			Value: "$organization",
		}, {
			Key:   "preserveNullAndEmptyArrays",
			Value: true,
		}},
	}}

	lookupOfficeStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "office"},
		{Key: "localField", Value: "office_id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "office"},
	}}}

	unwindOfficeStage := bson.D{{
		Key: "$unwind",
		Value: bson.D{{
			Key:   "path",
			Value: "$office",
		}, {
			Key:   "preserveNullAndEmptyArrays",
			Value: true,
		}},
	}}

	pipeline = append(pipeline, lookupOrganizationStage, unwindOrganizationStage, lookupOfficeStage, unwindOfficeStage)

	cur, err := p.db.Collection(personCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return model.PersonResponse{}, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, ctx)

	persons := make([]model.PersonResponse, 0)
	if err := cur.All(ctx, &persons); err != nil {
		return model.PersonResponse{}, err
	}

	if len(persons) == 0 {
		return model.PersonResponse{}, ErrNotFound
	}

	return persons[0], nil
}

func (p Person) QueryByID(ctx context.Context, tenantID, personID primitive.ObjectID) (model.Person, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: personID}}

	res := p.db.Collection(personCollection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Person{}, ErrNotFound
		}
		return model.Person{}, err
	}

	var person model.Person
	if err := res.Decode(&person); err != nil {
		return model.Person{}, err
	}

	return person, nil
}

func (p Person) Insert(ctx context.Context, person model.Person) error {

	if _, err := p.db.Collection(personCollection).InsertOne(ctx, person); err != nil {
		return err
	}

	return nil
}

func (p Person) BulkInsert(ctx context.Context, persons []model.Person) error {

	_persons := make([]interface{}, 0)
	for _, person := range persons {
		_persons = append(_persons, person)
	}

	if _, err := p.db.Collection(personCollection).InsertMany(ctx, _persons); err != nil {
		return err
	}

	return nil

}

func (p Person) Update(ctx context.Context, person model.Person) error {

	filter := bson.D{{Key: "tenant_id", Value: person.TenantID}, {Key: "_id", Value: person.ID}}
	update := bson.D{{Key: "$set", Value: person}}

	_, err := p.db.Collection(personCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (p Person) AddFieldValue(ctx context.Context, tenantID, fieldID primitive.ObjectID, value interface{}, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: value}, {Key: "updated", Value: now}}}}

	_, err := p.db.Collection(personCollection).UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (p Person) UpdateFieldValue(ctx context.Context, tenantID primitive.ObjectID, fieldID primitive.ObjectID, value interface{}, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: value}, {Key: "updated", Value: now}}}}

	if _, err := p.db.Collection(personCollection).UpdateMany(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (p Person) DeleteByIDs(ctx context.Context, tenantID primitive.ObjectID, personIDs []primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: bson.D{{Key: "$in", Value: personIDs}}}}
	_, err := p.db.Collection(personCollection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (p Person) DeleteByID(ctx context.Context, tenantID, personID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: personID}}
	_, err := p.db.Collection(personCollection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (p Person) TrashByID(ctx context.Context, tenantID, personID primitive.ObjectID, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: personID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "deleted", Value: true},
		{Key: "date_of_delete", Value: now},
	}}}
	_, err := p.db.Collection(personCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (p Person) TrashByIDs(ctx context.Context, tenantID primitive.ObjectID, personIDs []primitive.ObjectID, now time.Time) error {

	models := make([]mongo.WriteModel, 0)
	for _, personID := range personIDs {
		filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: personID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: true}, {Key: "date_of_delete", Value: now}}}}

		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}

	_, err := p.db.Collection(personCollection).BulkWrite(ctx, models)
	if err != nil {
		return err
	}

	return nil
}

func (p Person) TrashAll(ctx context.Context, tenantID primitive.ObjectID, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: true}, {Key: "date_of_delete", Value: now}}}}

	_, err := p.db.Collection(personCollection).UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (p Person) RestoreByID(ctx context.Context, tenantID, personID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: personID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: false}, {Key: "date_of_delete", Value: nil}}}}
	_, err := p.db.Collection(personCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (p Person) RestoreByIDs(ctx context.Context, tenantID primitive.ObjectID, personIDs []primitive.ObjectID) error {

	models := make([]mongo.WriteModel, 0)
	for _, personID := range personIDs {

		filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: personID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: false}, {Key: "date_of_delete", Value: nil}}}}

		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}

	_, err := p.db.Collection(personCollection).BulkWrite(ctx, models)
	if err != nil {
		return err
	}

	return nil
}

func (p Person) RestoreAll(ctx context.Context, tenantID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: false}, {Key: "date_of_delete", Value: nil}}}}

	_, err := p.db.Collection(personCollection).UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (p Person) UnsetFieldValues(ctx context.Context, tenantID primitive.ObjectID, fieldIDs []primitive.ObjectID, now time.Time) error {

	models := make([]mongo.WriteModel, 0)
	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	for _, fieldID := range fieldIDs {
		update := bson.D{{Key: "$unset", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: ""}, {Key: "updated", Value: now}}}}
		models = append(models, mongo.NewUpdateManyModel().SetFilter(filter).SetUpdate(update))
	}

	if _, err := p.db.Collection(personCollection).BulkWrite(ctx, models); err != nil {
		return err
	}

	return nil
}

func (p Person) UnsetAllFieldValues(ctx context.Context, tenantID primitive.ObjectID, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "field_values", Value: bson.D{}}, {Key: "updated", Value: now}}}}

	if _, err := p.db.Collection(personCollection).UpdateMany(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (p Person) PullChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error {

	// make bulk
	for _, choiceID := range choiceIDs {
		filter := bson.D{{Key: "tenant_id", Value: tenantID}}
		update := bson.D{
			{Key: "$pull", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: choiceID}}},
			{Key: "$set", Value: bson.D{{Key: "updated", Value: now}}}}
		if _, err := p.db.Collection(personCollection).UpdateMany(ctx, filter, update); err != nil {
			return err
		}
	}
	return nil
}

func (p Person) RemoveChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error {

	// make bulk
	for _, choiceID := range choiceIDs {
		filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: choiceID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: nil}, {Key: "updated", Value: now}}}}
		if _, err := p.db.Collection(personCollection).UpdateMany(ctx, filter, update); err != nil {
			return err
		}
	}

	return nil
}

func (p Person) Search(ctx context.Context, tenantID primitive.ObjectID, filter model.Filter, pageNumber, itemsPerPage int) ([]model.PersonResponse, error) {

	var pipeline mongo.Pipeline

	if filter.Sort.Field == "" {
		filter.Sort.Field = "name"
	}

	if filter.Sort.Order == 0 {
		filter.Sort.Order = 1
	}

	matchTenantStage := bson.D{{Key: "$match", Value: bson.D{{Key: "tenant_id", Value: tenantID}}}}
	limitStage := bson.D{{Key: "$limit", Value: itemsPerPage}}
	skipStage := bson.D{{Key: "$skip", Value: (pageNumber - 1) * itemsPerPage}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: filter.Sort.Field, Value: filter.Sort.Order}}}}

	mergeEmailStage := bson.D{{
		Key: "$set",
		Value: bson.D{{
			Key: "email",
			Value: bson.D{{
				Key:   "$concatArrays",
				Value: bson.A{"$others_email", bson.A{"$primary_email"}},
			}},
		}},
	}}

	pipeline = append(pipeline, matchTenantStage, mergeEmailStage)

	lookupOrganizationStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "organization"},
		{Key: "localField", Value: "organization_id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "organization"},
	}}}

	unwindOrganizationStage := bson.D{{
		Key: "$unwind",
		Value: bson.D{{
			Key:   "path",
			Value: "$organization",
		}, {
			Key:   "preserveNullAndEmptyArrays",
			Value: true,
		}},
	}}

	lookupOfficeStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "office"},
		{Key: "localField", Value: "office_id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "office"},
	}}}

	unwindOfficeStage := bson.D{{
		Key: "$unwind",
		Value: bson.D{{
			Key:   "path",
			Value: "$office",
		}, {
			Key:   "preserveNullAndEmptyArrays",
			Value: true,
		}},
	}}

	pipeline = append(pipeline, lookupOrganizationStage, unwindOrganizationStage, lookupOfficeStage, unwindOfficeStage)

	var OR bson.A
	for _, ANDFields := range filter.ORFields {
		var AND bson.A
		for _, filterField := range ANDFields.ANDFields {

			filter, err := filterField.BuildBSONFilter("person")
			if err != nil {
				return nil, err
			}
			AND = append(AND, filter)
		}
		OR = append(OR, bson.D{{Key: "$and", Value: AND}})
	}

	matchFilterStage := bson.D{{Key: "$match", Value: bson.D{{Key: "$or", Value: OR}}}}
	pipeline = append(pipeline, matchFilterStage, sortStage, skipStage, limitStage)

	cur, err := p.db.Collection(personCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, ctx)

	persons := make([]model.PersonResponse, 0)
	if err := cur.All(ctx, &persons); err != nil {
		return nil, err
	}

	return persons, nil
}

// Section

func (p Person) QuerySection(ctx context.Context, tenantID primitive.ObjectID) ([]model.Section, error) {
	return newSection(p.db, p.log, personSectionCollection).query(ctx, tenantID)
}

func (p Person) QuerySectionByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) (model.Section, error) {
	return newSection(p.db, p.log, personSectionCollection).queryByID(ctx, tenantID, sectionID)
}

func (p Person) CheckSectionByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {
	return newSection(p.db, p.log, personSectionCollection).checkByID(ctx, tenantID, sectionID)
}

func (p Person) QuerySectionByName(ctx context.Context, tenantID primitive.ObjectID, name string) (model.Section, error) {
	return newSection(p.db, p.log, personSectionCollection).queryByName(ctx, tenantID, name)
}

func (p Person) QuerySectionByTitle(ctx context.Context, tenantID primitive.ObjectID, title model.Translation) (model.Section, error) {
	return newSection(p.db, p.log, personSectionCollection).queryByTitle(ctx, tenantID, title)
}

func (p Person) InsertSection(ctx context.Context, ns model.Section) error {
	return newSection(p.db, p.log, personSectionCollection).insert(ctx, ns)
}

func (p Person) BulkInsertSection(ctx context.Context, nss []model.Section) error {
	return newSection(p.db, p.log, personSectionCollection).bulkInsert(ctx, nss)
}

func (p Person) UpdateSection(ctx context.Context, us model.Section) error {
	return newSection(p.db, p.log, personSectionCollection).update(ctx, us)
}

func (p Person) DeleteSectionByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {
	return newSection(p.db, p.log, personSectionCollection).deleteByID(ctx, tenantID, sectionID)
}

func (p Person) DeleteSectionByIDs(ctx context.Context, tenantID primitive.ObjectID, sectionIDs []primitive.ObjectID) error {
	return newSection(p.db, p.log, personSectionCollection).deleteByIDs(ctx, tenantID, sectionIDs)
}

func (p Person) DeleteSection(ctx context.Context, tenantID primitive.ObjectID) error {
	return newSection(p.db, p.log, personSectionCollection).delete(ctx, tenantID)
}
