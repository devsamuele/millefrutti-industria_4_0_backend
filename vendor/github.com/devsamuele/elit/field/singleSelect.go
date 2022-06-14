package field

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/devsamuele/elit/resperr"
	"github.com/devsamuele/elit/translation"
	"github.com/devsamuele/elit/utility/slice"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SingleSelect

type singleSelectInput struct {
	baseInput
	Choices []NewChoice `json:"choices"`
}

type singleSelectUpdate struct {
	baseUpdate
	Choices []UpdateChoice `json:"choices"`
}

type singleSelect struct {
	Base    base     `bson:",inline"`
	Choices []Choice `bson:"choices"`
}

type singleselectResponse struct {
	base
	Choices []Choice `json:"choices"`
}

func (ss *singleSelect) copy(f Fielder) error {
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, ss)
}

func (ss *singleSelect) UnmarshalBSON(data []byte) error {

	type singleSelectAlias singleSelect
	var ssa singleSelectAlias

	if err := bson.Unmarshal(data, &ssa); err != nil {
		return err
	}

	*ss = singleSelect(ssa)

	return nil
}

func (ss *singleSelect) MarshalBSON() ([]byte, error) {

	type singleSelectAlias singleSelect
	ssa := singleSelectAlias(*ss)
	return bson.Marshal(&ssa)
}

func (ss *singleSelect) MarshalJSON() ([]byte, error) {

	type singleSelectAlias singleSelect
	ssa := singleSelectAlias(*ss)
	ssr := singleselectResponse{
		base:    ssa.Base,
		Choices: ssa.Choices,
	}
	return json.Marshal(&ssr)
}

func (ss singleSelect) GetBase() base {
	return ss.Base
}

func (ss singleSelect) GetDefaultValue() interface{} {
	return nil
}

func NewSingleSelec(tenantID primitive.ObjectID, name string, label translation.Translation, custom, visible, editableValue, searchable, required bool, choices []Choice, sectionID, categoryID *primitive.ObjectID, now time.Time) Fielder {
	return singleSelect{
		Base: base{
			ID:            primitive.NewObjectID(),
			TenantID:      tenantID,
			Name:          name,
			Label:         label,
			Custom:        custom,
			Visible:       visible,
			EditableValue: editableValue,
			Searchable:    searchable,
			Type:          TypeSingleSelect,
			SectionID:     sectionID,
			CategoryID:    categoryID,
			Required:      required,
			Created:       now,
			Updated:       now,
		},
		Choices: choices,
	}
}
func (ss singleSelect) DecodeAndValidateValue(value interface{}) (interface{}, error) {

	if !ss.Base.EditableValue {
		return nil, resperr.Error{}
	}

	b, err := json.Marshal(&value)
	if err != nil {
		return nil, err
	}

	var choiceID *primitive.ObjectID
	if err := json.Unmarshal(b, &choiceID); err != nil {
		return nil, fmt.Errorf("invalid value: %w", err)
	}

	if choiceID == nil {
		return ss.GetDefaultValue(), nil
	}

	_, found := slice.Find(ss.Choices, func(choice Choice) bool {
		return choice.ID == *choiceID
	})
	if !found {
		return nil, resperr.Error{}
	}

	return choiceID, nil
}

func newSingleSelectFromInput(jsonInputB []byte, baseField base) (Fielder, error) {
	var ssi singleSelectInput
	if err := json.Unmarshal(jsonInputB, &ssi); err != nil {
		return nil, err
	}
	newChoices := make([]Choice, 0)
	for _, newChoice := range ssi.Choices {

		if newChoice.Value.IsEmpty() {
			return nil, resperr.Error{
				// TODO errs
			}
		}

		newChoices = append(newChoices, Choice{
			ID:     primitive.NewObjectID(),
			Value:  newChoice.Value,
			Custom: true,
		})
	}

	ss := singleSelect{
		Base:    baseField,
		Choices: newChoices,
	}

	return &ss, nil
}

func updateSingleSelectFromInput(ctx context.Context, jsonUpdateB []byte, oldField Fielder, oldBaseField base, cpr ChoicePullRemover, now time.Time) (Fielder, error) {

	if cpr == nil {
		return nil, errors.New("ChoicePullRemover not defined")
	}

	var ssu singleSelectUpdate
	if err := json.Unmarshal(jsonUpdateB, &ssu); err != nil {
		return nil, err
	}

	var oldSingleSelect singleSelect
	if err := oldSingleSelect.copy(oldField); err != nil {
		return nil, err
	}
	oldSingleSelect.Base = oldBaseField

	updatedChoices := make([]Choice, 0)
	if len(ssu.Choices) > 0 {
		for _, newChoice := range ssu.Choices {

			if newChoice.Value.IsEmpty() {
				return nil, resperr.Error{
					// TODO err
				}
			}

			if newChoice.ID == primitive.NilObjectID {
				updatedChoices = append(updatedChoices, Choice{
					ID:     primitive.NewObjectID(),
					Value:  newChoice.Value,
					Custom: true,
				})
			} else {
				oldChoice, found := slice.Find(oldSingleSelect.Choices, func(oldChoice Choice) bool {
					return oldChoice.ID == newChoice.ID
				})
				if found {
					if !oldChoice.Custom {
						// err: cannot update private choice
						return nil, resperr.Error{}
					}
					updatedChoices = append(updatedChoices, Choice{
						ID:     oldChoice.ID,
						Value:  newChoice.Value,
						Custom: oldChoice.Custom,
					})
				}
			}
		}

		removedChoices := slice.Filter(oldSingleSelect.Choices, func(oldChoice Choice) bool {
			_, found := slice.Find(updatedChoices, func(updatedChoice Choice) bool {
				return oldChoice.ID == updatedChoice.ID
			})
			if found {
				return false
			}

			if !found && !oldChoice.Custom {
				return false
			}

			return true
		})

		if len(removedChoices) > 0 {
			removedChoicesIDs := slice.Map(removedChoices, func(choice Choice) primitive.ObjectID {
				return choice.ID
			})

			if err := cpr.RemoveChoices(ctx, oldBaseField.TenantID, oldBaseField.ID, removedChoicesIDs, now); err != nil {
				return nil, err
			}
		}
	}

	ss := singleSelect{
		Base:    oldBaseField,
		Choices: updatedChoices,
	}

	return ss, nil

}
