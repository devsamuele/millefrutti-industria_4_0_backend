package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Organization ...
type Organization struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id"`
	TenantID    primitive.ObjectID     `json:"tenant_id" bson:"tenant_id"`
	Name        string                 `json:"name" bson:"name"`
	Phone       []string               `json:"phone" bson:"phone"`
	FieldValues map[string]interface{} `json:"field_values" bson:"field_values"`
	// Deleted      bool                   `json:"deleted,omitempty" bson:"deleted"`
	// DateOfDelete *time.Time             `json:"date_of_delete,omitempty" bson:"date_of_delete"`
	Created time.Time `json:"created" bson:"created"`
	Updated time.Time `json:"updated" bson:"updated"`
}

type OrganizationResponse struct {
	Organization `bson:"inline"`
	Offices      []Office `json:"offices" bson:"offices"`
}

type NewOrganization struct {
	Name *string `json:"name"`
}

type UpdateOrganization struct {
	Name        *string     `json:"name"`
	Phone       []string    `json:"phone"`
	FieldValues FieldValues `json:"field_values"`
}

func OrganizationFieldTypes() []string {
	return []string{FieldTypeText, FieldTypeRichText, FieldTypeAddress, FieldTypeURL,
		FieldTypeBoolean, FieldTypeNumber,
		FieldTypeCurrency, FieldTypeDuration, FieldTypePercent, FieldTypeRating,
		FieldTypeDatetime, FieldTypeSingleSelect, FieldTypeMultipleSelect, FieldTypeUser}
}
