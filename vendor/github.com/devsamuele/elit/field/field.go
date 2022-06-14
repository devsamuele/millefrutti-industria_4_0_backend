package field

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/devsamuele/elit/resperr"
	"github.com/devsamuele/elit/translation"
	"github.com/devsamuele/elit/utility/slice"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TypeText     = "text"      // string
	TypeRichText = "rich_text" // string
	TypeAddress  = "address"   // string
	TypeURL      = "url"       // string
	TypeEmail    = "email"     // []string
	TypePhone    = "phone"     // []string
	// FILE      = "file"      // string - objectID - uuid

	TypeBoolean = "boolean" // bool

	TypeForeignKey = "foreign_key" // ObjectID

	TypeNumber   = "number"   // float64
	TypeCurrency = "currency" // float64
	TypeDuration = "duration" // float64
	TypePercent  = "percent"  // float64
	TypeRating   = "rating"   // float64

	TypeDatetimeRange = "datetime_range" // object
	TypeDatetime      = "datetime"       // datetime

	TypeSingleSelect   = "single_select"   // fieldChoice
	TypeMultipleSelect = "multiple_select" // []fieldChoices

	TypeUser = "user" // []objectID
	//FieldTypeOrganization = "organization" // []objectID
	//FieldTypePerson       = "person"       // []objectID

	TypeCustomDataLink = "custom_data_link"
)

type Fielder interface {
	// UnmarshalBSON(data []byte) error
	// MarshalBSON() ([]byte, error)
	DecodeAndValidateValue(value interface{}) (interface{}, error)
	GetBase() base
	GetDefaultValue() interface{}
}

type ChoicePullRemover interface {
	RemoveChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error
	PullChoices(ctx context.Context, tenantID, fieldID primitive.ObjectID, choiceIDs []primitive.ObjectID, now time.Time) error
}

type IDChecker interface {
	CheckByID(ctx context.Context, tenantID, ID primitive.ObjectID) error
}

type LabelChecker interface {
	CheckByLabel(ctx context.Context, tenantID primitive.ObjectID, label translation.Translation) error
}

// type Builder struct {
// 	withSection          bool
// 	withCategory         bool
// 	validFieldTypes []string
// 	fkResource           []string
// 	store                Store
// }

// func NewBuilder(WithSection bool, WithCategory bool, ValidFieldTypes []string, ValidFKResource []string, store Store) Builder {
// 	b := Builder{
// 		withSection:          WithSection,
// 		withCategory:         WithCategory,
// 		validFieldTypes: ValidFieldTypes,
// 		fkResource:           ValidFKResource,
// 		store:                store,
// 	}

// 	return b
// }

func bsonDecoder(rawBson bson.Raw) (Fielder, error) {

	var field Fielder
	tt := struct {
		Type string `bson:"type"`
	}{}

	if err := bson.Unmarshal(rawBson, &tt); err != nil {
		return nil, err
	}

	switch tt.Type {
	case TypeText:
		var t text
		bson.Unmarshal(rawBson, &t)
		field = &t
	case TypeEmail:
		var e email
		bson.Unmarshal(rawBson, &e)
		field = &e
	case TypeSingleSelect:
		var ss singleSelect
		bson.Unmarshal(rawBson, &ss)
		field = &ss
	case TypeDatetimeRange:
		var dtr datetimeRange
		bson.Unmarshal(rawBson, &dtr)
		field = &dtr
	case TypeDatetime:
		var dt datetime
		bson.Unmarshal(rawBson, &dt)
		field = &dt
	case TypeForeignKey:
		var fk foreignKey
		bson.Unmarshal(rawBson, &fk)
		field = &fk
	}

	return field, nil
}

