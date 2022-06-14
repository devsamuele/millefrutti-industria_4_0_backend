package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/devsamuele/service-kit/web"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	FieldOptionIntegerFormat = "integer"
	FieldOptionDecimalFormat = "decimal"
)

const (
	FieldOptionHMDurationFormat  = "h:mm"
	FieldOptionHMSDurationFormat = "h:mm:ss"
)

const (
	FieldOptionLocalDateFormat    = "local"
	FieldOptionFriendlyDateFormat = "friendly"
	FieldOptionUSDateFormat       = "US"
	FieldOptionEuropeanDateFormat = "European"
	FieldOptionISODateFormat      = "ISO"
)

const (
	FieldOptionMinNumberPrecision = 1
	FieldOptionMaxNumberPrecision = 8
)

const (
	FieldTypeText     = "text"      // string
	FieldTypeRichText = "rich_text" // string
	FieldTypeAddress  = "address"   // string
	FieldTypeURL      = "url"       // string
	//FieldTypePhone    = "phone"     // string
	/*FieldTypeEmail    = "email"     // string*/
	// FILE      = "file"      // string - objectID - uuid

	FieldTypeBoolean = "boolean" // bool

	FieldTypeNumber   = "number"   // float64
	FieldTypeCurrency = "currency" // float64
	FieldTypeDuration = "duration" // float64
	FieldTypePercent  = "percent"  // float64
	FieldTypeRating   = "rating"   // float64

	FieldTypeDatetime = "datetime" // datetime

	FieldTypeSingleSelect   = "single_select"   // fieldChoice
	FieldTypeMultipleSelect = "multiple_select" // []fieldChoices

	FieldTypeUser = "user" // []objectID
	//FieldTypeOrganization = "organization" // []objectID
	//FieldTypePerson       = "person"       // []objectID

	FieldTypeCustomDataLink = "custom_data_link"
)

type ObjectIDs struct {
	IDs []primitive.ObjectID `json:"IDs"`
}

