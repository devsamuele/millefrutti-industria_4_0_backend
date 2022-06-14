package store

/*

import (
	"context"
	"github.com/devsamuele/contact-service/business/data/contact/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type PersonRole struct {
	db  *mongo.Database
	log *log.Logger
}

func NewPersonRole(db *mongo.Database, log *log.Logger) PersonRole {
	return PersonRole{
		db:  db,
		log: log,
	}
}

func (o PersonRole) Query(ctx context.Context, tenantID primitive.ObjectID) ([]model.PersonRole, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	cur, err := o.db.Collection(personRoleCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, ctx)

	personRoles := make([]model.PersonRole, 0)
	if err := cur.All(ctx, &personRoles); err != nil {
		return nil, err
	}

	return personRoles, nil
}

func (o PersonRole) QueryByID(ctx context.Context, tenantID, personRoleID primitive.ObjectID) (model.PersonRole, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: personRoleID}}
	res := o.db.Collection(personRoleCollection).FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.PersonRole{}, ErrNotFound
		}
		return model.PersonRole{}, err
	}

	var personRole model.PersonRole
	if err := res.Decode(&personRole); err != nil {
		return model.PersonRole{}, err
	}

	return personRole, nil
}

func (o PersonRole) Count(ctx context.Context, tenantID primitive.ObjectID) (int64, error) {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}}
	count, err := o.db.Collection(personRoleCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, err
}

func (o PersonRole) Insert(ctx context.Context, personRole model.PersonRole) error {

	if _, err := o.db.Collection(personRoleCollection).InsertOne(ctx, personRole); err != nil {
		return err
	}

	return nil
}

func (o PersonRole) BulkInsert(ctx context.Context, personRoles []model.PersonRole) error {

	_personRoles := make([]interface{}, 0)
	for _, personRole := range personRoles {
		_personRoles = append(_personRoles, personRole)
	}

	if _, err := o.db.Collection(personRoleCollection).InsertMany(ctx, _personRoles); err != nil {
		return err
	}

	return nil

}

func (o PersonRole) Update(ctx context.Context, personRole model.PersonRole) error {

	filter := bson.D{{Key: "tenant_id", Value: personRole.TenantID}, {Key: "_id", Value: personRole.ID}}
	update := bson.D{{Key: "$set", Value: personRole}}

	_, err := o.db.Collection(personRoleCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (o PersonRole) DeleteByIDs(ctx context.Context, tenantID primitive.ObjectID, personRoleIDs []primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: bson.D{{Key: "$in", Value: personRoleIDs}}}}
	_, err := o.db.Collection(personRoleCollection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (o PersonRole) DeleteByID(ctx context.Context, tenantID, personRoleID primitive.ObjectID) error {

	filter := bson.D{{Key: "tenant_id", Value: tenantID}, {Key: "_id", Value: personRoleID}}
	_, err := o.db.Collection(personRoleCollection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
*/