type base struct {
	ID            primitive.ObjectID      `json:"id" bson:"_id"`
	TenantID      primitive.ObjectID      `json:"tenant_id" bson:"tenant_id"`
	Name          string                  `json:"name,omitempty" bson:"name,omitempty"`
	Label         translation.Translation `json:"label" bson:"label"`
	Custom        bool                    `json:"custom" bson:"custom"`
	Visible       bool                    `json:"visible" bson:"visible"`
	EditableValue bool                    `json:"editable_value" bson:"editable_value"`
	Searchable    bool                    `json:"searchable" bson:"searchable"`
	Required      bool                    `json:"required" bson:"required"`
	Type          string                  `json:"type" bson:"type"`
	SectionID     *primitive.ObjectID     `json:"section_id,omitempty" bson:"section_id,omitempty"`
	CategoryID    *primitive.ObjectID     `json:"category_id,omitempty" bson:"category_id,omitempty"`
	Created       time.Time               `json:"created" bson:"created"`
	Updated       time.Time               `json:"updated" bson:"updated"`
}

type baseInput struct {
	Label         *translation.Translation `json:"label"`
	Visible       *bool                    `json:"visible"`
	EditableValue *bool                    `json:"editable_value"`
	Searchable    *bool                    `json:"searchable"`
	Type          *string                  `json:"type"`
	SectionID     *primitive.ObjectID      `json:"section_id,omitempty"`
	CategoryID    *primitive.ObjectID      `json:"category_id,omitempty"`
}

// func (bi *baseInput) UnmarshalJSON(b []byte) error {
// 	buffer := bytes.NewBuffer(b)
// 	decoder := json.NewDecoder(buffer)
// 	decoder.DisallowUnknownFields()

// 	type baseInputAlias baseInput
// 	bia := baseInputAlias(*bi)

// 	return decoder.Decode(&bia)
// }

type Choice struct {
	ID     primitive.ObjectID      `json:"id" bson:"id"`
	Value  translation.Translation `json:"value" bson:"value"`
	Custom bool                    `json:"custom" bson:"custom"`
}

type NewChoice struct {
	Value translation.Translation `json:"value"`
}

type UpdateChoice struct {
	ID    primitive.ObjectID      `json:"id"`
	Value translation.Translation `json:"value"`
}

type Input map[string]interface{}
type Update map[string]interface{}
type Values map[string]interface{}

// fieldTypes
// sectionStore
// categoryStore

type builderStageNew struct {
	ctx             context.Context
	tenantID        primitive.ObjectID
	now             time.Time
	validTypes      []string
	labelChecker    LabelChecker
	sectionChecker  IDChecker
	categoryChecker IDChecker
	jsonInput       Input
}

func (bs *builderStageNew) WithSection(sc IDChecker) *builderStageNew {
	bs.sectionChecker = sc
	return bs
}

func (bs *builderStageNew) WithCategory(cc IDChecker) *builderStageNew {
	bs.categoryChecker = cc
	return bs
}

