package field

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/devsamuele/elit/resperr"
	"github.com/devsamuele/elit/translation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type email base

func NewEmail(tenantID primitive.ObjectID, name string, label translation.Translation, custom, visible, editableValue, searchable, required bool, sectionID, categoryID *primitive.ObjectID, now time.Time) Fielder {
	return email{
		ID:            primitive.NewObjectID(),
		TenantID:      tenantID,
		Name:          name,
		Label:         label,
		Custom:        custom,
		Visible:       visible,
		EditableValue: editableValue,
		Searchable:    searchable,
		Type:          TypeEmail,
		SectionID:     sectionID,
		CategoryID:    categoryID,
		Required:      required,
		Created:       now,
		Updated:       now,
	}
}

func (e email) GetBase() base {
	return base(e)
}

func (e email) GetDefaultValue() interface{} {
	return nil
}

type emailValue struct {
	Primary   *string  `json:"primary" bson:"primary"`
	Secondary []string `json:"secondary" bson:"secondary"`
}

func (e email) DecodeAndValidateValue(value interface{}) (interface{}, error) {

	if !e.EditableValue {
		return nil, resperr.Error{}
	}

	b, err := json.Marshal(&value)
	if err != nil {
		return nil, err
	}

	var ev *emailValue
	if err := json.Unmarshal(b, &ev); err != nil {
		return nil, fmt.Errorf("invalid value: %w", err)
	}

	if ev == nil {
		return e.GetDefaultValue(), nil
	}

	if ev.Primary == nil && ev.Secondary != nil {
		return nil, resperr.Error{Message: "secondary email without primary not valid"}
	}

	if ev.Primary != nil {
		ok, err := checkEmail(*ev.Primary)
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, resperr.Error{}
		}
	}

	if ev.Secondary == nil {
		ev.Secondary = make([]string, 0)
	}

	for _, e := range ev.Secondary {
		ok, err := checkEmail(e)
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, resperr.Error{}
		}
	}

	return ev, nil
}

func checkEmail(email string) (bool, error) {
	matched, err := regexp.MatchString(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, email)
	if err != nil {
		return false, err
	}

	return matched, nil
}
