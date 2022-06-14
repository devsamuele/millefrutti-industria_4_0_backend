package field

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/devsamuele/elit/resperr"
	"github.com/devsamuele/elit/translation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type textResponse struct {
	base
	MinLenght int    `json:"min_lenght,omitempty"`
	MaxLenght int    `json:"max_length,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
}

type text struct {
	Base      base   `bson:"inline"`
	MinLenght int    `bson:"min_lenght,omitempty"`
	MaxLenght int    `bson:"max_length,omitempty"`
	Pattern   string `bson:"pattern,omitempty"`
}

func (t *text) UnmarshalBSON(data []byte) error {

	type textAlias text
	var ta textAlias

	if err := bson.Unmarshal(data, &ta); err != nil {
		return err
	}

	*t = text(ta)

	return nil
}

func (t *text) MarshalBSON() ([]byte, error) {

	type textAlias text
	ta := textAlias(*t)
	return bson.Marshal(&ta)
}

func (t *text) MarshalJSON() ([]byte, error) {

	type textAlias text
	ta := textAlias(*t)
	tr := textResponse{
		base:      ta.Base,
		MinLenght: ta.MinLenght,
		MaxLenght: ta.MaxLenght,
		Pattern:   ta.Pattern,
	}
	return json.Marshal(&tr)
}

func NewText(tenantID primitive.ObjectID, name string, label translation.Translation, custom, visible, editableValue, searchable, required bool, sectionID, categoryID *primitive.ObjectID, minLenght, maxLength int, pattern string, now time.Time) Fielder {
	return text{
		Base: base{
			ID:            primitive.NewObjectID(),
			TenantID:      tenantID,
			Name:          name,
			Label:         label,
			Custom:        custom,
			Visible:       visible,
			EditableValue: editableValue,
			Searchable:    searchable,
			Type:          TypeText,
			Required:      required,
			SectionID:     sectionID,
			CategoryID:    categoryID,
			Created:       now,
			Updated:       now,
		},
		MinLenght: minLenght,
		MaxLenght: maxLength,
		Pattern:   pattern,
	}
}

func (t text) GetBase() base {
	return t.Base
}

func (t text) GetDefaultValue() interface{} {
	return ""
}

func (t text) DecodeAndValidateValue(value interface{}) (interface{}, error) {

	if !t.Base.EditableValue {
		return nil, resperr.Error{}
	}

	b, err := json.Marshal(&value)
	if err != nil {
		return nil, err
	}

	var textValue *string
	if err := json.Unmarshal(b, &textValue); err != nil {
		return nil, fmt.Errorf("invalid value: %w", err)
	}

	if textValue == nil {
		return t.GetDefaultValue(), nil
	}

	err = resperr.Error{
		Message:      fmt.Sprintf("field[%v] is required", t.GetBase().ID.Hex()),
		Reason:       resperr.ErrReasonRequired,
		LocationType: "argument",
		Location:     fmt.Sprintf("field[%v]", t.GetBase().ID.Hex()),
	}

	if t.Base.Required {
		if t.MinLenght > 0 {
			if len(*textValue) < t.MinLenght {
				return nil, err
			}
		}

		if t.MaxLenght > 0 {
			if len(*textValue) > t.MaxLenght {
				return nil, err
			}
		}

		if t.Pattern != "" {
			matched, err := regexp.MatchString(t.Pattern, *textValue)
			if err != nil {
				return nil, err
			}
			if !matched {
				return nil, err
			}
		}
	}

	return textValue, nil
}
