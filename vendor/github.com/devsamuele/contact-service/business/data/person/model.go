package person

import (
	"time"

	"github.com/devsamuele/elit/field"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Person ...
type Person struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	TenantID    primitive.ObjectID `json:"tenant_id" bson:"tenant_id"`
	FieldValues field.Values       `json:"field_values" bson:"field_values"`
	Created     time.Time          `json:"created" bson:"created"`
	Updated     time.Time          `json:"updated" bson:"updated"`
}

type Input struct {
	FieldValues field.Values `json:"field_values"`
}

type Update struct {
	FieldValues field.Values `json:"field_values"`
}

type PersonResponse struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	TenantID     primitive.ObjectID `json:"tenant_id" bson:"tenant_id"`
	FieldValues  field.Values       `json:"field_values" bson:"field_values"`
	Organization *Organization      `json:"organization" bson:"organization"`
	Office       *Office            `json:"office" bson:"office"`
	Created      time.Time          `json:"created" bson:"created"`
	Updated      time.Time          `json:"updated" bson:"updated"`
}

type Organization struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	TenantID    primitive.ObjectID `json:"tenant_id" bson:"tenant_id"`
	FieldValues field.Values       `json:"field_values" bson:"field_values"`
	Created     time.Time          `json:"created" bson:"created"`
	Updated     time.Time          `json:"updated" bson:"updated"`
}

type Office struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	TenantID    primitive.ObjectID `json:"tenant_id" bson:"tenant_id"`
	FieldValues field.Values       `json:"field_values" bson:"field_values"`
	Created     time.Time          `json:"created" bson:"created"`
	Updated     time.Time          `json:"updated" bson:"updated"`
}

var fieldTypes = []string{
	field.TypeText, field.TypeRichText, field.TypeAddress, field.TypeURL,
	field.TypeBoolean, field.TypeNumber, field.TypeCurrency, field.TypeDuration,
	field.TypePercent, field.TypeRating, field.TypeDatetimeRange, field.TypeSingleSelect,
	field.TypeMultipleSelect, field.TypeUser, field.TypeEmail, field.TypePhone, field.TypeDatetime,
}
