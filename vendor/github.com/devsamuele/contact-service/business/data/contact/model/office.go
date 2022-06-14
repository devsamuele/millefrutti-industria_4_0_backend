package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Office ...
type Office struct {
	ID             primitive.ObjectID     `json:"id" bson:"_id"`
	TenantID       primitive.ObjectID     `json:"tenant_id" bson:"tenant_id"`
	OrganizationID primitive.ObjectID     `json:"organization_id" bson:"organization_id"`
	FieldValues    map[string]interface{} `json:"field_values" bson:"field_values"`
	Created        time.Time              `json:"created" bson:"created"`
	Updated        time.Time              `json:"updated" bson:"updated"`
}

type NewOffice struct {
	OrganizationID *primitive.ObjectID `json:"organization_id"`
	FieldValues    FieldValues         `json:"field_values"`
}

type UpdateOffice struct {
	FieldValues FieldValues `json:"field_values"`
}

func OfficeFieldTypes() []string {
	return []string{FieldTypeText, FieldTypeRichText, FieldTypeAddress, FieldTypeURL,
		FieldTypeBoolean, FieldTypeNumber,
		FieldTypeCurrency, FieldTypeDuration, FieldTypePercent, FieldTypeRating,
		FieldTypeDatetime, FieldTypeSingleSelect, FieldTypeMultipleSelect, FieldTypeUser}

}
