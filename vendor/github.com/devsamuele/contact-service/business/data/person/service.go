package person

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/devsamuele/contact-service/business/data/contact/store"
	"github.com/devsamuele/contact-service/business/sys/database"
	"github.com/devsamuele/elit/field"
	"github.com/devsamuele/elit/filter"
	"github.com/devsamuele/elit/resource"
	"github.com/devsamuele/elit/resperr"
	"github.com/devsamuele/elit/section"
	"github.com/devsamuele/elit/utility/slice"
	"github.com/devsamuele/service-kit/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	session                database.Session
	resourceStore          resource.Store
	resourceFieldStore     field.Store
	organizationFieldStore field.Store
	sectionStore           section.Store
}

func NewService(session database.Session, resourceStore resource.Store, resourceFieldStore field.Store, sectionStore section.Store, organizationFieldStore field.Store) Service {
	return Service{
		resourceStore:      resourceStore,
		resourceFieldStore: resourceFieldStore,
		sectionStore:       sectionStore,
		session:            session,
	}
}

// Field

func (s Service) QueryField(ctx context.Context, tenantID string) ([]field.Fielder, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, web.NewShutdownError("invalid tenant_id")
		}
	}

	if err != nil {
		return nil, err
	}

	return s.resourceFieldStore.Query(ctx, tenantOID)
}

func (s Service) QueryFieldByID(ctx context.Context, tenantID string, fieldID string) (field.Fielder, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, web.NewShutdownError("invalid tenant_id")
		}
	}

	fieldOID, err := primitive.ObjectIDFromHex(fieldID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, resperr.Error{
				Message:      "ID is not in its proper form",
				Reason:       resperr.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	f, err := s.resourceFieldStore.QueryByID(ctx, tenantOID, fieldOID)
	if err != nil {
		if errors.Is(err, field.ErrNotFound) {
			return nil, resperr.Error{
				Message:      "field not found",
				Reason:       resperr.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "field_id",
			}
		}
		return nil, err
	}

	return f, err
}

func (s Service) InsertField(ctx context.Context, tenantID string, nf field.Input, now time.Time) (field.Fielder, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, web.NewShutdownError("invalid tenant_id")
		}
	}

	fieldBuilder, err := field.NewBuilder(field.Config{
		WithSection:     true,
		ValidFieldTypes: fieldTypes,
	}, s.resourceFieldStore)
	if err != nil {
		return nil, err
	}

	field, err := fieldBuilder.NewFromInput(ctx, tenantOID, nf, now)
	if err != nil {
		return nil, err
	}

	// ---------------------------------------------------------------------------------

	sess, err := s.session.Start()
	if err != nil {
		return nil, err
	}

	defer sess.EndSession(ctx)

	txCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {

		err = s.resourceFieldStore.Insert(sessCtx, field)
		if err != nil {
			return nil, err
		}

		if err := s.resourceStore.AddFieldValue(sessCtx, tenantOID, field.GetBase().ID, field.GetDefaultValue(), now); err != nil {
			return nil, err
		}
		return nil, nil
	}

	if _, err := sess.WithTransaction(ctx, txCallback); err != nil {
		return nil, err
	}

	return field, nil
}