func (bs *builderStageNew) Build() (Fielder, error) {

	jsonInputB, err := json.Marshal(bs.jsonInput)
	if err != nil {
		return nil, err
	}

	var bi baseInput
	if err := json.Unmarshal(jsonInputB, &bi); err != nil {
		return nil, err
	}

	var baseField = base{
		ID:       primitive.NewObjectID(),
		TenantID: bs.tenantID,
		Custom:   true,
		Required: false,
		Created:  bs.now,
		Updated:  bs.now,
	}

	if bi.Type == nil {
		// TODO respERR
		return nil, errors.New("type is required")
	}

	_, found := slice.Find(bs.validTypes, func(fieldType string) bool {
		return fieldType == *bi.Type
	})

	if !found {
		return nil, errors.New("field type not supported")
	}

	baseField.Type = *bi.Type

	if bi.Label == nil {
		// TODO respERR
		return nil, errors.New("label is required")
	}

	if bi.Label.IsEmpty() {
		return nil, resperr.Error{}
	}

	// * default true
	if bi.Visible == nil {
		baseField.Visible = true
	} else {
		baseField.Visible = *bi.Visible
	}

	// * default true
	if bi.EditableValue == nil {
		baseField.EditableValue = true
	} else {
		baseField.EditableValue = *bi.EditableValue
	}

	// * default true
	if bi.Searchable == nil {
		baseField.Searchable = true
	} else {
		baseField.Searchable = *bi.Searchable
	}

	if err := bs.labelChecker.CheckByLabel(bs.ctx, bs.tenantID, *bi.Label); err != nil {
		if errors.Is(err, ErrLabelExists) {
			return nil, resperr.Error{
				Message:      "label already exist",
				Reason:       resperr.ErrReasonConflict,
				LocationType: "argument",
				Location:     "label",
			}
		}
		return nil, err
	}
	baseField.Label = *bi.Label

	if bs.sectionChecker != nil {
		if bi.SectionID == nil {
			// TODO respERR
			return nil, errors.New("section_id is required")
		}
		if err := bs.sectionChecker.CheckByID(bs.ctx, bs.tenantID, *bi.SectionID); err != nil {
			if errors.Is(err, ErrNotFound) {
				return nil, resperr.Error{
					Message:      "section not found",
					Reason:       resperr.ErrReasonNotFound,
					LocationType: "argument",
					Location:     "section_id",
				}
			}
			return nil, err
		}

		baseField.SectionID = bi.SectionID
	}
	// else {
	// 	if bi.SectionID != nil {
	// 		return nil, errors.New("decoding error")
	// 	}
	// }

	if bs.categoryChecker != nil {
		if bi.CategoryID == nil {
			return nil, errors.New("category_id is required")
		}

		// TODO add categoryStore

		baseField.CategoryID = bi.CategoryID
	}
	// else {
	// 	if bi.CategoryID != nil {
	// 		return nil, errors.New("decoding error")
	// 	}
	// }

	switch baseField.Type {
	case TypeText:
		return text{Base: baseField}, nil

	case TypeEmail:
		return email(baseField), nil

	// case typeRichText:
	// case typeAddress:
	// case typeURL:
	// case typeBoolean:
	// case typeNumber:
	// case typeCurrency:
	// case typeDuration:
	// case typePercent:
	// case typeRating:
	case TypeDatetime:
		return datetime(baseField), nil

	case TypeDatetimeRange:
		return datetimeRange(baseField), nil

	// case TypeForeignKey:
	// 	return newForeignKey(b.validFKResource, jsonInputB, baseField)

	case TypeSingleSelect:
		return newSingleSelectFromInput(jsonInputB, baseField)

	// case typeMultipleSelect:
	// case typeUser:
	// case typeCustomDataLink:
	default:
		return nil, errors.New("invalid field type")
	}
}

func NewFromInput(ctx context.Context, tenantID primitive.ObjectID, lc LabelChecker, jsonInput Input, validTypes []string, now time.Time) *builderStageNew {
	bs := builderStageNew{
		ctx:          ctx,
		tenantID:     tenantID,
		now:          now,
		validTypes:   validTypes,
		labelChecker: lc,
		jsonInput:    jsonInput,
	}

	return &bs
}

type builderStageUpdate struct {
	ctx               context.Context
	tenantID          primitive.ObjectID
	now               time.Time
	validTypes        []string
	oldField          Fielder
	labelChecker      LabelChecker
	sectionChecker    IDChecker
	categoryChecker   IDChecker
	choicePullRemover ChoicePullRemover
	jsonUpdate        Update
}

type baseUpdate struct {
	Label         *translation.Translation `json:"label"`
	Visible       *bool                    `json:"visible"`
	EditableValue *bool                    `json:"editable_value"`
	Searchable    *bool                    `json:"searchable"`
	SectionID     *primitive.ObjectID      `json:"section_id,omitempty"`
	CategoryID    *primitive.ObjectID      `json:"category_id,omitempty"`
}

func (bsu *builderStageUpdate) WithSection(sc IDChecker) *builderStageUpdate {
	bsu.sectionChecker = sc
	return bsu
}

func (bsu *builderStageUpdate) WithCategory(cc IDChecker) *builderStageUpdate {
	bsu.categoryChecker = cc
	return bsu
}