type Field struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	TenantID      primitive.ObjectID `json:"tenant_id" bson:"tenant_id"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	Label         Translation        `json:"label" bson:"label"`
	Type          string             `json:"type" bson:"type"`
	Options       interface{}        `json:"options,omitempty" bson:"options,omitempty"`
	Choices       []Choice           `json:"choices,omitempty" bson:"choices,omitempty"`
	Custom        bool               `json:"custom" bson:"custom"`
	Visible       bool               `json:"visible" bson:"visible"`
	EditableValue bool               `json:"editable_value" bson:"editable_value"`
	Searchable    bool               `json:"searchable" bson:"searchable"`
	Created       time.Time          `json:"created" bson:"created"`
	Updated       time.Time          `json:"updated" bson:"updated"`
}

type Fields []Field

func (fs Fields) FindIDByName(fieldName string) (bool, primitive.ObjectID) {
	found := false
	for _, f := range fs {
		if f.Name == fieldName {
			found = true
			return found, f.ID
		}
	}
	return found, primitive.NilObjectID
}

func (fs Fields) ToMap() map[string]Field {
	fMap := make(map[string]Field)

	for _, bf := range fs {
		fMap[bf.ID.Hex()] = bf
	}

	return fMap
}

type FieldWithSection struct {
	Field     `bson:",inline"`
	SectionID primitive.ObjectID `json:"section_id" bson:"section_id"`
}

type FieldsWithSection []FieldWithSection

func (fs FieldsWithSection) ToField() []Field {
	fields := make([]Field, 0)
	for _, fws := range fs {
		fields = append(fields, fws.Field)
	}
	return fields
}

func (fs FieldsWithSection) ToMap() map[string]FieldWithSection {
	fMap := make(map[string]FieldWithSection)

	for _, f := range fs {
		fMap[f.ID.Hex()] = f
	}

	return fMap
}

func (fs FieldsWithSection) FindIDByName(fieldName string) (bool, primitive.ObjectID) {
	found := false
	for _, f := range fs {
		if f.Name == fieldName {
			found = true
			return found, f.ID
		}
	}
	return found, primitive.NilObjectID
}

type NewChoice struct {
	Value Translation `json:"value"`
}

type UpdateChoice struct {
	ID    *primitive.ObjectID `json:"id"`
	Value Translation         `json:"value"`
}

type NewField struct {
	Label      *Translation `json:"label"`
	Type       *string      `json:"type"`
	Options    interface{}  `json:"options"`
	Choices    []NewChoice  `json:"choices"`
	Searchable *bool        `json:"searchable"`
}

type NewFieldWithSection struct {
	NewField
	SectionID *primitive.ObjectID `json:"section_id"`
}

type UpdateField struct {
	Label      *Translation   `json:"label"`
	Options    interface{}    `json:"options"`
	Choices    []UpdateChoice `json:"choices"`
	Searchable *bool          `json:"searchable"`
}

type UpdateFieldWithSection struct {
	UpdateField
	SectionID *primitive.ObjectID `json:"section_id"`
}

type NumberTypeOptions struct {
	Format    string   `json:"format" bson:"format"`       // decimal | integer | currency | percent | duration -> validation on list
	Precision int      `json:"precision" bson:"precision"` // 1|2|...validation -> num between 1 to n (only decimal, percent, currency)
	Negative  bool     `json:"negative" bson:"negative"`
	Default   *float64 `json:"default" bson:"default"`
}

type CurrencyTypeOptions struct {
	Symbol    string   `json:"format" bson:"symbol"`       // whatever -> no validation
	Precision int      `json:"precision" bson:"precision"` // 1|2|...validation -> num between 1 to n (only decimal, percent, currency)
	Negative  bool     `json:"negative" bson:"negative"`   //
	Default   *float64 `json:"default" bson:"default"`
}

type PercentTypeOptions struct {
	Precision int      `json:"precision" bson:"precision"` // 1|2|...validation -> num between 1 to n (only decimal, percent, currency)
	Negative  bool     `json:"negative" bson:"negative"`
	Default   *float64 `json:"default" bson:"default"`
}

type TextTypeOptions struct {
	Default *string `json:"default" bson:"default"`
}

type DurationTypeOptions struct {
	Format  string   `json:"format" bson:"format"`
	Default *float64 `json:"default" bson:"default"`
}

type DateTimeTypeOptions struct {
	Format  string     `json:"format" bson:"format"`
	Default *time.Time `json:"default" bson:"default"`
}

type RatingTypeOptions struct {
	Color string `json:"color" bson:"color"`
	Icon  string `json:"icon" bson:"icon"`
	Max   int    `json:"max" bson:"max"`
}

type SelectTypeOptions struct {
	InsertChoices bool `json:"insert_choices" bson:"insert_choices"`
}

type Choice struct {
	ID     primitive.ObjectID `json:"id" bson:"id"`
	Value  Translation        `json:"value" bson:"value"`
	Custom bool               `json:"custom" bson:"custom"`
}

type NewNumberTypeOptions struct {
	Format    *string  `json:"format"`    // decimal | integer | currency | percent | duration -> validation on list
	Precision *int     `json:"precision"` // 1|2|...validation -> num between 1 to n (only decimal, percent, currency)
	Negative  *bool    `json:"negative"`
	Default   *float64 `json:"default"`
}

func (nto NewNumberTypeOptions) toModelTypeOptions() NumberTypeOptions {
	return NumberTypeOptions{
		Format:    *nto.Format,
		Precision: *nto.Precision,
		Negative:  *nto.Negative,
		Default:   nto.Default,
	}
}

type UpdateNumberTypeOptions struct {
	Format    *string  `json:"format"`    // decimal | integer | currency | percent | duration -> validation on list
	Precision *int     `json:"precision"` // 1|2|...validation -> num between 1 to n (only decimal, percent, currency)
	Negative  *bool    `json:"negative"`
	Default   *float64 `json:"default"`
}

type NewCurrencyTypeOptions struct {
	Symbol    *string  `json:"format"`    // whatever -> no validation
	Precision *int     `json:"precision"` // 1|2|...validation -> num between 1 to n (only decimal, percent, currency)
	Negative  *bool    `json:"negative"`  //
	Default   *float64 `json:"default"`
}

func (cto NewCurrencyTypeOptions) toModelTypeOptions() CurrencyTypeOptions {
	return CurrencyTypeOptions{
		Symbol:    *cto.Symbol,
		Precision: *cto.Precision,
		Negative:  *cto.Negative,
		Default:   cto.Default,
	}
}

type UpdateCurrencyTypeOptions struct {
	Symbol    *string  `json:"format"`    // whatever -> no validation
	Precision *int     `json:"precision"` // 1|2|...validation -> num between 1 to n (only decimal, percent, currency)
	Negative  *bool    `json:"negative"`  //
	Default   *float64 `json:"default"`
}

type NewPercentTypeOptions struct {
	Precision *int     `json:"precision"` // 1|2|...validation -> num between 1 to n (only decimal, percent, currency)
	Negative  *bool    `json:"negative"`
	Default   *float64 `json:"default"`
}

func (pto NewPercentTypeOptions) toModelTypeOptions() PercentTypeOptions {
	return PercentTypeOptions{
		Precision: *pto.Precision,
		Negative:  *pto.Negative,
		Default:   pto.Default,
	}
}

type UpdatePercentTypeOptions struct {
	Precision *int     `json:"precision"` // 1|2|...validation -> num between 1 to n (only decimal, percent, currency)
	Negative  *bool    `json:"negative"`
	Default   *float64 `json:"default"`
}

type NewTextTypeOptions struct {
	Default *string `json:"default"`
}

func (tto NewTextTypeOptions) toModelTypeOptions() TextTypeOptions {
	return TextTypeOptions{
		Default: tto.Default,
	}
}

type UpdateTextTypeOptions struct {
	Default *string `json:"default"`
}

type NewDurationTypeOptions struct {
	Format  *string  `json:"format"`
	Default *float64 `json:"default"`
}

func (dto NewDurationTypeOptions) toModelTypeOptions() DurationTypeOptions {
	return DurationTypeOptions{
		Format:  *dto.Format,
		Default: dto.Default,
	}
}

type UpdateDurationTypeOptions struct {
	Format  *string  `json:"format"`
	Default *float64 `json:"default"`
}

type NewDateTimeTypeOptions struct {
	Format  *string    `json:"format"`
	Default *time.Time `json:"default"`
}

func (dto NewDateTimeTypeOptions) toModelTypeOptions() DateTimeTypeOptions {
	return DateTimeTypeOptions{
		Format:  *dto.Format,
		Default: dto.Default,
	}
}

type UpdateDateTimeTypeOptions struct {
	Format  *string    `json:"format"`
	Default *time.Time `json:"default"`
}

type NewRatingTypeOptions struct {
	Color *string `json:"color"`
	Icon  *string `json:"icon"`
	Max   *int    `json:"max"`
}

func (rto NewRatingTypeOptions) toModelTypeOptions() RatingTypeOptions {
	return RatingTypeOptions{
		Color: *rto.Color,
		Icon:  *rto.Icon,
		Max:   *rto.Max,
	}
}

type UpdateRatingTypeOptions struct {
	Color *string `json:"color"`
	Icon  *string `json:"icon"`
	Max   *int    `json:"max"`
}

type FieldValues map[string]interface{}

func (nf NewField) BuildChoices() []Choice {
	if len(nf.Choices) > 0 {
		newChoices := make([]Choice, 0)
		for _, newChoice := range nf.Choices {
			newChoices = append(newChoices, Choice{
				ID:     primitive.NewObjectID(),
				Value:  newChoice.Value,
				Custom: true,
			})
		}
		return newChoices
	}
	return nil
}

func (nf NewField) BuildOptions() interface{} {
	switch *nf.Type {
	case FieldTypeNumber:
		newOpts := nf.Options.(NewNumberTypeOptions)
		return newOpts.toModelTypeOptions()

	case FieldTypeCurrency:
		newOpts := nf.Options.(NewCurrencyTypeOptions)
		return newOpts.toModelTypeOptions()

	case FieldTypePercent:
		newOpts := nf.Options.(NewPercentTypeOptions)
		return newOpts.toModelTypeOptions()

	case FieldTypeDuration:
		newOpts := nf.Options.(NewDurationTypeOptions)
		return newOpts.toModelTypeOptions()

	case FieldTypeRating:
		newOpts := nf.Options.(NewRatingTypeOptions)
		return newOpts.toModelTypeOptions()

	case FieldTypeDatetime:
		newOpts := nf.Options.(NewDateTimeTypeOptions)
		return newOpts.toModelTypeOptions()

	case FieldTypeText:
		newOpts := nf.Options.(NewTextTypeOptions)
		return newOpts.toModelTypeOptions()

	case FieldTypeSingleSelect, FieldTypeMultipleSelect:
		return SelectTypeOptions{InsertChoices: true}

	default:
		return nil
	}
}

func (uf UpdateField) UpdateChoices(oldField Field) ([]Choice, []primitive.ObjectID) {

	newChoices := make([]Choice, 0)
	for _, choice := range uf.Choices {
		if choice.ID == nil {
			newChoices = append(newChoices, Choice{
				ID:     primitive.NewObjectID(),
				Value:  choice.Value,
				Custom: true,
			})
		} else {
			found := false
			for i := 0; i < len(oldField.Choices) && !found; i++ {
				if choice.ID.Hex() == oldField.Choices[i].ID.Hex() {
					newChoices = append(newChoices, Choice{
						ID:     oldField.Choices[i].ID,
						Value:  choice.Value,
						Custom: oldField.Choices[i].Custom,
					})
					found = true
				}
			}
		}
	}

	removedChoiceIDs := make([]primitive.ObjectID, 0)
	for _, oldChoice := range oldField.Choices {
		if !oldChoice.Custom {
			newChoices = append(newChoices, oldChoice)
		}
		found := false
		for i := 0; i < len(newChoices) && !found; i++ {
			if oldChoice.ID.Hex() == newChoices[i].ID.Hex() {
				found = true
			}
		}
		if !found {
			removedChoiceIDs = append(removedChoiceIDs, oldChoice.ID)
		}
	}

	return newChoices, removedChoiceIDs
}

func (uf UpdateField) UpdateOptions(oldField Field) interface{} {
	switch oldField.Type {
	case FieldTypeNumber:
		opts := oldField.Options.(NumberTypeOptions)
		updateOpts := uf.Options.(UpdateNumberTypeOptions)

		if updateOpts.Format != nil {
			opts.Format = *updateOpts.Format
		}

		if updateOpts.Precision != nil {
			opts.Precision = *updateOpts.Precision
		}

		if updateOpts.Default != nil {
			opts.Default = updateOpts.Default
		}

		return opts

	case FieldTypeCurrency:
		opts := oldField.Options.(CurrencyTypeOptions)
		var updateOpts UpdateCurrencyTypeOptions

		if updateOpts.Precision != nil {
			opts.Precision = *updateOpts.Precision
		}

		if updateOpts.Symbol != nil {
			opts.Symbol = *updateOpts.Symbol
		}

		if updateOpts.Default != nil {
			opts.Default = updateOpts.Default
		}

		return opts

	case FieldTypePercent:
		opts := oldField.Options.(PercentTypeOptions)
		var updateOpts UpdatePercentTypeOptions

		if updateOpts.Precision != nil {
			opts.Precision = *updateOpts.Precision
		}

		if updateOpts.Default != nil {
			opts.Default = updateOpts.Default
		}

		return opts

	case FieldTypeDuration:
		opts := oldField.Options.(DurationTypeOptions)
		var updateOpts UpdateDurationTypeOptions

		if updateOpts.Format != nil {
			opts.Format = *updateOpts.Format
		}

		if updateOpts.Default != nil {
			opts.Default = updateOpts.Default
		}

		return opts

	case FieldTypeRating:
		opts := oldField.Options.(RatingTypeOptions)
		var updateOpts UpdateRatingTypeOptions

		if updateOpts.Icon != nil {
			opts.Icon = *updateOpts.Icon
		}

		if updateOpts.Max != nil {
			opts.Max = *updateOpts.Max
		}

		if updateOpts.Color != nil {
			opts.Color = *updateOpts.Color
		}

		return opts

	case FieldTypeDatetime:
		opts := oldField.Options.(DateTimeTypeOptions)
		var updateOpts UpdateDateTimeTypeOptions

		if updateOpts.Format != nil {
			opts.Format = *updateOpts.Format
		}

		if updateOpts.Default != nil {
			opts.Default = updateOpts.Default
		}

		return opts

	case FieldTypeText:
		opts := oldField.Options.(TextTypeOptions)
		var updateOpts UpdateTextTypeOptions

		if updateOpts.Default != nil {
			opts.Default = updateOpts.Default
		}

		return opts

	default:
		return nil
	}
}

func validFieldOptionPrecision(precision int) bool {
	return precision <= FieldOptionMaxNumberPrecision && precision >= FieldOptionMinNumberPrecision
}

func validFieldOptionDateTimeFormat(format string) bool {
	return format == FieldOptionLocalDateFormat ||
		format == FieldOptionFriendlyDateFormat ||
		format == FieldOptionUSDateFormat ||
		format == FieldOptionEuropeanDateFormat ||
		format == FieldOptionISODateFormat
}

func validFieldOptionDurationFormat(format string) bool {
	return format == FieldOptionHMDurationFormat ||
		format == FieldOptionHMSDurationFormat

}

func validFieldOptionNumberFormat(format string) bool {
	return format == FieldOptionIntegerFormat ||
		format == FieldOptionDecimalFormat

}

func IsFieldWithArrayValue(fieldType string) bool {
	return fieldType == FieldTypeMultipleSelect || fieldType == FieldTypeUser
}

func IsSelectFieldType(fieldType string) bool {
	return fieldType == FieldTypeMultipleSelect || fieldType == FieldTypeSingleSelect
}

func (f *Field) UnmarshalBSON(data []byte) error {

	type AuxField Field
	var auxField AuxField

	if err := bson.Unmarshal(data, &auxField); err != nil {
		return err
	}

	if auxField.Options == nil {
		auxField.Options = nil
		*f = Field(auxField)
		return nil
	}

	b, err := bson.Marshal(auxField.Options)
	if err != nil {
		return err
	}

	switch auxField.Type {
	case FieldTypeNumber:
		var opts NumberTypeOptions
		if err := bson.Unmarshal(b, &opts); err != nil {
			return err
		}
		auxField.Options = opts

	case FieldTypeCurrency:
		var opts CurrencyTypeOptions
		if err := bson.Unmarshal(b, &opts); err != nil {
			return err
		}
		auxField.Options = opts

	case FieldTypePercent:
		var opts PercentTypeOptions
		if err := bson.Unmarshal(b, &opts); err != nil {
			return err
		}
		auxField.Options = opts

	case FieldTypeDuration:
		var opts DurationTypeOptions
		if err := bson.Unmarshal(b, &opts); err != nil {
			return err
		}
		auxField.Options = opts

	case FieldTypeRating:
		var opts RatingTypeOptions
		if err := bson.Unmarshal(b, &opts); err != nil {
			return err
		}
		auxField.Options = opts

	case FieldTypeDatetime:
		var opts DateTimeTypeOptions
		if err := bson.Unmarshal(b, &opts); err != nil {
			return err
		}
		auxField.Options = opts

	case FieldTypeText:
		var opts TextTypeOptions
		if err := bson.Unmarshal(b, &opts); err != nil {
			return err
		}
		auxField.Options = opts

	case FieldTypeSingleSelect, FieldTypeMultipleSelect:
		var opts SelectTypeOptions
		if err := bson.Unmarshal(b, &opts); err != nil {
			return err
		}
		auxField.Options = opts

	default:
		auxField.Options = nil
	}

	*f = Field(auxField)

	return nil
}

func (fws *FieldWithSection) UnmarshalBSON(data []byte) error {

	if err := fws.Field.UnmarshalBSON(data); err != nil {
		return err
	}

	var auxSectionID struct {
		SectionID primitive.ObjectID `bson:"section_id"`
	}

	if err := bson.Unmarshal(data, &auxSectionID); err != nil {
		return err
	}

	fws.SectionID = auxSectionID.SectionID

	return nil
}

func (nf *NewFieldWithSection) Validate(validTypes []string) error {

	if nf.SectionID == nil {
		return Error{
			Message:      "section is required",
			Reason:       ErrReasonRequired,
			LocationType: "argument",
			Location:     "section_id",
		}
	}

	if err := nf.NewField.Validate(validTypes); err != nil {
		return err
	}
	return nil
}

func (uf *UpdateFieldWithSection) Validate(oldField FieldWithSection) error {
	return uf.UpdateField.Validate(oldField.Field)
}

func (nf *NewField) Validate(validTypes []string) error {

	if nf.Type == nil {
		return Error{
			Message:      "type is required",
			Reason:       ErrReasonRequired,
			LocationType: "argument",
			Location:     "type",
		}
	}

	found := false
	for i := 0; i < len(validTypes) && !found; i++ {
		if *nf.Type == validTypes[i] {
			found = true
		}
	}
	if !found {
		return Error{
			Message:      "invalid type",
			Reason:       ErrReasonInvalidArgument,
			LocationType: "argument",
			Location:     "type",
		}
	}

	if nf.Label == nil || (nf.Label != nil && (*nf.Label).IsEmpty()) {
		return Error{
			Message:      "label is required",
			Reason:       ErrReasonRequired,
			LocationType: "argument",
			Location:     "label",
		}
	}

	if !(*nf.Label).IsValidLength(200) {
		return Error{
			Message:      "max length exceeded",
			Reason:       ErrReasonInvalidArgument,
			LocationType: "argument",
			Location:     "title",
		}
	}

	if !IsSelectFieldType(*nf.Type) && len(nf.Choices) > 0 {
		return Error{
			Message:      "choices not supported",
			Reason:       ErrReasonInvalidArgument,
			LocationType: "argument",
			Location:     "choices",
		}
	}

	if len(nf.Choices) > 0 {
		if len(nf.Choices) > 100 {
			return Error{
				Message:      "max length exceeded",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "choices",
			}
		}
		for _, newChoice := range nf.Choices {
			if newChoice.Value.IsEmpty() {
				return Error{
					Message:      "empty choice value",
					Reason:       ErrReasonInvalidArgument,
					LocationType: "argument",
					Location:     "choices",
				}
			}
			if newChoice.Value.IsValidLength(200) {
				return Error{
					Message:      "max value length exceeded",
					Reason:       ErrReasonInvalidArgument,
					LocationType: "argument",
					Location:     "choices",
				}
			}
		}
	}

	b, err := json.Marshal(nf.Options)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewBuffer(b))
	decoder.DisallowUnknownFields()

	switch *nf.Type {
	case FieldTypeNumber:
		var newOpts NewNumberTypeOptions
		if err := decoder.Decode(&newOpts); err != nil {
			return err
		}

		if newOpts.Format == nil {
			return Error{
				Message:      "format is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "options.format",
			}
		}

		if newOpts.Precision == nil {
			return Error{
				Message:      "precision is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "options.precision",
			}
		}

		if newOpts.Negative == nil {
			return Error{
				Message:      "negative is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "options.negative",
			}
		}

		if !validFieldOptionNumberFormat(*newOpts.Format) {
			return Error{
				Message:      "invalid format",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.format",
			}
		}

		if !validFieldOptionPrecision(*newOpts.Precision) {
			return Error{
				Message:      "invalid precision",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.precision",
			}
		}

		nf.Options = newOpts

	case FieldTypeCurrency:
		var newOpts NewCurrencyTypeOptions
		if err := decoder.Decode(&newOpts); err != nil {
			return err
		}

		if newOpts.Symbol == nil {
			return Error{
				Message:      "symbol is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "options.symbol",
			}
		}

		if newOpts.Precision == nil {
			return Error{
				Message:      "precision is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "options.precision",
			}
		}

		if newOpts.Negative == nil {
			return Error{
				Message:      "negative is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "options.negative",
			}
		}

		if !validFieldOptionPrecision(*newOpts.Precision) {
			return Error{
				Message:      "invalid precision",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.precision",
			}
		}

		nf.Options = newOpts

	case FieldTypePercent:
		var newOpts NewPercentTypeOptions
		if err := decoder.Decode(&newOpts); err != nil {
			return err
		}

		if newOpts.Precision == nil {
			return Error{
				Message:      "precision is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "options.precision",
			}
		}

		if newOpts.Negative == nil {
			return Error{
				Message:      "negative is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "options.negative",
			}
		}

		if !validFieldOptionPrecision(*newOpts.Precision) {
			return Error{
				Message:      "invalid precision",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.precision",
			}
		}

		nf.Options = newOpts

	case FieldTypeDuration:
		var newOpts NewDurationTypeOptions
		if err := decoder.Decode(&newOpts); err != nil {
			return err
		}

		if newOpts.Format == nil {
			return Error{
				Message:      "format is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "options.format",
			}
		}

		if !validFieldOptionDurationFormat(*newOpts.Format) {
			return Error{
				Message:      "invalid format",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.format",
			}
		}

		nf.Options = newOpts

	case FieldTypeRating:
		var newOpts NewRatingTypeOptions
		if err := decoder.Decode(&newOpts); err != nil {
			return err
		}

		if newOpts.Icon == nil {
			return Error{
				Message:      "icon is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "options.icon",
			}
		}

		if newOpts.Max == nil {
			return Error{
				Message:      "max is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "options.max",
			}
		}

		if newOpts.Color == nil {
			return Error{
				Message:      "color is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "options.color",
			}
		}

		if *newOpts.Max <= 0 {
			return Error{
				Message:      "invalid max value. Must be a positive number",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.max",
			}

		}

		nf.Options = newOpts

	case FieldTypeDatetime:
		var newOpts NewDateTimeTypeOptions
		if err := decoder.Decode(&newOpts); err != nil {
			return err
		}

		if newOpts.Format == nil {
			return Error{
				Message:      "format is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "options.format",
			}
		}

		if !validFieldOptionDateTimeFormat(*newOpts.Format) {
			return Error{
				Message:      "invalid format",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.format",
			}
		}

		nf.Options = newOpts

	case FieldTypeText:
		var newOpts NewTextTypeOptions
		if err := decoder.Decode(&newOpts); err != nil {
			return err
		}

		nf.Options = newOpts

	default:
		nf.Options = nil
	}

	return nil
}

func (uf *UpdateField) Validate(oldField Field) error {

	if uf.Label != nil {
		if (*uf.Label).IsEmpty() {
			return Error{
				Message:      "label is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "label",
			}
		}

		if !(*uf.Label).IsValidLength(200) {
			return Error{
				Message:      "max length exceeded",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "title",
			}
		}
	}

	if !IsSelectFieldType(oldField.Type) && len(uf.Choices) > 0 {
		return Error{
			Message:      "choices not supported",
			Reason:       ErrReasonInvalidArgument,
			LocationType: "argument",
			Location:     "choices",
		}
	}

	if len(uf.Choices) > 0 {
		selectOptions := oldField.Options.(SelectTypeOptions)
		if !selectOptions.InsertChoices {
			return Error{
				Message:      "choices not editable",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "choices",
			}
		}

		if len(uf.Choices) > 100 {
			return Error{
				Message:      "max length exceeded",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "choices",
			}
		}

		for _, choice := range uf.Choices {
			if (choice.Value).IsEmpty() {
				return Error{
					Message:      "empty choice value",
					Reason:       ErrReasonInvalidArgument,
					LocationType: "argument",
					Location:     "choices",
				}
			}
			if !choice.Value.IsValidLength(200) {
				return Error{
					Message:      "max value length exceeded",
					Reason:       ErrReasonInvalidArgument,
					LocationType: "argument",
					Location:     "choices",
				}
			}

			if choice.ID != nil {
				found := false
				for i := 0; i < len(oldField.Choices) && !found; i++ {
					if oldField.Choices[i].ID.Hex() == choice.ID.Hex() {
						if !oldField.Choices[i].Custom {
							return Error{
								Message:      "choice not editable",
								Reason:       ErrReasonInvalidArgument,
								LocationType: "argument",
								Location:     "choices",
							}
						}
						found = true
					}
				}
				if !found {
					return Error{
						Message:      "choice not found",
						Reason:       ErrReasonNotFound,
						LocationType: "argument",
						Location:     "choices",
					}
				}
			}
		}
	}

	if uf.Options != nil && !oldField.Custom {
		return Error{
			Message:      "unable to update private options",
			Reason:       ErrReasonInvalidArgument,
			LocationType: "argument",
			Location:     "options",
		}
	}

	b, err := json.Marshal(uf.Options)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewBuffer(b))
	decoder.DisallowUnknownFields()

	switch oldField.Type {
	case FieldTypeNumber:
		var updateOpts UpdateNumberTypeOptions
		if err := decoder.Decode(&updateOpts); err != nil {
			return err
		}
		if updateOpts.Format != nil && !validFieldOptionNumberFormat(*updateOpts.Format) {
			return Error{
				Message:      "invalid format",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.format",
			}
		}

		if updateOpts.Precision != nil && !validFieldOptionPrecision(*updateOpts.Precision) {
			return Error{
				Message:      "invalid precision",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.precision",
			}
		}

		if updateOpts.Negative != nil {
			return Error{
				Message:      "default value can not be update",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.default",
			}
		}

		uf.Options = updateOpts

	case FieldTypeCurrency:
		var updateOpts UpdateCurrencyTypeOptions
		if err := decoder.Decode(&updateOpts); err != nil {
			return err
		}

		if updateOpts.Precision != nil && !validFieldOptionPrecision(*updateOpts.Precision) {
			return Error{
				Message:      "invalid precision",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.precision",
			}
		}

		if updateOpts.Negative != nil {
			return Error{
				Message:      "default value can not be update",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.default",
			}
		}

		uf.Options = updateOpts

	case FieldTypePercent:
		var updateOpts UpdatePercentTypeOptions
		if err := decoder.Decode(&updateOpts); err != nil {
			return err
		}

		if updateOpts.Precision != nil && !validFieldOptionPrecision(*updateOpts.Precision) {
			return Error{
				Message:      "invalid precision",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.precision",
			}
		}

		if updateOpts.Negative != nil {
			return Error{
				Message:      "default value can not be update",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.default",
			}
		}

		uf.Options = updateOpts

	case FieldTypeDuration:
		var updateOpts UpdateDurationTypeOptions
		if err := decoder.Decode(&updateOpts); err != nil {
			return err
		}

		if updateOpts.Format != nil && !validFieldOptionDurationFormat(*updateOpts.Format) {
			return Error{
				Message:      "invalid format",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.format",
			}
		}

		uf.Options = updateOpts

	case FieldTypeRating:
		var updateOpts UpdateRatingTypeOptions
		if err := decoder.Decode(&updateOpts); err != nil {
			return err
		}

		if updateOpts.Max != nil && *updateOpts.Max <= 0 {
			return Error{
				Message:      "invalid max value. Must be a positive number",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.max",
			}
		}

		uf.Options = updateOpts

	case FieldTypeDatetime:
		var updateOpts UpdateDateTimeTypeOptions
		if err := decoder.Decode(&updateOpts); err != nil {
			return err
		}

		if updateOpts.Format != nil && !validFieldOptionDateTimeFormat(*updateOpts.Format) {
			return Error{
				Message:      "invalid format",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "options.format",
			}
		}

		uf.Options = updateOpts

	case FieldTypeText:
		var updateOpts UpdateTextTypeOptions
		if err := decoder.Decode(&updateOpts); err != nil {
			return err
		}

		uf.Options = updateOpts

	default:
		uf.Options = nil
	}

	return nil
}

// FieldValue

func (fv FieldValues) Validate(fields []Field) (FieldValues, error) {

	fieldsMap := Fields(fields).ToMap()
	if fv == nil {
		return make(FieldValues, 0), nil
	}

	for newFieldID := range fv {
		fieldID, err := primitive.ObjectIDFromHex(newFieldID)
		if err != nil {
			if errors.Is(err, primitive.ErrInvalidHex) {
				return nil, Error{
					Message:      "ID is not in its proper form",
					Reason:       ErrReasonInvalidArgument,
					LocationType: "argument",
					Location:     fmt.Sprintf("field_values.%v", fieldID),
				}
			}
			return nil, err
		}

		_, ok := fieldsMap[newFieldID]
		if !ok {
			return nil, Error{
				Message:      "field not found",
				Reason:       ErrReasonNotFound,
				LocationType: "argument",
				Location:     fmt.Sprintf("field[%v]", newFieldID),
			}
		}
	}

	fieldValues := make(map[string]interface{})
	for _, f := range fieldsMap {
		newFieldValue := fv[f.ID.Hex()]

		if !f.EditableValue {
			return nil, Error{
				Message:      "field not editable",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "",
				Location:     "",
			}
		}

		switch f.Type {
		case FieldTypeNumber:

			if newFieldValue == nil {
				fieldValues[f.ID.Hex()] = nil
			}

			b, err := json.Marshal(newFieldValue)
			if err != nil {
				return nil, err
			}

			opts := f.Options.(NumberTypeOptions)

			var number float64
			if err := json.Unmarshal(b, &number); err != nil {
				return nil, err
			}

			if !opts.Negative && number < 0 {
				return nil, Error{
					Message:      "value must be greater than 0",
					Reason:       ErrReasonInvalidArgument,
					LocationType: "argument",
					Location:     fmt.Sprintf("field_values.%v", f.ID.Hex()),
				}
			}

			fieldValues[f.ID.Hex()] = number

		case FieldTypeCurrency:

			fieldValues[f.ID.Hex()] = nil

			b, err := json.Marshal(newFieldValue)
			if err != nil {
				return nil, err
			}

			opts := f.Options.(CurrencyTypeOptions)

			var currency float64
			if err := json.Unmarshal(b, &currency); err != nil {
				return nil, err
			}

			if !opts.Negative && currency < 0 {
				return nil, Error{
					Message:      "value must be greater than 0",
					Reason:       ErrReasonInvalidArgument,
					LocationType: "argument",
					Location:     fmt.Sprintf("field_values.%v", f.ID.Hex()),
				}
			}

			fieldValues[f.ID.Hex()] = currency

		case FieldTypePercent:

			fieldValues[f.ID.Hex()] = nil

			b, err := json.Marshal(newFieldValue)
			if err != nil {
				return nil, err
			}

			opts := f.Options.(PercentTypeOptions)

			var percent float64
			if err := json.Unmarshal(b, &percent); err != nil {
				return nil, err
			}

			if !opts.Negative && percent < 0 {
				return nil, Error{
					Message:      "value must be greater than 0",
					Reason:       ErrReasonInvalidArgument,
					LocationType: "argument",
					Location:     fmt.Sprintf("field_values.%v", f.ID.Hex()),
				}
			}

			fieldValues[f.ID.Hex()] = percent

		case FieldTypeDuration:

			fieldValues[f.ID.Hex()] = nil

			b, err := json.Marshal(newFieldValue)
			if err != nil {
				return nil, err
			}

			var duration float64
			if err := json.Unmarshal(b, &duration); err != nil {
				return nil, err
			}

			fieldValues[f.ID.Hex()] = duration

		case FieldTypeRating:

			fieldValues[f.ID.Hex()] = nil

			b, err := json.Marshal(newFieldValue)
			if err != nil {
				return nil, err
			}

			opts := f.Options.(RatingTypeOptions)
			var rating float64
			if err := json.Unmarshal(b, &rating); err != nil {
				return nil, err
			}

			if rating > float64(opts.Max) || rating < 0 {
				return nil, Error{
					Message:      fmt.Sprintf("value must be greater than 0 or less than equal %v", opts.Max),
					Reason:       ErrReasonInvalidArgument,
					LocationType: "argument",
					Location:     fmt.Sprintf("field_values.%v", f.ID.Hex()),
				}
			}

			fieldValues[f.ID.Hex()] = rating

		case FieldTypeDatetime:

			fieldValues[f.ID.Hex()] = nil

			b, err := json.Marshal(newFieldValue)
			if err != nil {
				return nil, err
			}

			var datetime time.Time
			if err := json.Unmarshal(b, &datetime); err != nil {
				return nil, err
			}

			fieldValues[f.ID.Hex()] = datetime

		case FieldTypeBoolean:

			fieldValues[f.ID.Hex()] = nil

			b, err := json.Marshal(newFieldValue)
			if err != nil {
				return nil, err
			}

			var boolean bool
			if err := json.Unmarshal(b, &boolean); err != nil {
				return nil, err
			}

			fieldValues[f.ID.Hex()] = boolean

		case FieldTypeRichText:

			fieldValues[f.ID.Hex()] = nil

			b, err := json.Marshal(newFieldValue)
			if err != nil {
				return nil, err
			}

			var text string
			if err := json.Unmarshal(b, &text); err != nil {
				return nil, err
			}

			fieldValues[f.ID.Hex()] = text

		case FieldTypeText, FieldTypeAddress, FieldTypeURL:

			fieldValues[f.ID.Hex()] = nil

			b, err := json.Marshal(newFieldValue)
			if err != nil {
				return nil, err
			}

			var text string
			if err := json.Unmarshal(b, &text); err != nil {
				return nil, err
			}

			if len(text) > 250 {
				return nil, Error{
					Message:      "max length exceeded",
					Reason:       ErrReasonInvalidArgument,
					LocationType: "argument",
					Location:     fmt.Sprintf("field_values[%v]", f.ID.Hex()),
				}
			}

			fieldValues[f.ID.Hex()] = text

			/*
				case FieldTypeEmail:

					fieldValues[f.ID.Hex()] = nil

					b, err := json.Marshal(newFieldValue)
					if err != nil {
						return nil, err
					}

					emails := make([]Email, 0)
					if err := json.Unmarshal(b, &emails); err != nil {
						return nil, err
					}

					for _, email := range emails {
						if len(email.Value) > 250 {
							return nil, Error{
								Message:      "max length exceeded",
								Reason:       ErrReasonInvalidArgument,
								LocationType: "argument",
								Location:     fmt.Sprintf("field_values[%v]", f.ID.Hex()),
							}
						}
						valid, err := CheckEmail(email.Value)
						if err != nil {
							return nil, err
						}
						if !valid {
							return nil, Error{
								Message:      "email not valid",
								Reason:       ErrReasonInvalidArgument,
								LocationType: "argument",
								Location:     fmt.Sprintf("field_values.%v", f.ID.Hex()),
							}
						}
					}

					fieldValues[f.ID.Hex()] = emails*/

		case FieldTypeUser:

			if newFieldValue == nil {
				fieldValues[f.ID.Hex()] = bson.A{}
			}

			b, err := json.Marshal(newFieldValue)
			if err != nil {
				return nil, err
			}

			var users []primitive.ObjectID
			if err := json.Unmarshal(b, &users); err != nil {
				return nil, err
			}

			// TODO CHECK USERS

			fieldValues[f.ID.Hex()] = users
			/*
				case FieldTypePerson:

					if newFieldValue == nil {
						fieldValues[f.ID.Hex()] = bson.A{}
					}

					b, err := json.Marshal(newFieldValue)
					if err != nil {
						return nil, err
					}

					var persons []primitive.ObjectID
					if err := json.Unmarshal(b, &persons); err != nil {
						return nil, err
					}

					fieldValues[f.ID.Hex()] = persons

				case FieldTypeOrganization:

					if newFieldValue == nil {
						fieldValues[f.ID.Hex()] = bson.A{}
					}

					b, err := json.Marshal(newFieldValue)
					if err != nil {
						return nil, err
					}

					var organizations []primitive.ObjectID
					if err := json.Unmarshal(b, &organizations); err != nil {
						return nil, err
					}

					fieldValues[f.ID.Hex()] = organizations*/

		case FieldTypeSingleSelect:

			fieldValues[f.ID.Hex()] = nil

			b, err := json.Marshal(newFieldValue)
			if err != nil {
				return nil, err
			}

			var choiceID primitive.ObjectID
			if err := json.Unmarshal(b, &choiceID); err != nil {
				return nil, err
			}

			found := false
			for i := 0; i < len(f.Choices) && !found; i++ {
				if f.Choices[i].ID.Hex() == choiceID.Hex() {
					found = true
				}
			}
			if !found {
				return nil, Error{
					Message:      "choice not found",
					Reason:       ErrReasonNotFound,
					LocationType: "argument",
					Location:     fmt.Sprintf("field_values.%v", f.ID.Hex()),
				}
			}

			fieldValues[f.ID.Hex()] = choiceID

		case FieldTypeMultipleSelect:

			if newFieldValue == nil {
				fieldValues[f.ID.Hex()] = bson.A{}
			}

			b, err := json.Marshal(newFieldValue)
			if err != nil {
				return nil, err
			}

			var choiceIDs []primitive.ObjectID
			if err := json.Unmarshal(b, &choiceIDs); err != nil {
				return nil, err
			}

			for _, choiceID := range choiceIDs {
				found := false
				for i := 0; i < len(f.Choices); i++ {
					if choiceID.Hex() == f.Choices[i].ID.Hex() {
						found = true
					}
				}
				if !found {
					return nil, Error{
						Message:      "choice not found",
						Reason:       ErrReasonNotFound,
						LocationType: "argument",
						Location:     fmt.Sprintf("field_values.%v", f.ID.Hex()),
					}
				}
			}

			fieldValues[f.ID.Hex()] = choiceIDs

		default:
			return nil, web.NewShutdownError("invalid field type")
		}
	}

	return fieldValues, nil
}
