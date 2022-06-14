package store

import (
	"context"
	"fmt"
	"github.com/devsamuele/contact-service/business/data/contact/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Office struct {
	db  *mongo.Database
	log *log.Logger
}

func NewOffice(db *mongo.Database, log *log.Logger) Office {
	return Office{
		db:  db,
		log: log,
	}
}

func (o Office) QueryField(ctx context.Context, tenantID primitive.ObjectID) ([]model.Field, error) {
	return newField(o.db, o.log, officeFieldCollection).query(ctx, tenantID)
}

func (o Office) QueryFieldByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) (model.Field, error) {
	return newField(o.db, o.log, officeFieldCollection).queryByID(ctx, tenantID, fieldID)
}

func (o Office) QueryFieldByName(ctx context.Context, tenantID primitive.ObjectID, name string) (model.Field, error) {
	return newField(o.db, o.log, officeFieldCollection).queryByName(ctx, tenantID, name)
}

func (o Office) QueryFieldByLabel(ctx context.Context, tenantID primitive.ObjectID, label model.Translation) (model.Field, error) {
	return newField(o.db, o.log, officeFieldCollection).queryByLabel(ctx, tenantID, label)
}

func (o Office) CheckFieldByLabel(ctx context.Context, tenantID primitive.ObjectID, label model.Translation) error {
	return newFieldWithSection(o.db, o.log, organizationFieldCollection).checkByLabel(ctx, tenantID, label)
}

func (o Office) InsertField(ctx context.Context, nf model.Field) error {
	return newField(o.db, o.log, officeFieldCollection).insert(ctx, nf)
}

func (o Office) BulkInsertField(ctx context.Context, nfs []model.Field) error {
	return newField(o.db, o.log, officeFieldCollection).bulkInsert(ctx, nfs)
}

func (o Office) UpdateField(ctx context.Context, uf model.Field) error {
	return newField(o.db, o.log, officeFieldCollection).update(ctx, uf)
}

func (o Office) DeleteFieldByID(ctx context.Context, tenantID, fieldID primitive.ObjectID) error {
	return newField(o.db, o.log, officeFieldCollection).deleteByID(ctx, tenantID, fieldID)
}

func (o Office) DeleteFieldByIDs(ctx context.Context, tenantID primitive.ObjectID, fieldIDs []primitive.ObjectID) error {
	return newField(o.db, o.log, officeFieldCollection).deleteByIDs(ctx, tenantID, fieldIDs)
}

func (o Office) DeleteField(ctx context.Context, tenantID primitive.ObjectID) error {
	return newField(o.db, o.log, officeFieldCollection).delete(ctx, tenantID)
}

// OFFICE

func (o Office) QueryByID(ctx context.Context, tenantID, officeID primitive.ObjectID) (model.Office, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: officeID}}
	res := o.db.Collection(officeCollection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Office{}, ErrNotFound
		}
		return model.Office{}, err
	}

	var office model.Office
	if err := res.Decode(&office); err != nil {
		return model.Office{}, err
	}

	return office, nil
}

func (o Office) QueryByOrganizationID(ctx context.Context, tenantID, organizationID primitive.ObjectID) ([]model.Office, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "organization_id", Value: organizationID}}
	cur, err := o.db.Collection(officeCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, ctx)

	offices := make([]model.Office, 0)
	if err := cur.All(ctx, &offices); err != nil {
		return nil, err
	}

	return offices, nil
}

func (o Office) Count(ctx context.Context, tenantID primitive.ObjectID) (int64, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	count, err := o.db.Collection(officeCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, err
}

func (o Office) Insert(ctx context.Context, office model.Office) error {

	if _, err := o.db.Collection(officeCollection).InsertOne(ctx, office); err != nil {
		return err
	}

	return nil
}

func (o Office) BulkInsert(ctx context.Context, offices []model.Office) error {

	_offices := make([]interface{}, 0)
	for _, office := range offices {
		_offices = append(_offices, office)
	}

	if _, err := o.db.Collection(officeCollection).InsertMany(ctx, _offices); err != nil {
		return err
	}

	return nil

}

func (o Office) Update(ctx context.Context, office model.Office) error {

	filter := bson.D{{Key: "tenant_id", Value: office.TenantID}, {Key: "_id", Value: office.ID}}
	update := bson.D{{Key: "$set", Value: office}}

	_, err := o.db.Collection(officeCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (o Office) AddFieldValue(ctx context.Context, tenantID, fieldID primitive.ObjectID, value interface{}, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: value}, {Key: "updated", Value: now}}}}

	_, err := o.db.Collection(officeCollection).UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (o Office) UpdateFieldValue(ctx context.Context, tenantID primitive.ObjectID, fieldID primitive.ObjectID, value interface{}, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: value}, {Key: "updated", Value: now}}}}

	if _, err := o.db.Collection(officeCollection).UpdateMany(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (o Office) DeleteByIDs(ctx context.Context, tenantID primitive.ObjectID, officeIDs []primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: bson.D{{Key: "$in", Value: officeIDs}}}}
	_, err := o.db.Collection(officeCollection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (o Office) DeleteByOrganizationID(ctx context.Context, tenantID primitive.ObjectID, organizationID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "organization_id", Value: organizationID}}
	_, err := o.db.Collection(officeCollection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (o Office) DeleteByID(ctx context.Context, tenantID, officeID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: officeID}}
	_, err := o.db.Collection(officeCollection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (o Office) UnsetFieldValues(ctx context.Context, tenantID primitive.ObjectID, fieldIDs []primitive.ObjectID, now time.Time) error {

	models := make([]mongo.WriteModel, 0)
	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	for _, fieldID := range fieldIDs {
		update := bson.D{{Key: "$unset", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: ""}, {Key: "updated", Value: now}}}}
		models = append(models, mongo.NewUpdateManyModel().SetFilter(filter).SetUpdate(update))
	}

	if _, err := o.db.Collection(officeCollection).BulkWrite(ctx, models); err != nil {
		return err
	}

	return nil
}

func (o Office) UnsetAllFieldValues(ctx context.Context, tenantID primitive.ObjectID, now time.Time) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "field_values", Value: bson.D{}}, {Key: "updated", Value: now}}}}

	if _, err := o.db.Collection(officeCollection).UpdateMany(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (o Office) PullChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error {

	// make bulk
	for _, choiceID := range choiceIDs {
		filter := bson.D{{Key: "tenant_id", Value: tenantID}}
		update := bson.D{
			{Key: "$pull", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: choiceID}}},
			{Key: "$set", Value: bson.D{{Key: "updated", Value: now}}}}
		if _, err := o.db.Collection(officeCollection).UpdateMany(ctx, filter, update); err != nil {
			return err
		}
	}
	return nil
}

func (o Office) RemoveChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error {

	// make bulk
	for _, choiceID := range choiceIDs {
		filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: choiceID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("field_values.%v", fieldID.Hex()), Value: nil}, {Key: "updated", Value: now}}}}
		if _, err := o.db.Collection(officeCollection).UpdateMany(ctx, filter, update); err != nil {
			return err
		}
	}

	return nil
}
