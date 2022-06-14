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

type Organization struct {
	db  *mongo.Database
	log *log.Logger
}

func NewOrganization(db *mongo.Database, log *log.Logger) Organization {
	return Organization{
		db:  db,
		log: log,
	}
}

// Field

func (o Organization) QueryField(ctx context.Context, tenantID primitive.ObjectID) ([]model.FieldWithSection, error) {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).query(ctx, tenantID)
}

func (o Organization) QueryFieldByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) (model.FieldWithSection, error) {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).queryByID(ctx, tenantID, fieldID)
}

func (o Organization) QueryFieldBySectionID(ctx context.Context, tenantID, sectionID primitive.ObjectID) ([]model.FieldWithSection, error) {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).queryBySectionID(ctx, tenantID, sectionID)
}

func (o Organization) QueryFieldByName(ctx context.Context, tenantID primitive.ObjectID, name string) (model.FieldWithSection, error) {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).queryByName(ctx, tenantID, name)
}

func (o Organization) QueryFieldByLabel(ctx context.Context, tenantID primitive.ObjectID, label model.Translation) (model.FieldWithSection, error) {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).queryByLabel(ctx, tenantID, label)
}

func (o Organization) CheckFieldByLabel(ctx context.Context, tenantID primitive.ObjectID, label model.Translation) error {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).checkByLabel(ctx, tenantID, label)
}

func (o Organization) InsertField(ctx context.Context, nf model.FieldWithSection) error {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).insert(ctx, nf)
}

func (o Organization) BulkInsertField(ctx context.Context, nfs []model.FieldWithSection) error {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).bulkInsert(ctx, nfs)
}

func (o Organization) UpdateField(ctx context.Context, field model.FieldWithSection) error {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).update(ctx, field)
}

func (o Organization) DeleteFieldByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) error {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).deleteByID(ctx, tenantID, fieldID)
}

func (o Organization) DeleteFieldByIDs(ctx context.Context, tenantID primitive.ObjectID, fieldIDs []primitive.ObjectID) error {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).deleteByIDs(ctx, tenantID, fieldIDs)
}

func (o Organization) DeleteFieldBySectionID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).deleteBySectionID(ctx, tenantID, sectionID)
}

func (o Organization) DeleteFieldBySectionIDs(ctx context.Context, tenantID primitive.ObjectID, sectionIDs []primitive.ObjectID) error {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).deleteBySectionIDs(ctx, tenantID, sectionIDs)
}

func (o Organization) DeleteField(ctx context.Context, tenantID primitive.ObjectID) error {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).delete(ctx, tenantID)
}

// Organization

func (o Organization) Count(ctx context.Context, tenantID primitive.ObjectID) (int64, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	count, err := o.db.Collection(organizationCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, err
}

func (o Organization) QueryByID(ctx context.Context, tenantID, organizationID primitive.ObjectID) (model.Organization, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: organizationID}}
	res := o.db.Collection(organizationCollection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Organization{}, ErrNotFound
		}
		return model.Organization{}, err
	}

	var organization model.Organization
	if err := res.Decode(&organization); err != nil {
		return model.Organization{}, err
	}

	return organization, nil
}

func (o Organization) Insert(ctx context.Context, organization model.Organization) error {

	if _, err := o.db.Collection(organizationCollection).InsertOne(ctx, organization); err != nil {
		return err
	}

	return nil
}

func (o Organization) BulkInsert(ctx context.Context, organizations []model.Organization) error {

	_organizations := make([]interface{}, 0)
	for _, organization := range organizations {
		_organizations = append(_organizations, organization)
	}

	if _, err := o.db.Collection(organizationCollection).InsertMany(ctx, _organizations); err != nil {
		return err
	}

	return nil

}

