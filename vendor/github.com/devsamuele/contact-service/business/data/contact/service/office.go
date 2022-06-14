package service

import (
	"context"
	"errors"
	"github.com/devsamuele/contact-service/business/data/contact/model"
	"github.com/devsamuele/contact-service/business/data/contact/store"
	"github.com/devsamuele/contact-service/business/sys/database"
	"github.com/devsamuele/service-kit/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Office struct {
	organizationStore store.Organization
	officeStore       store.Office
	session           database.Session
}

func NewOffice(session database.Session, organizationStore store.Organization, officeStore store.Office) Office {
	return Office{
		organizationStore: organizationStore,
		officeStore:       officeStore,
		session:           session,
	}
}

// Field

func (s Office) QueryField(ctx context.Context, tenantID string) ([]model.Field, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, web.NewShutdownError("invalid tenant_id")
		}
	}

	fields, err := s.officeStore.QueryField(ctx, tenantOID)
	if err != nil {
		return nil, err
	}

	return fields, err

}

func (s Office) QueryFieldByID(ctx context.Context, tenantID string, fieldID string) (model.Field, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Field{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	fieldOID, err := primitive.ObjectIDFromHex(fieldID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Field{}, model.Error{
				Message:      "ID is not in its proper form",
				Reason:       model.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	field, err := s.officeStore.QueryFieldByID(ctx, tenantOID, fieldOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Field{}, model.Error{
				Message:      "field not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "field_id",
			}
		}
		return model.Field{}, err
	}

	return field, err
}

func (s Office) InsertField(ctx context.Context, tenantID string, nf model.NewField, now time.Time) (model.Field, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Field{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	err = s.officeStore.CheckFieldByLabel(ctx, tenantOID, *nf.Label)
	if err != nil {
		if !errors.Is(err, store.ErrNotFound) {
			return model.Field{}, err
		}
	} else {
		return model.Field{}, model.Error{
			Message:      "label already exist",
			Reason:       model.ErrReasonConflict,
			LocationType: "argument",
			Location:     "label",
		}
	}

	if err := nf.Validate(model.OfficeFieldTypes()); err != nil {
		return model.Field{}, err
	}

	field := model.Field{
		ID:            primitive.NewObjectID(),
		TenantID:      tenantOID,
		Name:          "",
		Type:          *nf.Type,
		Label:         *nf.Label,
		Visible:       true,
		Searchable:    true,
		EditableValue: true,
		Created:       now,
		Updated:       now,
	}

	if nf.Searchable != nil {
		field.Searchable = *nf.Searchable
	}

	field.Choices = nf.BuildChoices()
	field.Options = nf.BuildOptions()

	// ---------------------------------------------------------------------------------

	sess, err := s.session.Start()
	if err != nil {
		return model.Field{}, err
	}

	defer sess.EndSession(ctx)

	txCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {

		err = s.officeStore.InsertField(ctx, field)
		if err != nil {
			return model.Field{}, err
		}

		count, err := s.officeStore.Count(ctx, tenantOID)
		if err != nil {
			return model.Field{}, err
		}

		if count > 0 {
			var value interface{}
			if model.IsFieldWithArrayValue(field.Type) {
				value = bson.A{}
			}

			if err := s.officeStore.AddFieldValue(ctx, tenantOID, field.ID, value, now); err != nil {
				return model.Field{}, err
			}
		}
		return nil, nil
	}

	if _, err := sess.WithTransaction(ctx, txCallback); err != nil {
		return model.Field{}, err
	}

	return field, nil
}

func (s Office) UpdateField(ctx context.Context, tenantID string, fieldID string, uf model.UpdateField, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	fieldOID, err := primitive.ObjectIDFromHex(fieldID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Error{
				Message:      "ID is not in its proper form",
				Reason:       model.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	oldField, err := s.officeStore.QueryFieldByID(ctx, tenantOID, fieldOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Error{
				Message:      "field not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return err
	}

	if err := s.officeStore.CheckFieldByLabel(ctx, tenantOID, *uf.Label); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Error{
				Message:      "label already exist",
				Reason:       model.ErrReasonConflict,
				LocationType: "argument",
				Location:     "label",
			}
		}
		return err
	}

	if err := uf.Validate(oldField); err != nil {
		return err
	}

	if uf.Label != nil {
		oldField.Label = *uf.Label
	}

	if uf.Searchable != nil {
		oldField.Searchable = *uf.Searchable
	}

	updatedChoices, removedChoiceIDs := uf.UpdateChoices(oldField)
	oldField.Choices = updatedChoices

	oldField.Options = uf.UpdateOptions(oldField)
	oldField.Updated = now

	sess, err := s.session.Start()
	if err != nil {
		return err
	}

	defer sess.EndSession(ctx)

	txCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {

		if len(removedChoiceIDs) > 0 {
			switch oldField.Type {
			case model.FieldTypeSingleSelect:
				if len(removedChoiceIDs) > 0 {
					if err := s.officeStore.RemoveChoices(sessCtx, tenantOID, fieldOID, removedChoiceIDs, now); err != nil {
						return nil, err
					}
				}

			case model.FieldTypeMultipleSelect:
				if len(removedChoiceIDs) > 0 {
					if err := s.officeStore.PullChoices(sessCtx, tenantOID, fieldOID, removedChoiceIDs, now); err != nil {
						return nil, err
					}
				}
			}
		}

		if err := s.officeStore.UpdateField(ctx, oldField); err != nil {
			return nil, err
		}

		return nil, nil
	}

	if _, err := sess.WithTransaction(ctx, txCallback); err != nil {
		return err
	}

	return nil
}

func (s Office) DeleteFieldByID(ctx context.Context, tenantID, fieldID string, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	fieldOID, err := primitive.ObjectIDFromHex(fieldID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Error{
				Message:      "ID is not in its proper form",
				Reason:       model.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	_, err = s.officeStore.QueryFieldByID(ctx, tenantOID, fieldOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Error{}
		}
		return err
	}

	sess, err := s.session.Start()
	if err != nil {
		return err
	}

	txCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		err = s.officeStore.DeleteFieldByID(sessCtx, tenantOID, fieldOID)
		if err != nil {
			return nil, err
		}

		err = s.officeStore.UnsetFieldValues(sessCtx, tenantOID, []primitive.ObjectID{fieldOID}, now)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	_, err = sess.WithTransaction(ctx, txCallback)
	if err != nil {
		return err
	}

	return nil
}

func (s Office) DeleteFieldByIDs(ctx context.Context, tenantID string, fieldIDs model.ObjectIDs, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	sess, err := s.session.Start()
	if err != nil {
		return err
	}

	txCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		err = s.officeStore.DeleteFieldByIDs(sessCtx, tenantOID, fieldIDs.IDs)
		if err != nil {
			return nil, err
		}

		err = s.officeStore.UnsetFieldValues(sessCtx, tenantOID, fieldIDs.IDs, now)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	_, err = sess.WithTransaction(ctx, txCallback)
	if err != nil {
		return err
	}

	return nil
}

// Office

func (s Office) QueryByID(ctx context.Context, tenantID, officeID string) (model.Office, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Office{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	officeOID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Office{}, model.Error{
				Message:      "ID is not in its proper form",
				Reason:       model.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	office, err := s.officeStore.QueryByID(ctx, tenantOID, officeOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Office{}, model.Error{
				Message:      "office not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return model.Office{}, err
	}

	return office, nil
}

func (s Office) Create(ctx context.Context, tenantID string, no model.NewOffice, now time.Time) (model.Office, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Office{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	if no.OrganizationID == nil {
		return model.Office{}, model.Error{
			Message:      "organization_id is required",
			Reason:       model.ErrReasonRequired,
			LocationType: "argument",
			Location:     "organization_id",
		}
	}

	_, err = s.organizationStore.QueryByID(ctx, tenantOID, *no.OrganizationID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Office{}, model.Error{
				Message:      "organization not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "argument",
				Location:     "organization_id",
			}
		}
	}

	fields, err := s.officeStore.QueryField(ctx, tenantOID)
	if err != nil {
		return model.Office{}, err
	}

	decodedFieldValues, err := no.FieldValues.Validate(fields)
	if err != nil {
		return model.Office{}, err
	}

	newOffice := model.Office{
		ID:             primitive.NewObjectID(),
		TenantID:       tenantOID,
		FieldValues:    decodedFieldValues,
		OrganizationID: *no.OrganizationID,
		Created:        now,
		Updated:        now,
	}

	err = s.officeStore.Insert(ctx, newOffice)
	if err != nil {
		return model.Office{}, err
	}

	// orgBytes, err := json.Marshal(org)
	// if err != nil {
	// 	return organization.Organization{}, err
	// }

	// type AmqpMessage struct {
	// 	ID       string             `json:"id"`
	// 	TenantID primitive.ObjectID `json:"tenant_id"`
	// 	Resource string             `json:"resource"`
	// 	Action   string             `json:"action"`
	// 	Data     []byte             `json:"data"`
	// }

	// m := AmqpMessage{
	// 	ID:       uuid.NewString(),
	// 	TenantID: tenantID,
	// 	Resource: "organization",
	// 	Action:   "create",
	// 	Data:     orgBytes,
	// }

	// mBytes, err := json.Marshal(m)
	// if err != nil {
	// 	return organization.Organization{}, err
	// }

	// s.ch.Publish("service_message", "", true, false, amqp.Publishing{
	// 	DeliveryMode: amqp.Persistent,
	// 	ContentType:  "application/json",
	// 	MessageId:    uuid.NewString(),
	// 	Body:         mBytes,
	// })

	return newOffice, nil
}

func (s Office) Update(ctx context.Context, tenantID, officeID string, uo model.UpdateOffice, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	officeOID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Error{
				Message:      "ID is not in its proper form",
				Reason:       model.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	office, err := s.officeStore.QueryByID(ctx, tenantOID, officeOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Error{
				Message:      "office not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return err
	}

	fields, err := s.officeStore.QueryField(ctx, tenantOID)
	if err != nil {
		return err
	}

	decodedFieldValues, err := uo.FieldValues.Validate(fields)
	if err != nil {
		return err
	}

	for fieldID, fieldValue := range decodedFieldValues {
		office.FieldValues[fieldID] = fieldValue
	}

	office.Updated = now

	err = s.officeStore.Update(ctx, office)
	if err != nil {
		return err
	}

	return nil
}

func (s Office) Delete(ctx context.Context, tenantID, officeID string) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	officeOID, err := primitive.ObjectIDFromHex(officeID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Error{
				Message:      "ID is not in its proper form",
				Reason:       model.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	_, err = s.officeStore.QueryByID(ctx, tenantOID, officeOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Error{
				Message:      "office not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return err
	}

	err = s.officeStore.DeleteByID(ctx, tenantOID, officeOID)
	if err != nil {
		return err
	}

	return nil
}
