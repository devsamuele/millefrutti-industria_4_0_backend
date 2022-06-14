package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Person ...
type Person struct {
	ID             primitive.ObjectID     `json:"id" bson:"_id"`
	TenantID       primitive.ObjectID     `json:"tenant_id" bson:"tenant_id"`
	OrganizationID *primitive.ObjectID    `json:"organization_id" bson:"organization_id"`
	OfficeID       *primitive.ObjectID    `json:"office_id" bson:"office_id"`
	Name           string                 `json:"name" bson:"name"`
	PrimaryEmail   string                 `json:"primary_email" bson:"primary_email"`
	OthersEmail    []string               `json:"others_email" bson:"others_email"`
	Phone          []string               `json:"phone" bson:"phone"`
	FieldValues    map[string]interface{} `json:"field_values" bson:"field_values"`
	// Deleted        bool                   `json:"deleted,omitempty" bson:"deleted"`
	// DateOfDelete   *time.Time             `json:"date_of_delete,omitempty" bson:"date_of_delete"`
	Created time.Time `json:"created" bson:"created"`
	Updated time.Time `json:"updated" bson:"updated"`
}

type PersonResponse struct {
	ID           primitive.ObjectID     `json:"id" bson:"_id"`
	TenantID     primitive.ObjectID     `json:"tenant_id" bson:"tenant_id"`
	Organization *Organization          `json:"organization" bson:"organization"`
	Office       *Office                `json:"office" bson:"office"`
	Name         string                 `json:"name" bson:"name"`
	PrimaryEmail string                 `json:"primary_email" bson:"primary_email"`
	OthersEmail  []string               `json:"others_email" bson:"others_email"`
	Phone        []string               `json:"phone" bson:"phone"`
	FieldValues  map[string]interface{} `json:"field_values" bson:"field_values"`
	Deleted      bool                   `json:"deleted,omitempty" bson:"deleted"`
	DateOfDelete *time.Time             `json:"date_of_delete,omitempty" bson:"date_of_delete"`
	Created      time.Time              `json:"created" bson:"created"`
	Updated      time.Time              `json:"updated" bson:"updated"`
}

type NewPerson struct {
	Name         *string `json:"name"`
	PrimaryEmail *string `json:"primary_email"`
}

type UpdatePerson struct {
	Name           *string             `json:"name"`
	PrimaryEmail   *string             `json:"primary_email"`
	OthersEmail    []string            `json:"others_email"`
	FieldValues    FieldValues         `json:"field_values"`
	Phone          []string            `json:"phone"`
	OrganizationID *primitive.ObjectID `json:"organization_id"`
	OfficeID       *primitive.ObjectID `json:"office_id"`
}

func PersonFieldTypes() []string {
	return []string{FieldTypeText, FieldTypeRichText, FieldTypeAddress, FieldTypeURL,
		FieldTypeBoolean, FieldTypeNumber,
		FieldTypeCurrency, FieldTypeDuration, FieldTypePercent, FieldTypeRating,
		FieldTypeDatetime, FieldTypeSingleSelect, FieldTypeMultipleSelect, FieldTypeUser}
}