func (o Organization) Update(ctx context.Context, organization model.Organization) error {

	filter := bson.D{{Key: "tenant_id", Value: organization.TenantID}, {Key: "_id", Value: organization.ID}}
	update := bson.D{{Key: "$set", Value: organization}}

	_, err := o.db.Collection(organizationCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (o Organization) AddFieldValue(ctx context.Context, tenantID, fieldID primitive.ObjectID, value interface{}, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: value}, {Key: "updated", Value: now}}}}

	_, err := o.db.Collection(organizationCollection).UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (o Organization) UpdateFieldValue(ctx context.Context, tenantID primitive.ObjectID, fieldID primitive.ObjectID, value interface{}, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: value}, {Key: "updated", Value: now}}}}

	if _, err := o.db.Collection(organizationCollection).UpdateMany(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (o Organization) DeleteByIDs(ctx context.Context, tenantID primitive.ObjectID, organizationIDs []primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: bson.D{{Key: "$in", Value: organizationIDs}}}}
	_, err := o.db.Collection(organizationCollection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (o Organization) DeleteByID(ctx context.Context, tenantID, organizationID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: organizationID}}
	_, err := o.db.Collection(organizationCollection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (o Organization) TrashByID(ctx context.Context, tenantID, organizationID primitive.ObjectID, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: organizationID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "deleted", Value: true},
		{Key: "date_of_delete", Value: now},
	}}}
	_, err := o.db.Collection(organizationCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (o Organization) TrashByIDs(ctx context.Context, tenantID primitive.ObjectID, organizationIDs []primitive.ObjectID, now time.Time) error {

	models := make([]mongo.WriteModel, 0)
	for _, organizationID := range organizationIDs {
		filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: organizationID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: true}, {Key: "date_of_delete", Value: now}}}}

		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}

	_, err := o.db.Collection(organizationCollection).BulkWrite(ctx, models)
	if err != nil {
		return err
	}

	return nil
}

func (o Organization) TrashAll(ctx context.Context, tenantID primitive.ObjectID, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: true}, {Key: "date_of_delete", Value: now}}}}

	_, err := o.db.Collection(organizationCollection).UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (o Organization) RestoreByID(ctx context.Context, tenantID, organizationID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: organizationID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: false}, {Key: "date_of_delete", Value: nil}}}}
	_, err := o.db.Collection(organizationCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (o Organization) RestoreByIDs(ctx context.Context, tenantID primitive.ObjectID, organizationIDs []primitive.ObjectID) error {

	models := make([]mongo.WriteModel, 0)
	for _, organizationID := range organizationIDs {

		filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: organizationID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: false}, {Key: "date_of_delete", Value: nil}}}}

		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}

	_, err := o.db.Collection(organizationCollection).BulkWrite(ctx, models)
	if err != nil {
		return err
	}

	return nil
}

func (o Organization) RestoreAll(ctx context.Context, tenantID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: false}, {Key: "date_of_delete", Value: nil}}}}

	_, err := o.db.Collection(organizationCollection).UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (o Organization) UnsetFieldValues(ctx context.Context, tenantID primitive.ObjectID, fieldIDs []primitive.ObjectID, now time.Time) error {

	models := make([]mongo.WriteModel, 0)
	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	for _, fieldID := range fieldIDs {
		update := bson.D{{Key: "$unset", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: ""}, {Key: "updated", Value: now}}}}
		models = append(models, mongo.NewUpdateManyModel().SetFilter(filter).SetUpdate(update))
	}

	if _, err := o.db.Collection(organizationCollection).BulkWrite(ctx, models); err != nil {
		return err
	}

	return nil
}

func (o Organization) UnsetAllFieldValues(ctx context.Context, tenantID primitive.ObjectID, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "field_values", Value: bson.D{}}, {Key: "updated", Value: now}}}}

	if _, err := o.db.Collection(organizationCollection).UpdateMany(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (o Organization) PullChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error {

	// make bulk
	for _, choiceID := range choiceIDs {
		filter := bson.D{{Key: "tenant_id", Value: tenantID}}
		update := bson.D{
			{Key: "$pull", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: choiceID}}},
			{Key: "$set", Value: bson.D{{Key: "updated", Value: now}}}}
		if _, err := o.db.Collection(organizationCollection).UpdateMany(ctx, filter, update); err != nil {
			return err
		}
	}
	return nil
}

func (o Organization) RemoveChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error {

	// make bulk
	for _, choiceID := range choiceIDs {
		filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: choiceID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: nil}, {Key: "updated", Value: now}}}}
		if _, err := o.db.Collection(organizationCollection).UpdateMany(ctx, filter, update); err != nil {
			return err
		}
	}

	return nil
}

