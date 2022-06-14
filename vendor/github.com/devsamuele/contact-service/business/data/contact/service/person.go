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

type Person struct {
	personStore       store.Person
	organizationStore store.Organization
	officeStore       store.Office
	session           database.Session
}

func NewPerson(session database.Session, personStore store.Person, organizationStore store.Organization, officeStore store.Office) Person {
	return Person{
		personStore:       personStore,
		organizationStore: organizationStore,
		officeStore:       officeStore,
		session:           session,
	}
}

// Field

func (s Person) QueryField(ctx context.Context, tenantID string) ([]model.FieldWithSection, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, web.NewShutdownError("invalid tenant_id")
		}
	}

	fields, err := s.personStore.QueryField(ctx, tenantOID)
	if err != nil {
		return nil, err
	}

	return fields, err
}

func (s Person) QueryFieldByID(ctx context.Context, tenantID string, fieldID string) (model.FieldWithSection, error) {

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

	field, err := s.personStore.QueryFieldByID(ctx, tenantOID, fieldOID)
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

func (s Person) InsertField(ctx context.Context, tenantID string, nf model.NewFieldWithSection, now time.Time) (model.FieldWithSection, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.FieldWithSection{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	err = s.personStore.CheckFieldByLabel(ctx, tenantOID, *nf.Label)
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

	err = s.personStore.CheckSectionByID(ctx, tenantOID, *nf.SectionID)
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

	if err := nf.Validate(model.PersonFieldTypes()); err != nil {
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

		err = s.personStore.InsertField(sessCtx, field)
		if err != nil {
			return model.FieldWithSection{}, err
		}

		count, err := s.personStore.Count(sessCtx, tenantOID)
		if err != nil {
			return model.FieldWithSection{}, err
		}

		if count > 0 {
			var value interface{}
			if model.IsFieldWithArrayValue(field.Type) {
				value = bson.A{}
			}

			if err := s.personStore.AddFieldValue(sessCtx, tenantOID, field.ID, value, now); err != nil {
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

func (s Person) UpdateField(ctx context.Context, tenantID string, fieldID string, uf model.UpdateFieldWithSection, now time.Time) error {

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

	oldField, err := s.personStore.QueryFieldByID(ctx, tenantOID, fieldOID)
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

	if err := s.personStore.CheckFieldByLabel(ctx, tenantOID, *uf.Label); err != nil {
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
		if err := s.personStore.CheckSectionByID(ctx, tenantOID, *uf.SectionID); err != nil {
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
					if err := s.personStore.RemoveChoices(sessCtx, tenantOID, fieldOID, removedChoiceIDs, now); err != nil {
						return nil, err
					}
				}

			case model.FieldTypeMultipleSelect:
				if len(removedChoiceIDs) > 0 {
					if err := s.personStore.PullChoices(sessCtx, tenantOID, fieldOID, removedChoiceIDs, now); err != nil {
						return nil, err
					}
				}
			}
		}

		if err := s.personStore.UpdateField(ctx, oldField); err != nil {
			return nil, err
		}

		return nil, nil
	}

	if _, err := sess.WithTransaction(ctx, txCallback); err != nil {
		return err
	}

	return nil
}

func (s Person) DeleteFieldByID(ctx context.Context, tenantID, fieldID string, now time.Time) error {

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

	_, err = s.personStore.QueryFieldByID(ctx, tenantOID, fieldOID)
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
		err = s.personStore.DeleteFieldByID(sessCtx, tenantOID, fieldOID)
		if err != nil {
			return nil, err
		}

		err = s.personStore.UnsetFieldValues(sessCtx, tenantOID, []primitive.ObjectID{fieldOID}, now)
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

func (s Person) DeleteFieldByIDs(ctx context.Context, tenantID string, fieldIDs model.ObjectIDs, now time.Time) error {

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
		err = s.personStore.DeleteFieldByIDs(sessCtx, tenantOID, fieldIDs.IDs)
		if err != nil {
			return nil, err
		}

		err = s.personStore.UnsetFieldValues(sessCtx, tenantOID, fieldIDs.IDs, now)
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

// Person

func (s Person) QueryByID(ctx context.Context, tenantID, personID string) (model.PersonResponse, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.PersonResponse{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	personOID, err := primitive.ObjectIDFromHex(personID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.PersonResponse{}, model.Error{
				Message:      "ID is not in its proper form",
				Reason:       model.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	person, err := s.personStore.QueryResponseByID(ctx, tenantOID, personOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.PersonResponse{}, model.Error{
				Message:      "person not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return model.PersonResponse{}, err
	}

	return person, nil
}

func (s Person) Create(ctx context.Context, tenantID string, np model.NewPerson, now time.Time) (model.Person, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Person{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	fields, err := s.personStore.QueryField(ctx, tenantOID)
	if err != nil {
		return model.Person{}, err
	}

	if np.Name == nil || len(*np.Name) < 3 {
		return model.Person{}, model.Error{
			Message:      "name is required",
			Reason:       model.ErrReasonRequired,
			LocationType: "argument",
			Location:     "name",
		}
	}

	if np.PrimaryEmail == nil {
		return model.Person{}, model.Error{
			Message:      "primary_email is required",
			Reason:       model.ErrReasonRequired,
			LocationType: "argument",
			Location:     "primary_id",
		}
	}

	match, err := model.CheckEmail(*np.PrimaryEmail)
	if err != nil {
		return model.Person{}, err
	}

	if !match {
		return model.Person{}, model.Error{
			Message:      "email not valid",
			Reason:       model.ErrReasonInvalidArgument,
			LocationType: "argument",
			Location:     "primary_email",
		}
	}

	filter := model.Filter{
		Sort:     model.Sort{},
		ORFields: make([]model.ANDFilterField, 0),
	}

	filter.ORFields = append(filter.ORFields, model.ANDFilterField{ANDFields: []model.FilterField{
		{Resource: "person",
			FieldPath: "email",
			FieldType: model.SearchFieldTypeString,
			Value:     np.PrimaryEmail,
			MatchType: model.FilterMatchTypeEqual},
	}})

	persons, err := s.personStore.Search(ctx, tenantOID, filter, 1, 1)
	if err != nil {
		return model.Person{}, err
	}
	if len(persons) != 0 {
		return model.Person{}, model.Error{
			Message:      "email already exists",
			Reason:       model.ErrReasonConflict,
			LocationType: "argument",
			Location:     "primary_email",
		}
	}

	fieldValues := make(model.FieldValues, 0)
	for _, field := range fields {
		fieldValues[field.ID.Hex()] = nil
		if model.IsFieldWithArrayValue(field.Type) {
			fieldValues[field.ID.Hex()] = bson.A{}
		}
	}

	newPerson := model.Person{
		ID:             primitive.NewObjectID(),
		TenantID:       tenantOID,
		Name:           *np.Name,
		PrimaryEmail:   *np.PrimaryEmail,
		FieldValues:    fieldValues,
		OfficeID:       nil,
		OrganizationID: nil,
		OthersEmail:    make([]string, 0),
		Phone:          make([]string, 0),
		// Deleted:        false,
		// DateOfDelete:   nil,
		Created: now,
		Updated: now,
	}

	err = s.personStore.Insert(ctx, newPerson)
	if err != nil {
		return model.Person{}, err
	}

	return newPerson, nil
}

func (s Person) Update(ctx context.Context, tenantID, personID string, up model.UpdatePerson, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	personOID, err := primitive.ObjectIDFromHex(personID)
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

	person, err := s.personStore.QueryByID(ctx, tenantOID, personOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Error{
				Message:      "person not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return err
	}

	if up.Name != nil {
		if len(*up.Name) < 3 {
			return model.Error{
				Message:      "name is required",
				Reason:       model.ErrReasonRequired,
				LocationType: "argument",
				Location:     "name",
			}
		}
		person.Name = *up.Name
	}

	emails := make([]string, 0)
	if up.PrimaryEmail != nil {
		match, err := model.CheckEmail(*up.PrimaryEmail)
		if err != nil {
			return err
		}
		if !match {
			return model.Error{
				Message:      "email not valid",
				Reason:       model.ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "primary_email",
			}
		}
		emails = []string{*up.PrimaryEmail}
	}

	if up.OthersEmail != nil {
		for _, email := range up.OthersEmail {
			match, err := model.CheckEmail(email)
			if err != nil {
				return err
			}
			if !match {
				return model.Error{
					Message:      "email not valid",
					Reason:       model.ErrReasonInvalidArgument,
					LocationType: "argument",
					Location:     "others_email",
				}
			}
		}
		emails = append(emails, up.OthersEmail...)
	}

	filter := model.Filter{
		Sort:     model.Sort{},
		ORFields: make([]model.ANDFilterField, 0),
	}

	for _, email := range emails {
		filter.ORFields = append(filter.ORFields, model.ANDFilterField{ANDFields: []model.FilterField{
			{Resource: "person",
				FieldPath: "email",
				FieldType: model.SearchFieldTypeString,
				Value:     email,
				MatchType: model.FilterMatchTypeEqual},
		}})
	}

	if up.Phone != nil {
		for _, phone := range up.Phone {
			filter.ORFields = append(filter.ORFields, model.ANDFilterField{ANDFields: []model.FilterField{
				{Resource: "person",
					FieldPath: "phone",
					FieldType: model.SearchFieldTypeString,
					Value:     phone,
					MatchType: model.FilterMatchTypeEqual},
			}})
		}
	}

	if len(emails) > 0 || up.Phone != nil {
		persons, err := s.personStore.Search(ctx, tenantOID, model.Filter{
			Sort: model.Sort{},
			ORFields: []model.ANDFilterField{{[]model.FilterField{
				{
					Resource:  "person",
					FieldPath: "email",
					FieldType: model.SearchFieldTypeStringList,
					Value:     emails,
					MatchType: model.FilterMatchTypeAnyOf,
				},
			}}},
		}, 1, 1)
		if err != nil {
			return err
		}
		if len(persons) != 0 && persons[0].ID.Hex() != personID {
			if persons[0].PrimaryEmail == *up.PrimaryEmail {
				return model.Error{
					Message:      "email already exists",
					Reason:       model.ErrReasonConflict,
					LocationType: "argument",
					Location:     "primary_email",
				}
			}

			for _, email := range persons[0].OthersEmail {
				for _, newEmail := range up.OthersEmail {
					if email == newEmail {
						return model.Error{
							Message:      "email already exists",
							Reason:       model.ErrReasonConflict,
							LocationType: "argument",
							Location:     "others_email",
						}
					}
				}
			}

			for _, phone := range persons[0].Phone {
				for _, newPhone := range up.Phone {
					if phone == newPhone {
						return model.Error{
							Message:      "phone already exists",
							Reason:       model.ErrReasonConflict,
							LocationType: "argument",
							Location:     "phone",
						}
					}
				}
			}
		}
	}

	if up.PrimaryEmail != nil {
		person.PrimaryEmail = *up.PrimaryEmail
	}

	if up.OthersEmail != nil {
		person.OthersEmail = up.OthersEmail
	}

	if up.Phone != nil {
		person.Phone = up.Phone
	}

	if up.OrganizationID != nil {
		_, err := s.organizationStore.QueryByID(ctx, tenantOID, *up.OrganizationID)
		if errors.Is(err, store.ErrNotFound) {
			return model.Error{
				Message:      "organization not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "argument",
				Location:     "organization_id",
			}
		}
		person.OrganizationID = up.OrganizationID
	}

	if up.OfficeID != nil {
		if person.OrganizationID == nil {
			return model.Error{
				Message:      "organization_id is required",
				Reason:       model.ErrReasonRequired,
				LocationType: "argument",
				Location:     "organization_id",
			}
		}
		office, err := s.officeStore.QueryByID(ctx, tenantOID, *up.OfficeID)
		if errors.Is(err, store.ErrNotFound) {
			return model.Error{
				Message:      "office not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "argument",
				Location:     "office_id",
			}
		}

		if *person.OrganizationID != office.OrganizationID {
			return model.Error{
				Message:      fmt.Sprintf("office%v does not belong to organization[%v]", office.ID, *person.OrganizationID),
				Reason:       model.ErrReasonNotFound,
				LocationType: "argument",
				Location:     "office_id - organization_id",
			}
		}
		person.OfficeID = up.OfficeID
	}

	fields, err := s.personStore.QueryField(ctx, tenantOID)
	if err != nil {
		return err
	}

	decodedFieldValues, err := up.FieldValues.Validate(model.FieldsWithSection(fields).ToField())
	if err != nil {
		return err
	}

	for fieldID, fieldValue := range decodedFieldValues {
		person.FieldValues[fieldID] = fieldValue
	}

	person.Updated = now
	err = s.personStore.Update(ctx, person)
	if err != nil {
		return err
	}

	return nil
}

// func (s Person) Trash(ctx context.Context, tenantID, personID string, now time.Time) error {

// 	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
// 	if err != nil {
// 		if errors.Is(err, primitive.ErrInvalidHex) {
// 			return web.NewShutdownError("invalid tenant_id")
// 		}
// 	}

// 	personOID, err := primitive.ObjectIDFromHex(personID)
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

// 	_, err = s.personStore.QueryByID(ctx, tenantOID, personOID)
// 	if err != nil {
// 		if errors.Is(err, store.ErrNotFound) {
// 			return model.Error{
// 				Message:      "person not found",
// 				Reason:       model.ErrReasonNotFound,
// 				LocationType: "parameter",
// 				Location:     "url",
// 			}
// 		}
// 		return err
// 	}

// 	err = s.personStore.TrashByID(ctx, tenantOID, personOID, now)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s Person) Restore(ctx context.Context, tenantID, personID string) error {

// 	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
// 	if err != nil {
// 		if errors.Is(err, primitive.ErrInvalidHex) {
// 			return web.NewShutdownError("invalid tenant_id")
// 		}
// 	}

// 	personOID, err := primitive.ObjectIDFromHex(personID)
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

// 	_, err = s.personStore.QueryByID(ctx, tenantOID, personOID)
// 	if err != nil {
// 		if errors.Is(err, store.ErrNotFound) {
// 			return model.Error{
// 				Message:      "person not found",
// 				Reason:       model.ErrReasonNotFound,
// 				LocationType: "parameter",
// 				Location:     "url",
// 			}
// 		}
// 		return err
// 	}

// 	err = s.personStore.RestoreByID(ctx, tenantOID, personOID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (s Person) Delete(ctx context.Context, tenantID, personID string) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	personOID, err := primitive.ObjectIDFromHex(personID)
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

	_, err = s.personStore.QueryByID(ctx, tenantOID, personOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.Error{
				Message:      "person not found",
				Reason:       model.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return err
	}

	err = s.personStore.DeleteByID(ctx, tenantOID, personOID)
	if err != nil {
		return err
	}

	return nil
}

func (s Person) Search(ctx context.Context, tenantID string, nf model.NewFilter) ([]model.PersonResponse, error) {

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
			case "person":
				if len(splitPath) == 1 {
					switch splitPath[0] {
					case "_id":
						ff.FieldType = model.SearchFieldTypeObjectID
					case "created", "updated":
						ff.FieldType = model.SearchFieldTypeDatetime
					case "email", "phone", "name":
						ff.FieldType = model.SearchFieldTypeString
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

					f, err := s.personStore.QueryFieldByID(ctx, tenantOID, fieldID)
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
			case "organization":
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

					f, err := s.personStore.QueryFieldByID(ctx, tenantOID, fieldID)
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

					f, err := s.personStore.QueryFieldByID(ctx, tenantOID, fieldID)
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

	persons, err := s.personStore.Search(ctx, tenantOID, filter, nf.PageNumber, nf.ItemsPerPage)
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

	return persons, nil
}

// Section

func (s Person) QuerySection(ctx context.Context, tenantID string) ([]model.Section, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, web.NewShutdownError("invalid tenant_id")
		}
	}

	sections, err := s.personStore.QuerySection(ctx, tenantOID)
	if err != nil {
		return nil, err
	}

	return sections, err
}

func (s Person) QuerySectionByID(ctx context.Context, tenantID, sectionID string) (model.Section, error) {

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

	section, err := s.personStore.QuerySectionByID(ctx, tenantOID, sectionOID)
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

func (s Person) CreateSection(ctx context.Context, tenantID string, ns model.NewSection, now time.Time) (model.Section, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return model.Section{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	if err := ns.Validate(); err != nil {
		return model.Section{}, err
	}

	_, err = s.personStore.QuerySectionByTitle(ctx, tenantOID, *ns.Title)
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

	err = s.personStore.InsertSection(ctx, newSection)
	if err != nil {
		return model.Section{}, err
	}

	return newSection, nil
}

func (s Person) UpdateSection(ctx context.Context, tenantID, sectionID string, us model.UpdateSection, now time.Time) error {

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

	oldSection, err := s.personStore.QuerySectionByID(ctx, tenantOID, sectionOID)
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

		section, err := s.personStore.QuerySectionByTitle(ctx, tenantOID, *us.Title)
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

	err = s.personStore.UpdateSection(ctx, oldSection)
	if err != nil {
		return err
	}
	return nil
}

func (s Person) DeleteSectionByID(ctx context.Context, tenantID, sectionID string, now time.Time) error {

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

	_, err = s.personStore.QuerySectionByID(ctx, tenantOID, sectionOID)
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
		err = s.personStore.DeleteSectionByID(sessCtx, tenantOID, sectionOID)
		if err != nil {
			return nil, err
		}

		fields, err := s.personStore.QueryFieldBySectionID(sessCtx, tenantOID, sectionOID)
		if err != nil {
			return nil, err
		}

		if len(fields) > 0 {
			err = s.personStore.DeleteFieldBySectionID(sessCtx, tenantOID, sectionOID)
			if err != nil {
				return nil, err
			}

			fieldIDs := make([]primitive.ObjectID, 0)
			for _, f := range fields {
				fieldIDs = append(fieldIDs, f.ID)
			}

			err = s.personStore.UnsetFieldValues(sessCtx, tenantOID, fieldIDs, now)
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