func (bsu *builderStageUpdate) Update() (Fielder, error) {
	jsonUpdateB, err := json.Marshal(bsu.jsonUpdate)
	if err != nil {
		return nil, err
	}

	var bu baseUpdate
	if err := json.Unmarshal(jsonUpdateB, &bu); err != nil {
		return nil, err
	}

	// oldField, err := fieldStore.QueryByID(ctx, tenantID, fieldID)
	// if err != nil {
	// 	if errors.Is(err, ErrNotFound) {
	// 		return nil, resperr.Error{
	// 			Message:      "field not found",
	// 			Reason:       resperr.ErrReasonNotFound,
	// 			LocationType: "parameter",
	// 			Location:     "url",
	// 		}
	// 	}
	// 	return nil, err
	// }

	oldBaseField := bsu.oldField.GetBase()

	if bu.Label != nil {
		if err := bsu.labelChecker.CheckByLabel(bsu.ctx, bsu.tenantID, *bu.Label); err != nil {
			if errors.Is(err, ErrLabelExists) {
				return nil, resperr.Error{
					Message:      "label already exist",
					Reason:       resperr.ErrReasonConflict,
					LocationType: "argument",
					Location:     "label",
				}
			}
			return nil, err
		}
		oldBaseField.Label = *bu.Label
	}

	if bu.EditableValue != nil {
		oldBaseField.EditableValue = *bu.EditableValue
	}

	if bu.Searchable != nil {
		oldBaseField.Searchable = *bu.Searchable
	}

	if bu.Visible != nil {
		oldBaseField.Visible = *bu.Visible
	}

	if bsu.sectionChecker != nil {
		if bu.SectionID != nil {
			if err := bsu.sectionChecker.CheckByID(bsu.ctx, bsu.tenantID, *bu.SectionID); err != nil {
				if errors.Is(err, ErrNotFound) {
					return nil, resperr.Error{
						Message:      "section not found",
						Reason:       resperr.ErrReasonNotFound,
						LocationType: "argument",
						Location:     "section_id",
					}
				}
				return nil, err
			}
			oldBaseField.SectionID = bu.SectionID
		}
	} else {
		if bu.SectionID != nil {
			return nil, errors.New("decoding error")
		}
	}

	if bsu.categoryChecker != nil {
		if bu.CategoryID != nil {
			// TODO add categoryStore

			oldBaseField.CategoryID = bu.CategoryID
		}
	} else {
		if bu.CategoryID != nil {
			return nil, errors.New("decoding error")
		}
	}

	oldBaseField.Updated = bsu.now

	switch oldBaseField.Type {
	case TypeText:
		return text{
			Base: oldBaseField,
		}, nil

	case TypeEmail:
		return email(oldBaseField), nil

	// case typeRichText:
	// case typeAddress:
	// case typeURL:
	// case typeBoolean:
	// case typeNumber:
	// case typeCurrency:
	// case typeDuration:
	// case typePercent:
	// case typeRating:
	case TypeDatetime:
		return datetime(oldBaseField), nil

	case TypeDatetimeRange:
		return datetimeRange(oldBaseField), nil

	// case TypeForeignKey:
	// 	return updateForeignKey(b.validFKResource, jsonUpdateB, oldBaseField)

	case TypeSingleSelect:
		return updateSingleSelectFromInput(bsu.ctx, jsonUpdateB, bsu.oldField, oldBaseField, bsu.choicePullRemover, bsu.now)

	// case typeMultipleSelect:
	// case typeUser:
	// case typeCustomDataLink:
	default:
		return nil, errors.New("invalid field type")
	}
}

func UpdateFromInput(ctx context.Context, tenantID primitive.ObjectID, oldField Fielder, jsonUpdate Update, lc LabelChecker, cpr ChoicePullRemover, validTypes []string, now time.Time) *builderStageUpdate {
	bsu := builderStageUpdate{
		ctx:               ctx,
		tenantID:          tenantID,
		now:               now,
		validTypes:        validTypes,
		labelChecker:      lc,
		oldField:          oldField,
		choicePullRemover: cpr,
		jsonUpdate:        jsonUpdate,
	}

	return &bsu
}
