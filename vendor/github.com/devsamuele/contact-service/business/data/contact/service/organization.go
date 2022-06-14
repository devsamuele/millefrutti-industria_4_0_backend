package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/devsamuele/contact-service/business/data/contact/model"
	"github.com/devsamuele/contact-service/business/data/contact/store"
	"github.com/devsamuele/contact-service/business/sys/database"
	"github.com/devsamuele/service-kit/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Organization struct {
	organizationStore store.Organization
	officeStore       store.Office
	session           database.Session
}

func NewOrganization(session database.Session, organizationStore store.Organization, officeStore store.Office) Organization {
	return Organization{
		organizationStore: organizationStore,
		officeStore:       officeStore,
		session:           session,
	}
}

// Field

func (s Organization) QueryField(ctx context.Context, tenantID string) ([]model.FieldWithSection, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, web.NewShutdownError("invalid tenant_id")
		}
	}

	fields, err := s.organizationStore.QueryField(ctx, tenantOID)
	if err != nil {
		return nil, err
	}

	return fields, err
}

func (s Organization) QueryFieldByID(ctx context.Context, tenantID string, fieldID string) (model.FieldWithSection, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.FieldWithSection{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	fieldOID, err := primitive.ObjectIDFromHex(fieldID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.FieldWithSection{}, model.Error{
				Message:      "ID is not in its proper form",
				Reason:       model.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	field, err := s.organizationStore.QueryFieldByID(ctx, tenantOID, fieldOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.FieldWithSection{}, model.Error{
				Message:      "field not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "field_id",
			}
		}
		return model.FieldWithSection{}, err
	}

	return field, err
}

func (s Organization) InsertField(ctx context.Context, tenantID string, nf model.NewFieldWithSection, now time.Time) (model.FieldWithSection, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.FieldWithSection{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	err = s.organizationStore.CheckFieldByLabel(ctx, tenantOID, *nf.Label)
	if err != nil {
		if !errors.Is(err, store.ErrNotFound) {
			return model.FieldWithSection{}, err
		}
	} else {
		return model.FieldWithSection{}, model.Error{
			Message:      "label already exist",
			Reason:       model.ErrReasonConflict,
			LocationType: "argument",
			Location:     "label",
		}
	}

	err = s.organizationStore.CheckSectionByID(ctx, tenantOID, *nf.SectionID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.FieldWithSection{}, model.Error{
				Message:      "section not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "argument",
				Location:     "section_id",
			}
		}
		return model.FieldWithSection{}, err
	}

	if err := nf.Validate(model.OrganizationFieldTypes()); err != nil {
		return model.FieldWithSection{}, err
	}

	field := model.FieldWithSection{
		Field: model.Field{
			ID:            primitive.NewObjectID(),
			TenantID:      tenantOID,
			Name:          "",
			Type:          *nf.Type,
			Label:         *nf.Label,
			Custom:        true,
			Visible:       true,
			Searchable:    true,
			EditableValue: true,
			Created:       now,
			Updated:       now,
		},
		SectionID: *nf.SectionID,
	}

	if nf.Searchable != nil {
		field.Searchable = *nf.Searchable
	}

	field.Choices = nf.BuildChoices()
	field.Options = nf.BuildOptions()

	// ---------------------------------------------------------------------------------

	sess, err := s.session.Start()
	if err != nil {
		return model.FieldWithSection{}, err
	}

	defer sess.EndSession(ctx)

	txCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {

		err = s.organizationStore.InsertField(ctx, field)
		if err != nil {
			return model.FieldWithSection{}, err
		}

		count, err := s.organizationStore.Count(ctx, tenantOID)
		if err != nil {
			return model.FieldWithSection{}, err
		}

		if count > 0 {
			var value interface{}
			if model.IsFieldWithArrayValue(field.Type) {
				value = bson.A{}
			}

			if err := s.organizationStore.AddFieldValue(ctx, tenantOID, field.ID, value, now); err != nil {
				return model.FieldWithSection{}, err
			}
		}
		return nil, nil
	}

	if _, err := sess.WithTransaction(ctx, txCallback); err != nil {
		return model.FieldWithSection{}, err
	}

	return field, nil
}

func (s Organization) UpdateField(ctx context.Context, tenantID string, fieldID string, uf model.UpdateFieldWithSection, now time.Time) error {

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

	oldField, err := s.organizationStore.QueryFieldByID(ctx, tenantOID, fieldOID)
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

	if err := s.organizationStore.CheckFieldByLabel(ctx, tenantOID, *uf.Label); err != nil {
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

	if uf.SectionID != nil {
		if err := s.organizationStore.CheckSectionByID(ctx, tenantOID, *uf.SectionID); err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return model.Error{
					Message:      "section not found",
					Reason:       model.ErrReasonNotFound,
					LocationType: "argument",
					Location:     "section_id",
				}
			}
		}
	}

	if err := uf.Validate(oldField); err != nil {
		return err
	}

	if uf.Label != nil {
		oldField.Label = *uf.Label
	}

	if uf.SectionID != nil {
		oldField.SectionID = *uf.SectionID
	}

	if uf.Searchable != nil {
		oldField.Searchable = *uf.Searchable
	}

	updatedChoices, removedChoiceIDs := uf.UpdateChoices(oldField.Field)
	oldField.Choices = updatedChoices

	oldField.Options = uf.UpdateOptions(oldField.Field)
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
					if err := s.organizationStore.RemoveChoices(sessCtx, tenantOID, fieldOID, removedChoiceIDs, now); err != nil {
						return nil, err
					}
				}

			case model.FieldTypeMultipleSelect:
				if len(removedChoiceIDs) > 0 {
					if err := s.organizationStore.PullChoices(sessCtx, tenantOID, fieldOID, removedChoiceIDs, now); err != nil {
						return nil, err
					}
				}
			}
		}

		if err := s.organizationStore.UpdateField(ctx, oldField); err != nil {
			return nil, err
		}

		return nil, nil
	}

	if _, err := sess.WithTransaction(ctx, txCallback); err != nil {
		return err
	}

	return nil
}

