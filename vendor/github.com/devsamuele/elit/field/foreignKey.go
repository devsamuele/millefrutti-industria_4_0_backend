package field

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/devsamuele/elit/resperr"
	"github.com/devsamuele/elit/translation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type foreignKeyInput struct {
// 	baseInput
// 	Resource *string `json:"resource"`
// }

// type foreignKeyUpdate baseUpdate

type foreignKeyResponse struct {
	base
	Resource string `json:"resource"`
}

type foreignKey struct {
	Base     base   `bson:",inline"`
	Resource string `bson:"resource"`
}

func (fk *foreignKey) UnmarshalBSON(data []byte) error {

	type foreignKeyAlias foreignKey
	var fka foreignKeyAlias

	if err := bson.Unmarshal(data, &fka); err != nil {
		return err
	}

	*fk = foreignKey(fka)

	return nil
}

func (fk *foreignKey) MarshalBSON() ([]byte, error) {

	type foreignKeyAlias foreignKey
	fka := foreignKeyAlias(*fk)
	return bson.Marshal(&fka)
}

func (fk *foreignKey) MarshalJSON() ([]byte, error) {

	type foreignKeyAlias foreignKey
	fka := foreignKeyAlias(*fk)
	fkr := foreignKeyResponse{
		base:     fka.Base,
		Resource: fka.Resource,
	}
	return json.Marshal(&fkr)
}

func NewForeignKey(tenantID primitive.ObjectID, resource string, label translation.Translation, custom, visible, editableValue, searchable, required bool, sectionID, categoryID *primitive.ObjectID, now time.Time) Fielder {
	return foreignKey{
		Base: base{
			ID:            primitive.NewObjectID(),
			TenantID:      tenantID,
			Name:          resource,
			Label:         label,
			Custom:        custom,
			Visible:       visible,
			EditableValue: editableValue,
			Searchable:    searchable,
			Type:          TypeForeignKey,
			Required:      required,
			SectionID:     sectionID,
			CategoryID:    categoryID,
			Created:       now,
			Updated:       now,
		},
		Resource: resource,
	}
}

func (fk foreignKey) GetBase() base {
	return fk.Base
}

func (fk foreignKey) GetDefaultValue() interface{} {
	return nil
}

func (fk foreignKey) DecodeAndValidateValue(value interface{}) (interface{}, error) {

	if !fk.Base.EditableValue {
		return nil, resperr.Error{}
	}

	b, err := json.Marshal(&value)
	if err != nil {
		return nil, err
	}

	var foreignKeyValue *primitive.ObjectID
	if err := json.Unmarshal(b, &foreignKeyValue); err != nil {
		return nil, fmt.Errorf("invalid value: %w", err)
	}

	if foreignKeyValue == nil {
		return fk.GetDefaultValue(), nil
	}

	return foreignKeyValue, nil
}