func (s Service) UpdateField(ctx context.Context, tenantID string, fieldID string, uf field.Update, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	fieldOID, err := primitive.ObjectIDFromHex(fieldID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return resperr.Error{
				Message:      "ID is not in its proper form",
				Reason:       resperr.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	fieldBuilder, err := field.NewBuilder(field.Config{
		WithSection:     true,
		ValidFieldTypes: fieldTypes,
	}, s.resourceFieldStore)
	if err != nil {
		return err
	}

	sess, err := s.session.Start()
	if err != nil {
		return err
	}

	defer sess.EndSession(ctx)

	txCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {

		updatedField, err := fieldBuilder.UpdateFromInput(sessCtx, tenantOID, fieldOID, uf, s.resourceStore, now)
		if err != nil {
			return nil, err
		}

		if err := s.resourceFieldStore.Update(sessCtx, tenantOID, fieldOID, updatedField); err != nil {
			return nil, err
		}

		return nil, nil
	}

	if _, err := sess.WithTransaction(ctx, txCallback); err != nil {
		return err
	}

	return nil
}

func (s Service) DeleteFieldByID(ctx context.Context, tenantID, fieldID string, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	fieldOID, err := primitive.ObjectIDFromHex(fieldID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return resperr.Error{
				Message:      "ID is not in its proper form",
				Reason:       resperr.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	_, err = s.resourceFieldStore.QueryByID(ctx, tenantOID, fieldOID)
	if err != nil {
		if errors.Is(err, field.ErrNotFound) {
			return resperr.Error{}
		}
		return err
	}

	sess, err := s.session.Start()
	if err != nil {
		return err
	}

	txCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {

		err = s.resourceFieldStore.DeleteByID(sessCtx, tenantOID, fieldOID)
		if err != nil {
			return nil, err
		}

		err = s.resourceStore.UnsetFieldValues(sessCtx, tenantOID, []primitive.ObjectID{fieldOID}, now)
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

func (s Service) DeleteFieldByIDs(ctx context.Context, tenantID string, fieldIDs []primitive.ObjectID, now time.Time) error {

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
		err = s.resourceFieldStore.DeleteByIDs(sessCtx, tenantOID, fieldIDs)
		if err != nil {
			return nil, err
		}

		err = s.resourceStore.UnsetFieldValues(sessCtx, tenantOID, fieldIDs, now)
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

// resource.Resource

func (s Service) QueryByID(ctx context.Context, tenantID, resourceID string) (resource.Resource, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return resource.Resource{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	resourceOID, err := primitive.ObjectIDFromHex(resourceID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return resource.Resource{}, resperr.Error{
				Message:      "ID is not in its proper form",
				Reason:       resperr.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	res, err := s.resourceStore.QueryByID(ctx, tenantOID, resourceOID)
	if err != nil {
		if errors.Is(err, resource.ErrNotFound) {
			return resource.Resource{}, resperr.Error{
				Message:      "resource not found",
				Reason:       resperr.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return resource.Resource{}, err
	}

	return res, nil
}

func (s Service) Create(ctx context.Context, tenantID string, np Input, now time.Time) (resource.Resource, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return resource.Resource{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	fields, err := s.resourceFieldStore.Query(ctx, tenantOID)
	if err != nil {
		return resource.Resource{}, err
	}

	// match, err := CheckEmail(*np.PrimaryEmail)
	// if err != nil {
	// 	return resource.Resource{}, err
	// }

	// if !match {
	// 	return resource.Resource{}, Error{
	// 		Message:      "email not valid",
	// 		Reason:       ErrReasonInvalidArgument,
	// 		LocationType: "argument",
	// 		Location:     "primary_email",
	// 	}
	// }

	// filter := Filter{
	// 	Sort:     Sort{},
	// 	ORFields: make([]ANDFilterField, 0),
	// }

	// filter.ORFields = append(filter.ORFields, ANDFilterField{ANDFields: []FilterField{
	// 	{Resource: "resource",
	// 		FieldPath: "email",
	// 		FieldType: SearchFieldTypeString,
	// 		Value:     np.PrimaryEmail,
	// 		MatchType: FilterMatchTypeEqual},
	// }})

	// resources, err := s.resourceStore.Search(ctx, tenantOID, filter, 1, 1)
	// if err != nil {
	// 	return resource.Resource{}, err
	// }
	// if len(resources) != 0 {
	// 	return resource.Resource{}, Error{
	// 		Message:      "email already exists",
	// 		Reason:       ErrReasonConflict,
	// 		LocationType: "argument",
	// 		Location:     "primary_email",
	// 	}
	// }

	fieldValues := make(field.Values, 0)
	for _, f := range fields {
		fieldValue, ok := np.FieldValues[f.GetBase().ID.Hex()]
		if !ok {
			if f.GetBase().Required {
				return resource.Resource{}, resperr.Error{
					Message:      fmt.Sprintf("field[%v] is required", f.GetBase().ID.Hex()),
					Reason:       resperr.ErrReasonRequired,
					LocationType: "argument",
					Location:     fmt.Sprintf("field[%v]", f.GetBase().ID.Hex()),
				}
			}
			fieldValues[f.GetBase().ID.Hex()] = f.GetDefaultValue()
		} else {
			decodedValue, err := f.DecodeAndValidateValue(ctx, fieldValue, s.resourceStore)
			if err != nil {
				return resource.Resource{}, err
			}

			if f.GetBase().Name == "name" {
				name := decodedValue.(*string)
				if len(*name) < 3 {
					return resource.Resource{}, resperr.Error{
						Message:      fmt.Sprintf("field[%v] is required", f.GetBase().ID.Hex()),
						Reason:       resperr.ErrReasonRequired,
						LocationType: "argument",
						Location:     fmt.Sprintf("field[%v]", f.GetBase().ID.Hex()),
					}
				}
			}

			fieldValues[f.GetBase().ID.Hex()] = decodedValue
		}
	}

	res := resource.Resource{
		ID:          primitive.NewObjectID(),
		TenantID:    tenantOID,
		FieldValues: fieldValues,
		Created:     now,
		Updated:     now,
	}

	err = s.resourceStore.Insert(ctx, res)
	if err != nil {
		return resource.Resource{}, err
	}

	return res, nil
}

func (s Service) Update(ctx context.Context, tenantID, resourceID string, up Update, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	resourceOID, err := primitive.ObjectIDFromHex(resourceID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return resperr.Error{
				Message:      "ID is not in its proper form",
				Reason:       resperr.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	resource, err := s.resourceStore.QueryByID(ctx, tenantOID, resourceOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return resperr.Error{
				Message:      "resource not found",
				Reason:       resperr.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return err
	}

	// emails := make([]string, 0)
	// if up.PrimaryEmail != nil {
	// 	match, err := CheckEmail(*up.PrimaryEmail)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if !match {
	// 		return Error{
	// 			Message:      "email not valid",
	// 			Reason:       ErrReasonInvalidArgument,
	// 			LocationType: "argument",
	// 			Location:     "primary_email",
	// 		}
	// 	}
	// 	emails = []string{*up.PrimaryEmail}
	// }

	// if up.OthersEmail != nil {
	// 	for _, email := range up.OthersEmail {
	// 		match, err := CheckEmail(email)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if !match {
	// 			return Error{
	// 				Message:      "email not valid",
	// 				Reason:       ErrReasonInvalidArgument,
	// 				LocationType: "argument",
	// 				Location:     "others_email",
	// 			}
	// 		}
	// 	}
	// 	emails = append(emails, up.OthersEmail...)
	// }

	// filter := Filter{
	// 	Sort:     Sort{},
	// 	ORFields: make([]ANDFilterField, 0),
	// }

	// for _, email := range emails {
	// 	filter.ORFields = append(filter.ORFields, ANDFilterField{ANDFields: []FilterField{
	// 		{Resource: "resource",
	// 			FieldPath: "email",
	// 			FieldType: SearchFieldTypeString,
	// 			Value:     email,
	// 			MatchType: FilterMatchTypeEqual},
	// 	}})
	// }

	// if up.Phone != nil {
	// 	for _, phone := range up.Phone {
	// 		filter.ORFields = append(filter.ORFields, ANDFilterField{ANDFields: []FilterField{
	// 			{Resource: "resource",
	// 				FieldPath: "phone",
	// 				FieldType: SearchFieldTypeString,
	// 				Value:     phone,
	// 				MatchType: FilterMatchTypeEqual},
	// 		}})
	// 	}
	// }

	// if len(emails) > 0 || up.Phone != nil {
	// 	resources, err := s.resourceStore.Search(ctx, tenantOID, Filter{
	// 		Sort: Sort{},
	// 		ORFields: []ANDFilterField{{[]FilterField{
	// 			{
	// 				Resource:  "resource",
	// 				FieldPath: "email",
	// 				FieldType: SearchFieldTypeStringList,
	// 				Value:     emails,
	// 				MatchType: FilterMatchTypeAnyOf,
	// 			},
	// 		}}},
	// 	}, 1, 1)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if len(resources) != 0 && resources[0].ID.Hex() != resourceID {
	// 		if resources[0].PrimaryEmail == *up.PrimaryEmail {
	// 			return Error{
	// 				Message:      "email already exists",
	// 				Reason:       ErrReasonConflict,
	// 				LocationType: "argument",
	// 				Location:     "primary_email",
	// 			}
	// 		}

	// 		for _, email := range resources[0].OthersEmail {
	// 			for _, newEmail := range up.OthersEmail {
	// 				if email == newEmail {
	// 					return Error{
	// 						Message:      "email already exists",
	// 						Reason:       ErrReasonConflict,
	// 						LocationType: "argument",
	// 						Location:     "others_email",
	// 					}
	// 				}
	// 			}
	// 		}

	// 		for _, phone := range resources[0].Phone {
	// 			for _, newPhone := range up.Phone {
	// 				if phone == newPhone {
	// 					return Error{
	// 						Message:      "phone already exists",
	// 						Reason:       ErrReasonConflict,
	// 						LocationType: "argument",
	// 						Location:     "phone",
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// if up.PrimaryEmail != nil {
	// 	resource.PrimaryEmail = *up.PrimaryEmail
	// }

	// if up.OthersEmail != nil {
	// 	resource.OthersEmail = up.OthersEmail
	// }

	// if up.Phone != nil {
	// 	resource.Phone = up.Phone
	// }

	// if up.OrganizationID != nil {
	// 	_, err := s.organizationStore.QueryByID(ctx, tenantOID, *up.OrganizationID)
	// 	if errors.Is(err, store.ErrNotFound) {
	// 		return Error{
	// 			Message:      "organization not found",
	// 			Reason:       ErrReasonNotFound,
	// 			LocationType: "argument",
	// 			Location:     "organization_id",
	// 		}
	// 	}
	// 	resource.OrganizationID = up.OrganizationID
	// }

	// if up.OfficeID != nil {
	// 	if resource.OrganizationID == nil {
	// 		return Error{
	// 			Message:      "organization_id is required",
	// 			Reason:       ErrReasonRequired,
	// 			LocationType: "argument",
	// 			Location:     "organization_id",
	// 		}
	// 	}
	// 	office, err := s.officeStore.QueryByID(ctx, tenantOID, *up.OfficeID)
	// 	if errors.Is(err, store.ErrNotFound) {
	// 		return Error{
	// 			Message:      "office not found",
	// 			Reason:       ErrReasonNotFound,
	// 			LocationType: "argument",
	// 			Location:     "office_id",
	// 		}
	// 	}

	// 	if *resource.OrganizationID != office.OrganizationID {
	// 		return Error{
	// 			Message:      fmt.Sprintf("office%v does not belong to organization[%v]", office.ID, *resource.OrganizationID),
	// 			Reason:       ErrReasonNotFound,
	// 			LocationType: "argument",
	// 			Location:     "office_id - organization_id",
	// 		}
	// 	}
	// 	resource.OfficeID = up.OfficeID
	// }

	fields, err := s.resourceFieldStore.Query(ctx, tenantOID)
	if err != nil {
		return err
	}

	for _, f := range fields {
		fieldValue, ok := up.FieldValues[f.GetBase().ID.Hex()]
		if ok {
			decodedValue, err := f.DecodeAndValidateValue(ctx, fieldValue, s.resourceStore)
			if err != nil {
				return err
			}

			if f.GetBase().Name == "name" {
				name := decodedValue.(*string)
				if len(*name) < 3 {
					return resperr.Error{
						Message:      fmt.Sprintf("field[%v] is required", f.GetBase().ID.Hex()),
						Reason:       resperr.ErrReasonRequired,
						LocationType: "argument",
						Location:     fmt.Sprintf("field[%v]", f.GetBase().ID.Hex()),
					}
				}
			}
			resource.FieldValues[f.GetBase().ID.Hex()] = decodedValue
		}
	}

	if err := s.resourceStore.Update(ctx, resource); err != nil {
		return err
	}

	return nil
}

func (s Service) Delete(ctx context.Context, tenantID, resourceID string) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	resourceOID, err := primitive.ObjectIDFromHex(resourceID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return resperr.Error{
				Message:      "ID is not in its proper form",
				Reason:       resperr.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	_, err = s.resourceStore.QueryByID(ctx, tenantOID, resourceOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return resperr.Error{
				Message:      "resource not found",
				Reason:       resperr.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return err
	}

	err = s.resourceStore.DeleteByID(ctx, tenantOID, resourceOID)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) Search(ctx context.Context, tenantID string, nf filter.Input) ([]resource.ResourceResponse, error) {

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

	resourceFields, err := s.resourceFieldStore.Query(ctx, tenantOID)
	if err != nil {
		return nil, err
	}

	// organizationFields, err := s.organizationFieldStore.Query(ctx, tenantOID)
	// if err != nil {
	// 	return nil, err
	// }

	builder := filter.NewBuilder(filter.Config{
		MainResource: filter.Resource{
			Name:   "person",
			Fields: resourceFields,
		},
	})

	ORGroup, err := builder.NewFromInput(nf)
	if err != nil {
		return nil, err
	}

	resources, err := s.resourceStore.Search(ctx, tenantOID, ORGroup, builder.FkResources(), nf.Sort, nf.PageNumber, nf.ItemsPerPage)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

// Section

func (s Service) QuerySection(ctx context.Context, tenantID string) ([]section.Section, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, web.NewShutdownError("invalid tenant_id")
		}
	}

	sections, err := s.sectionStore.Query(ctx, tenantOID)
	if err != nil {
		return nil, err
	}

	return sections, err
}

func (s Service) QuerySectionByID(ctx context.Context, tenantID, sectionID string) (section.Section, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return section.Section{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	sectionOID, err := primitive.ObjectIDFromHex(sectionID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return section.Section{}, resperr.Error{
				Message:      "ID is not in its proper form",
				Reason:       resperr.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	sec, err := s.sectionStore.QueryByID(ctx, tenantOID, sectionOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return section.Section{}, resperr.Error{
				Message:      "section not found",
				Reason:       resperr.ErrReasonNotFound,
				LocationType: "parameter",
				Location:     "url",
			}
		}
		return section.Section{}, err
	}

	return sec, err
}

func (s Service) CreateSection(ctx context.Context, tenantID string, ns section.NewSection, now time.Time) (section.Section, error) {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return section.Section{}, web.NewShutdownError("invalid tenant_id")
		}
	}

	if err := ns.Validate(); err != nil {
		return section.Section{}, err
	}

	// _, err = s.resourceStore.QuerySectionByTitle(ctx, tenantOID, *ns.Title)
	// if err != nil {
	// 	if !errors.Is(err, store.ErrNotFound) {
	// 		return Section{}, err
	// 	}
	// } else {
	// 	return Section{}, Error{
	// 		Message:      "title already exist",
	// 		Reason:       ErrReasonConflict,
	// 		LocationType: "argument",
	// 		Location:     "title",
	// 	}
	// }

	newSection := section.Section{
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

	err = s.sectionStore.Insert(ctx, newSection)
	if err != nil {
		return section.Section{}, err
	}

	return newSection, nil
}

func (s Service) UpdateSection(ctx context.Context, tenantID, sectionID string, us section.UpdateSection, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	sectionOID, err := primitive.ObjectIDFromHex(sectionID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return resperr.Error{
				Message:      "ID is not in its proper form",
				Reason:       resperr.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	oldSection, err := s.sectionStore.QueryByID(ctx, tenantOID, sectionOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return resperr.Error{
				Message:      "section not found",
				Reason:       resperr.ErrReasonNotFound,
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

		// section, err := s.resourceStore.QuerySectionByTitle(ctx, tenantOID, *us.Title)
		// if err != nil {
		// 	if !errors.Is(err, store.ErrNotFound) {
		// 		return err
		// 	}
		// } else if oldSection.ID != section.ID {
		// 	return Error{
		// 		Message:      "title already exist",
		// 		Reason:       ErrReasonConflict,
		// 		LocationType: "argument",
		// 		Location:     "title",
		// 	}
		// }
		oldSection.Title = *us.Title
	}

	if us.Description != nil {
		oldSection.Description = us.Description
	}

	oldSection.Updated = now

	err = s.sectionStore.Update(ctx, oldSection)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) DeleteSectionByID(ctx context.Context, tenantID, sectionID string, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	sectionOID, err := primitive.ObjectIDFromHex(sectionID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return resperr.Error{
				Message:      "ID is not in its proper form",
				Reason:       resperr.ErrReasonInvalidParameter,
				LocationType: "parameter",
				Location:     "url",
			}
		}
	}

	_, err = s.sectionStore.QueryByID(ctx, tenantOID, sectionOID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return resperr.Error{
				Message:      "section not found",
				Reason:       resperr.ErrReasonNotFound,
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
		err = s.sectionStore.DeleteByID(sessCtx, tenantOID, sectionOID)
		if err != nil {
			return nil, err
		}

		fields, err := s.resourceFieldStore.QueryBySectionID(sessCtx, tenantOID, sectionOID)
		if err != nil {
			return nil, err
		}

		if len(fields) > 0 {
			err = s.resourceFieldStore.DeleteBySectionID(sessCtx, tenantOID, sectionOID)
			if err != nil {
				return nil, err
			}

			fieldIDs := slice.Map(fields, func(field field.Fielder) primitive.ObjectID {
				return field.GetBase().ID
			})

			err = s.resourceStore.UnsetFieldValues(sessCtx, tenantOID, fieldIDs, now)
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