func (s Organization) DeleteFieldByID(ctx context.Context, tenantID, fieldID string, now time.Time) error {

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

	_, err = s.organizationStore.QueryFieldByID(ctx, tenantOID, fieldOID)
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
		err = s.organizationStore.DeleteFieldByID(sessCtx, tenantOID, fieldOID)
		if err != nil {
			return nil, err
		}

		err = s.organizationStore.UnsetFieldValues(sessCtx, tenantOID, []primitive.ObjectID{fieldOID}, now)
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

func (s Organization) DeleteFieldByIDs(ctx context.Context, tenantID string, fieldIDs model.ObjectIDs, now time.Time) error {

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
		err = s.organizationStore.DeleteFieldByIDs(sessCtx, tenantOID, fieldIDs.IDs)
		if err != nil {
			return nil, err
		}

		err = s.organizationStore.UnsetFieldValues(sessCtx, tenantOID, fieldIDs.IDs, now)
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

// Organization

func (s Organization) QueryByID(ctx context.Context, tenantID, organizationID string) (model.OrganizationResponse, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.OrganizationResponse{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	organizationOID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.OrganizationResponse{}, model.Error{
				Message:      "ID is not in its proper form",
				Reason:       model.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	organization, err := s.organizationStore.QueryByID(ctx, tenantOID, organizationOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.OrganizationResponse{}, model.Error{
				Message:      "organization not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return model.OrganizationResponse{}, err
	}

	offices, err := s.officeStore.QueryByOrganizationID(ctx, tenantOID, organizationOID)

	return model.OrganizationResponse{
		Organization: organization,
		Offices:      offices,
	}, nil
}

func (s Organization) Create(ctx context.Context, tenantID string, no model.NewOrganization, now time.Time) (model.Organization, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Organization{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	fields, err := s.organizationStore.QueryField(ctx, tenantOID)
	if err != nil {
		return model.Organization{}, err
	}

	if no.Name == nil || len(*no.Name) < 3 {
		return model.Organization{}, model.Error{
			Message:      "name is required",
			Reason:       model.ErrReasonRequired,
			LocationType: "argument",
			Location:     "name",
		}
	}

	fieldValues := make(model.FieldValues, 0)
	for _, field := range fields {
		fieldValues[field.ID.Hex()] = nil
		if model.IsFieldWithArrayValue(field.Type) {
			fieldValues[field.ID.Hex()] = bson.A{}
		}
	}

	newOrganization := model.Organization{
		ID:          primitive.NewObjectID(),
		TenantID:    tenantOID,
		Name:        *no.Name,
		Phone:       make([]string, 0),
		FieldValues: fieldValues,
		// Deleted:      false,
		// DateOfDelete: nil,
		Created: now,
		Updated: now,
	}

	err = s.organizationStore.Insert(ctx, newOrganization)
	if err != nil {
		return model.Organization{}, err
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

	return newOrganization, nil
}

func (s Organization) Update(ctx context.Context, tenantID, organizationID string, uo model.UpdateOrganization, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	organizationOID, err := primitive.ObjectIDFromHex(organizationID)
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

	org, err := s.organizationStore.QueryByID(ctx, tenantOID, organizationOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Error{
				Message:      "organization not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return err
	}

	if uo.Name != nil {
		if len(*uo.Name) < 3 {
			return model.Error{
				Message:      "name is required",
				Reason:       model.ErrReasonRequired,
				LocationType: "argument",
				Location:     "name",
			}
		}
		org.Name = *uo.Name
	}

	if uo.Phone != nil {
		filter := model.Filter{
			Sort:     model.Sort{},
			ORFields: make([]model.ANDFilterField, 0),
		}

		for _, phone := range uo.Phone {
			filter.ORFields = append(filter.ORFields, model.ANDFilterField{ANDFields: []model.FilterField{
				{Resource: "person",
					FieldPath: "phone",
					FieldType: model.SearchFieldTypeString,
					Value:     phone,
					MatchType: model.FilterMatchTypeEqual},
			}})
		}

		organization, err := s.organizationStore.Search(ctx, tenantOID, filter, 1, 1)
		if err != nil {
			return err
		}
		if len(organization) != 0 && organization[0].ID.Hex() != organizationID {
			return model.Error{
				Message:      "phone already exists",
				Reason:       model.ErrReasonConflict,
				LocationType: "argument",
				Location:     "phone",
			}
		}
	}

	if uo.Phone != nil {
		org.Phone = uo.Phone
	}

	fields, err := s.organizationStore.QueryField(ctx, tenantOID)
	if err != nil {
		return err
	}

	decodedFieldValues, err := uo.FieldValues.Validate(model.FieldsWithSection(fields).ToField())
	if err != nil {
		return err
	}

	for fieldID, fieldValue := range decodedFieldValues {
		org.FieldValues[fieldID] = fieldValue
	}

	org.Updated = now

	err = s.organizationStore.Update(ctx, org)
	if err != nil {
		return err
	}

	return nil
}

// func (s Organization) Trash(ctx context.Context, tenantID, organizationID string, now time.Time) error {

// 	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
// 	if err != nil {
// 		if errors.Is(err, primitive.ErrInvalidHex) {
// 			return web.NewShutdownError("invalid tenant_id")
// 		}
// 	}

// 	organizationOID, err := primitive.ObjectIDFromHex(organizationID)
// 	if err != nil {
// 		if errors.Is(err, primitive.ErrInvalidHex) {
// 			return model.Error{
// 				Message:      "ID is not in its proper form",
// 				Reason:       model.ErrReasonInvalidParameter,
// 				LocationType: "parameter",
// 				Location:     "url",
// 			}
// 		}
// 	}

// 	_, err = s.organizationStore.QueryByID(ctx, tenantOID, organizationOID)
// 	if err != nil {
// 		if errors.Is(err, store.ErrNotFound) {
// 			return model.Error{
// 				Message:      "organization not found",
// 				Reason:       model.ErrReasonNotFound,
// 				LocationType: "parameter",
// 				Location:     "url",
// 			}
// 		}
// 		return err
// 	}

// 	err = s.organizationStore.TrashByID(ctx, tenantOID, organizationOID, now)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s Organization) Restore(ctx context.Context, tenantID, organizationID string) error {

// 	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
// 	if err != nil {
// 		if errors.Is(err, primitive.ErrInvalidHex) {
// 			return web.NewShutdownError("invalid tenant_id")
// 		}
// 	}

// 	organizationOID, err := primitive.ObjectIDFromHex(organizationID)
// 	if err != nil {
// 		if errors.Is(err, primitive.ErrInvalidHex) {
// 			return model.Error{
// 				Message:      "ID is not in its proper form",
// 				Reason:       model.ErrReasonInvalidParameter,
// 				LocationType: "parameter",
// 				Location:     "url",
// 			}
// 		}
// 	}

// 	_, err = s.organizationStore.QueryByID(ctx, tenantOID, organizationOID)
// 	if err != nil {
// 		if errors.Is(err, store.ErrNotFound) {
// 			return model.Error{
// 				Message:      "organization not found",
// 				Reason:       model.ErrReasonNotFound,
// 				LocationType: "parameter",
// 				Location:     "url",
// 			}
// 		}
// 		return err
// 	}

// 	err = s.organizationStore.RestoreByID(ctx, tenantOID, organizationOID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (s Organization) Delete(ctx context.Context, tenantID, organizationID string) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	organizationOID, err := primitive.ObjectIDFromHex(organizationID)
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

	_, err = s.organizationStore.QueryByID(ctx, tenantOID, organizationOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Error{
				Message:      "organization not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return err
	}

	sess, err := s.session.Start()
	if err != nil {
		return err
	}

	defer sess.EndSession(ctx)

	txCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		err = s.organizationStore.DeleteByID(sessCtx, tenantOID, organizationOID)
		if err != nil {
			return nil, err
		}

		if err := s.officeStore.DeleteByOrganizationID(sessCtx, tenantOID, organizationOID); err != nil {
			return nil, err
		}
		return nil, nil
	}

	if _, err := sess.WithTransaction(ctx, txCallback); err != nil {
		return err
	}

	return nil
}

func (s Organization) Search(ctx context.Context, tenantID string, nf model.NewFilter) ([]model.Organization, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, web.NewShutdownError("invalid tenant_id")
		}
	}

	if nf.PageNumber <= 0 {
		nf.PageNumber = 1
	}

	if nf.ItemsPerPage <= 0 {
		nf.ItemsPerPage = 50
	}

	// TODO VALIDATION

	filter := model.Filter{
		Sort: nf.Sort,
	}

	resourceMap := make(map[string]string)

	filter.ORFields = make([]model.ANDFilterField, 0)
	for _, ANDFields := range nf.ORFields {
		filterFields := make([]model.FilterField, 0)
		for _, filterField := range ANDFields.ANDFields {

			ff := model.FilterField{
				Resource:  filterField.Resource,
				FieldPath: filterField.FieldPath,
				Value:     filterField.Value,
				MatchType: filterField.MatchType,
			}

			resourceMap[filterField.Resource] = filterField.Resource

			splitPath := strings.Split(filterField.FieldPath, ".")

			switch filterField.Resource {
			case "organization":
				if len(splitPath) == 1 {
					switch splitPath[0] {
					case "_id":
						ff.FieldType = model.SearchFieldTypeObjectID
					case "name", "phone":
						ff.FieldType = model.SearchFieldTypeString
					case "created", "updated":
						ff.FieldType = model.SearchFieldTypeDatetime
					default:
						return nil, model.Error{
							Message:      fmt.Sprintf("invalid field path %v", filterField.FieldPath),
							Reason:       model.ErrReasonInvalidArgument,
							LocationType: "argument",
							Location:     "field_path",
						}
					}
				} else if len(splitPath) > 1 {
					fieldID, err := primitive.ObjectIDFromHex(splitPath[len(splitPath)-1])
					if err != nil {
						if errors.Is(err, primitive.ErrInvalidHex) {
							return nil, model.Error{
								Message:      fmt.Sprintf("invalid field path %v", filterField.FieldPath),
								Reason:       model.ErrReasonInvalidArgument,
								LocationType: "argument",
								Location:     "field_path",
							}
						}
						return nil, err
					}

					f, err := s.organizationStore.QueryFieldByID(ctx, tenantOID, fieldID)
					if err != nil {
						if errors.Is(err, store.ErrNotFound) {
							return nil, model.Error{
								Message:      fmt.Sprintf("invalid field path %v", filterField.FieldPath),
								Reason:       model.ErrReasonNotFound,
								LocationType: "argument",
								Location:     "field_path",
							}
						}
						return nil, err
					}
					ff.FieldType = f.Type

					searchFieldType := model.ToSearchFieldType(ff.FieldType)
					if !model.ValidMatchType(searchFieldType, ff.MatchType) {
						if err != nil {
							if errors.Is(err, model.ErrInvalidMatchType) {
								return nil, model.Error{
									Message:      "invalid match type",
									Reason:       model.ErrReasonInvalidArgument,
									LocationType: "argument",
									Location:     "match_type",
								}
							}
						}
					}
					ff.FieldType = searchFieldType

				} else {
					return nil, model.Error{
						Message:      fmt.Sprintf("invalid field path %v", filterField.FieldPath),
						Reason:       model.ErrReasonInvalidArgument,
						LocationType: "argument",
						Location:     "field_path",
					}
				}

			case "office":
				if len(splitPath) > 1 {
					fieldID, err := primitive.ObjectIDFromHex(splitPath[len(splitPath)-1])
					if err != nil {
						if errors.Is(err, primitive.ErrInvalidHex) {
							return nil, model.Error{
								Message:      fmt.Sprintf("invalid field path %v", filterField.FieldPath),
								Reason:       model.ErrReasonInvalidArgument,
								LocationType: "argument",
								Location:     "field_path",
							}
						}
						return nil, err
					}

					f, err := s.officeStore.QueryFieldByID(ctx, tenantOID, fieldID)
					if err != nil {
						if errors.Is(err, store.ErrNotFound) {
							return nil, model.Error{
								Message:      fmt.Sprintf("invalid field path %v", filterField.FieldPath),
								Reason:       model.ErrReasonNotFound,
								LocationType: "argument",
								Location:     "field_path",
							}
						}
						return nil, err
					}
					ff.FieldType = f.Type

					searchFieldType := model.ToSearchFieldType(ff.FieldType)
					if !model.ValidMatchType(searchFieldType, ff.MatchType) {
						if err != nil {
							if errors.Is(err, model.ErrInvalidMatchType) {
								return nil, model.Error{
									Message:      "invalid match type",
									Reason:       model.ErrReasonInvalidArgument,
									LocationType: "argument",
									Location:     "match_type",
								}
							}
						}
					}
					ff.FieldType = searchFieldType

				} else {
					return nil, model.Error{
						Message:      fmt.Sprintf("invalid field path %v", filterField.FieldPath),
						Reason:       model.ErrReasonInvalidArgument,
						LocationType: "argument",
						Location:     "field_path",
					}
				}

			default:
				return nil, model.Error{
					Message:      fmt.Sprintf("resource %v not valid", filterField.Resource),
					Reason:       model.ErrReasonInvalidArgument,
					LocationType: "argument",
					Location:     "resource",
				}
			}

			filterFields = append(filterFields, ff)
		}
		filter.ORFields = append(filter.ORFields, model.ANDFilterField{ANDFields: filterFields})
	}

	organizations, err := s.organizationStore.Search(ctx, tenantOID, filter, nf.PageNumber, nf.ItemsPerPage)
	if err != nil {
		if errors.Is(err, model.ErrInvalidMatchType) {
			log.Println("invalid")
			return nil, model.Error{
				Message:      "invalid match type",
				Reason:       model.ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "match_type",
			}
		}
		return nil, err
	}

	return organizations, nil
}

// Section

func (s Organization) QuerySection(ctx context.Context, tenantID string) ([]model.Section, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, web.NewShutdownError("invalid tenant_id")
		}
	}

	sections, err := s.organizationStore.QuerySection(ctx, tenantOID)
	if err != nil {
		return nil, err
	}

	return sections, err
}

func (s Organization) QuerySectionByID(ctx context.Context, tenantID, sectionID string) (model.Section, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Section{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	sectionOID, err := primitive.ObjectIDFromHex(sectionID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Section{}, model.Error{
				Message:      "ID is not in its proper form",
				Reason:       model.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	section, err := s.organizationStore.QuerySectionByID(ctx, tenantOID, sectionOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Section{}, model.Error{
				Message:      "section not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return model.Section{}, err
	}

	return section, err
}

func (s Organization) CreateSection(ctx context.Context, tenantID string, ns model.NewSection, now time.Time) (model.Section, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Section{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	if err := ns.Validate(); err != nil {
		return model.Section{}, err
	}

	_, err = s.organizationStore.QuerySectionByTitle(ctx, tenantOID, *ns.Title)
	if err != nil {
		if !errors.Is(err, store.ErrNotFound) {
			return model.Section{}, err
		}
	} else {
		return model.Section{}, model.Error{
			Message:      "title already exist",
			Reason:       model.ErrReasonConflict,
			LocationType: "argument",
			Location:     "title",
		}
	}

	newSection := model.Section{
		ID:       primitive.NewObjectID(),
		TenantID: tenantOID,
		Name:     "",
		Title:    *ns.Title,
		Custom:   true,
		Created:  now,
		Updated:  now,
	}

	if ns.Description != nil {
		newSection.Description = ns.Description
	}

	err = s.organizationStore.InsertSection(ctx, newSection)
	if err != nil {
		return model.Section{}, err
	}

	return newSection, nil
}

func (s Organization) UpdateSection(ctx context.Context, tenantID, sectionID string, us model.UpdateSection, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	sectionOID, err := primitive.ObjectIDFromHex(sectionID)
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

	oldSection, err := s.organizationStore.QuerySectionByID(ctx, tenantOID, sectionOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Error{
				Message:      "section not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return err
	}

	if err := us.Validate(oldSection); err != nil {
		return err
	}

	if us.Title != nil {

		section, err := s.organizationStore.QuerySectionByTitle(ctx, tenantOID, *us.Title)
		if err != nil {
			if !errors.Is(err, store.ErrNotFound) {
				return err
			}
		} else if oldSection.ID != section.ID {
			return model.Error{
				Message:      "title already exist",
				Reason:       model.ErrReasonConflict,
				LocationType: "argument",
				Location:     "title",
			}
		}
		oldSection.Title = *us.Title
	}

	if us.Description != nil {
		oldSection.Description = us.Description
	}

	oldSection.Updated = now

	err = s.organizationStore.UpdateSection(ctx, oldSection)
	if err != nil {
		return err
	}
	return nil
}

func (s Organization) DeleteSectionByID(ctx context.Context, tenantID, sectionID string, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	sectionOID, err := primitive.ObjectIDFromHex(sectionID)
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

	_, err = s.organizationStore.QuerySectionByID(ctx, tenantOID, sectionOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Error{
				Message:      "section not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return err
	}

	sess, err := s.session.Start()
	if err != nil {
		return err
	}

	txCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		err = s.organizationStore.DeleteSectionByID(sessCtx, tenantOID, sectionOID)
		if err != nil {
			return nil, err
		}

		fields, err := s.organizationStore.QueryFieldBySectionID(sessCtx, tenantOID, sectionOID)
		if err != nil {
			return nil, err
		}

		if len(fields) > 0 {
			err = s.organizationStore.DeleteFieldBySectionID(sessCtx, tenantOID, sectionOID)
			if err != nil {
				return nil, err
			}

			fieldIDs := make([]primitive.ObjectID, 0)
			for _, f := range fields {
				fieldIDs = append(fieldIDs, f.ID)
			}

			err = s.organizationStore.UnsetFieldValues(sessCtx, tenantOID, fieldIDs, now)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	_, err = sess.WithTransaction(ctx, txCallback)
	if err != nil {
		return err
	}

	return nil
}
