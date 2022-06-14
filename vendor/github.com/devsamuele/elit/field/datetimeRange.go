package field

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/devsamuele/elit/resperr"
	"github.com/devsamuele/elit/translation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type datetimeRange base

func NewDatetimeRange(tenantID primitive.ObjectID, name string, label translation.Translation, custom, visible, editableValue, searchable, required bool, sectionID, categoryID *primitive.ObjectID, now time.Time) Fielder {
	return datetimeRange{
		ID:            primitive.NewObjectID(),
		TenantID:      tenantID,
		Name:          name,
		Label:         label,
		Custom:        custom,
		Visible:       visible,
		EditableValue: editableValue,
		Required:      required,
		Searchable:    searchable,
		Type:          TypeDatetimeRange,
		SectionID:     sectionID,
		CategoryID:    categoryID,
		Created:       now,
		Updated:       now,
	}
}

func (dtr datetimeRange) GetBase() base {
	return base(dtr)
}

func (dtr datetimeRange) GetDefaultValue() interface{} {
	return nil
}

type datetimeRangeValue struct {
	Start *time.Time `json:"start,omitempty" bson:"start,omitempty"`
	End   *time.Time `json:"end,omitempty" bson:"end,omitempty"`
}

func (dtr datetimeRange) DecodeAndValidateValue(value interface{}) (interface{}, error) {

	if !dtr.EditableValue {
		return nil, resperr.Error{}
	}

	b, err := json.Marshal(&value)
	if err != nil {
		return nil, err
	}

	var dtrv *datetimeRangeValue
	if err := json.Unmarshal(b, &dtrv); err != nil {
		return nil, fmt.Errorf("invalid value: %w", err)
	}

	if dtrv == nil || (dtrv.Start == nil && dtrv.End == nil) {
		return dtr.GetDefaultValue(), nil
	}

	if dtrv.Start != nil && dtrv.End != nil {
		if dtrv.End.Before(*dtrv.Start) {
			return nil, resperr.Error{}
		}
	}

	return dtrv, nil
}
