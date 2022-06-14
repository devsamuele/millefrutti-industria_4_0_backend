package service

import (
	"context"
	"errors"
	"time"

	"github.com/devsamuele/contact-service/business/sys/database"
	"github.com/devsamuele/elit/field"
	"github.com/devsamuele/elit/resource"
	"github.com/devsamuele/elit/section"
	"github.com/devsamuele/elit/translation"
	"github.com/devsamuele/service-kit/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	session      database.Session
	personStore  resource.Store
	fieldStore   field.Store
	sectionStore section.Store
}

func NewAdmin(session database.Session, personStore resource.Store, fieldStore field.Store, sectionStore section.Store) Admin {
	return Admin{
		session:      session,
		personStore:  personStore,
		fieldStore:   fieldStore,
		sectionStore: sectionStore,
	}
}

func (s Admin) InitTenant(ctx context.Context, tenantID string, now time.Time) error {

	tenantOID, err := primitive.ObjectIDFromHex(tenantID)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return web.NewShutdownError("invalid tenant_id")
		}
	}

	// organizationSection := model.Section{
	// 	ID:       primitive.NewObjectID(),
	// 	TenantID: tenantOID,
	// 	Name:     "info",
	// 	Title: model.Translation{
	// 		It: "info",
	// 		En: "info",
	// 		De: "info",
	// 		Fr: "info",
	// 		Es: "info",
	// 	},
	// 	Description: nil,
	// 	Custom:      false,
	// 	Created:     now,
	// 	Updated:     now,
	// }

	// if err := s.organizationStore.InsertSection(ctx, organizationSection); err != nil {
	// 	return err
	// }

	// organizationNameField := model.FieldWithSection{
	// 	Field: model.Field{
	// 		ID:       primitive.NewObjectID(),
	// 		TenantID: tenantOID,
	// 		Name:     "name",
	// 		Label: model.Translation{
	// 			It: "name",
	// 			De: "name",
	// 			En: "name",
	// 			Fr: "name",
	// 			Es: "name",
	// 		},
	// 		Type:          "text",
	// 		Custom:        false,
	// 		Visible:       true,
	// 		EditableValue: true,
	// 		Searchable:    true,
	// 		Created:       now,
	// 		Updated:       now,
	// 	},
	// 	SectionID: organizationSection.ID,
	// }

	// if err := s.organizationStore.InsertField(ctx, organizationNameField); err != nil {
	// 	return err
	// }

	// organizationStatusField := model.FieldWithSection{
	// 	Field: model.Field{
	// 		ID:       primitive.NewObjectID(),
	// 		TenantID: tenantOID,
	// 		Name:     "lead_status",
	// 		Label: model.Translation{
	// 			It: "lead status",
	// 			De: "lead status",
	// 			En: "lead status",
	// 			Fr: "lead status",
	// 			Es: "lead status",
	// 		},
	// 		Type: "single_select",
	// 		Choices: []model.Choice{
	// 			{
	// 				ID: primitive.NewObjectID(),
	// 				Value: model.Translation{
	// 					It: "potenziale",
	// 					En: "potenziale",
	// 					De: "potenziale",
	// 					Fr: "potenziale",
	// 					Es: "potenziale",
	// 				},
	// 				Custom: false,
	// 			},
	// 			{
	// 				ID: primitive.NewObjectID(),
	// 				Value: model.Translation{
	// 					It: "acquisito",
	// 					En: "acquisito",
	// 					De: "acquisito",
	// 					Fr: "acquisito",
	// 					Es: "acquisito",
	// 				},
	// 				Custom: false,
	// 			},
	// 		},
	// 		Custom:        false,
	// 		Options:       model.SelectTypeOptions{InsertChoices: true},
	// 		Visible:       true,
	// 		EditableValue: true,
	// 		Searchable:    true,
	// 		Created:       now,
	// 		Updated:       now,
	// 	},
	// 	SectionID: organizationSection.ID,
	// }

	// if err := s.organizationStore.InsertField(ctx, organizationStatusField); err != nil {
	// 	return err
	// }

	// Person

	personInfoSection := section.Section{
		ID:       primitive.NewObjectID(),
		TenantID: tenantOID,
		Name:     "info",
		Title: translation.Translation{
			It: "info",
			En: "info",
			De: "info",
			Fr: "info",
			Es: "info",
		},
		Description: nil,
		Custom:      false,
		Created:     now,
		Updated:     now,
	}

	if err := s.sectionStore.Insert(ctx, personInfoSection); err != nil {
		return err
	}

	personNameField := field.NewText(tenantOID, "name", translation.Translation{
		It: "name",
		En: "name",
		De: "name",
		Fr: "name",
		Es: "name",
	}, false, true, true, true, true, &personInfoSection.ID, nil, now)

	if err := s.fieldStore.Insert(ctx, personNameField); err != nil {
		return err
	}

	organizationFKField := field.NewForeignKey(tenantOID, "organization", translation.Translation{
		It: "organization_id",
		En: "organization_id",
		De: "organization_id",
		Fr: "organization_id",
		Es: "organization_id",
	}, false, false, true, false, false, nil, nil, now)

	if err := s.fieldStore.Insert(ctx, organizationFKField); err != nil {
		return err
	}

	officeFKField := field.NewForeignKey(tenantOID, "office", translation.Translation{
		It: "office_id",
		En: "office_id",
		De: "office_id",
		Fr: "office_id",
		Es: "office_id",
	}, false, false, true, false, false, nil, nil, now)

	if err := s.fieldStore.Insert(ctx, officeFKField); err != nil {
		return err
	}

	// // Office
	// officeNameField := model.Field{
	// 	ID:       primitive.NewObjectID(),
	// 	TenantID: tenantOID,
	// 	Name:     "name",
	// 	Label: model.Translation{
	// 		It: "name",
	// 		De: "name",
	// 		En: "name",
	// 		Fr: "name",
	// 		Es: "name",
	// 	},
	// 	Type:          "text",
	// 	Custom:        false,
	// 	Visible:       true,
	// 	EditableValue: true,
	// 	Searchable:    true,
	// 	Created:       now,
	// 	Updated:       now,
	// }

	// if err := s.officeStore.InsertField(ctx, officeNameField); err != nil {
	// 	return err
	// }

	return nil

}
