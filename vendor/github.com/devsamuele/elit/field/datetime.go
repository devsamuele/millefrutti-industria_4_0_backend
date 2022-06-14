package field

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/devsamuele/elit/resperr"
	"github.com/devsamuele/elit/translation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type datetime base

func NewDatetime(tenantID primitive.ObjectID, name string, label translation.Translation, custom, visible, editableValue, searchable, required bool, sectionID, categoryID *primitive.ObjectID, now time.Time) Fielder {
	return datetime{
		ID:            primitive.NewObjectID(),
		TenantID:      tenantID,
		Name:          name,
		Label:         label,
		Custom:        custom,
		Visible:       visible,
		EditableValue: editableValue,
		Searchable:    searchable,
		Required:      required,
		Type:          TypeDatetime,
		SectionID:     sectionID,
		CategoryID:    categoryID,
		Created:       now,
		Updated:       now,
	}
}

func (d datetime) GetBase() base {
	return base(d)
}

func (d datetime) GetDefaultValue() interface{} {
	return nil
}

func (d datetime) DecodeAndValidateValue(value interface{}) (interface{}, error) {

	if !d.EditableValue {
		return nil, resperr.Error{}
	}

	b, err := json.Marshal(&value)
	if err != nil {
		return nil, err
	}

	var dv *time.Time
	if err := json.Unmarshal(b, &dv); err != nil {
		return nil, fmt.Errorf("invalid value: %w", err)
	}

	if dv == nil {
		return d.GetDefaultValue(), nil
	}

	return dv, nil
}