func (o Organization) Search(ctx context.Context, tenantID primitive.ObjectID, filter model.Filter, pageNumber, itemsPerPage int) ([]model.Organization, error) {

	var pipeline mongo.Pipeline

	matchTenantStage := bson.D{{Key: "$match", Value: bson.D{{Key: "tenant_id", Value: tenantID}}}}
	limitStage := bson.D{{Key: "$limit", Value: itemsPerPage}}
	skipStage := bson.D{{Key: "$skip", Value: (pageNumber - 1) * itemsPerPage}}

	pipeline = append(pipeline, matchTenantStage)

	// TODO involvedResources
	lookupOfficeStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "office"},
		{Key: "localField", Value: "_id"},
		{Key: "foreignField", Value: "organization_id"},
		{Key: "as", Value: "office"},
	}}}
	pipeline = append(pipeline, lookupOfficeStage)

	var OR bson.A
	for _, ANDFields := range filter.ORFields {
		var AND bson.A
		for _, filterField := range ANDFields.ANDFields {

			filter, err := filterField.BuildBSONFilter("organization")
			if err != nil {
				return nil, err
			}
			AND = append(AND, filter)
		}
		OR = append(OR, bson.D{{Key: "$and", Value: AND}})
	}

	matchFilterStage := bson.D{{Key: "$match", Value: bson.D{{Key: "$or", Value: OR}}}}
	pipeline = append(pipeline, matchFilterStage, skipStage, limitStage)

	cur, err := o.db.Collection(organizationCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, ctx)

	organizations := make([]model.Organization, 0)
	if err := cur.All(ctx, &organizations); err != nil {
		return nil, err
	}

	return organizations, nil
}

// Section

func (o Organization) QuerySection(ctx context.Context, tenantID primitive.ObjectID) ([]model.Section, error) {
	return newSection(o.db, o.log, organizationSectionCollection).query(ctx, tenantID)
}

func (o Organization) QuerySectionByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) (model.Section, error) {
	return newSection(o.db, o.log, organizationSectionCollection).queryByID(ctx, tenantID, sectionID)
}

func (o Organization) CheckSectionByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {
	return newSection(o.db, o.log, organizationSectionCollection).checkByID(ctx, tenantID, sectionID)
}

func (o Organization) QuerySectionByName(ctx context.Context, tenantID primitive.ObjectID, name string) (model.Section, error) {
	return newSection(o.db, o.log, organizationSectionCollection).queryByName(ctx, tenantID, name)
}

func (o Organization) QuerySectionByTitle(ctx context.Context, tenantID primitive.ObjectID, title model.Translation) (model.Section, error) {
	return newSection(o.db, o.log, organizationSectionCollection).queryByTitle(ctx, tenantID, title)
}

func (o Organization) InsertSection(ctx context.Context, ns model.Section) error {
	return newSection(o.db, o.log, organizationSectionCollection).insert(ctx, ns)
}

func (o Organization) BulkInsertSection(ctx context.Context, nss []model.Section) error {
	return newSection(o.db, o.log, organizationSectionCollection).bulkInsert(ctx, nss)
}

func (o Organization) UpdateSection(ctx context.Context, us model.Section) error {
	return newSection(o.db, o.log, organizationSectionCollection).update(ctx, us)
}

func (o Organization) DeleteSectionByID(ctx context.Context, tenantID, sectionID primitive.ObjectID) error {
	return newSection(o.db, o.log, organizationSectionCollection).deleteByID(ctx, tenantID, sectionID)
}

func (o Organization) DeleteSectionByIDs(ctx context.Context, tenantID primitive.ObjectID, sectionIDs []primitive.ObjectID) error {
	return newSection(o.db, o.log, organizationSectionCollection).deleteByIDs(ctx, tenantID, sectionIDs)
}

func (o Organization) DeleteSection(ctx context.Context, tenantID primitive.ObjectID) error {
	return newSection(o.db, o.log, organizationSectionCollection).delete(ctx, tenantID)
}
