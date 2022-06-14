package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidMatchType = errors.New("invalid match_type")

const (
	FilterMatchTypeEqual            = "eq"
	FilterMatchTypeNotEqual         = "ne"
	FilterMatchTypeContains         = "contains"
	FilterMatchTypeNotContains      = "not_contains"
	FilterMatchTypeStartWith        = "start_with"
	FilterMatchTypeEndWith          = "end_with"
	FilterMatchTypeEmpty            = "empty"
	FilterMatchTypeNotEmpty         = "not_empty"
	FilterMatchTypeGreaterThan      = "gt"
	FilterMatchTypeGreaterThanEqual = "gte"
	FilterMatchTypeLessThan         = "lt"
	FilterMatchTypeLessThanEqual    = "lte"
	FilterMatchTypeAnyOf            = "any_of"
	FilterMatchTypeNoneOf           = "none_of"
)

const (
	SearchFieldTypeString       = "string"
	SearchFieldTypeStringList   = "string_list"
	SearchFieldTypeNumber       = "number"
	SearchFieldTypeDatetime     = "datetime"
	SearchFieldTypeBoolean      = "boolean"
	SearchFieldTypeObjectID     = "objectID"
	SearchFieldTypeObjectIDList = "objectID_list"
)

func ToSearchFieldType(fieldType string) string {
	switch fieldType {
	case FieldTypeText, FieldTypeAddress, FieldTypeRichText, FieldTypeURL:
		return SearchFieldTypeString
	case FieldTypeNumber, FieldTypeCurrency, FieldTypeDuration, FieldTypePercent, FieldTypeRating:
		return SearchFieldTypeNumber
	case FieldTypeBoolean:
		return SearchFieldTypeBoolean
	case FieldTypeDatetime:
		return SearchFieldTypeDatetime
	case FieldTypeUser:
		return SearchFieldTypeObjectIDList
	case FieldTypeSingleSelect, FieldTypeMultipleSelect:
		return SearchFieldTypeStringList
	default:
		return ""
	}
}

type Filter struct {
	Sort     Sort             `json:"sort"`
	ORFields []ANDFilterField `json:"or_fields"`
}

type ANDFilterField struct {
	ANDFields []FilterField `json:"and_fields" bson:"ANDFields"`
}

type FilterField struct {
	Resource  string      `json:"resource" bson:"resource"`
	FieldPath string      `json:"field_path" bson:"field_path"`
	FieldType string      `json:"field_type" bson:"fieldType"`
	Value     interface{} `json:"value" bson:"value"`
	MatchType string      `json:"match_type" bson:"matchType"`
}

type Sort struct {
	Field string `json:"field" bson:"field"`
	Order int    `json:"order" bson:"order"`
}

type NewFilter struct {
	PageNumber   int                 `json:"page_number"`
	ItemsPerPage int                 `json:"items_per_page"`
	Sort         Sort                `json:"sort"`
	ORFields     []NewANDFilterField `json:"or_fields"`
}

type NewANDFilterField struct {
	ANDFields []NewFilterField `json:"and_fields"`
}

type NewFilterField struct {
	Resource  string      `json:"resource"`
	FieldPath string      `json:"field_path"`
	Value     interface{} `json:"value"`
	MatchType string      `json:"match_type"`
}

func (filterField FilterField) BuildBSONFilter(currentResource string) (bson.D, error) {
	b, err := json.Marshal(filterField.Value)
	if err != nil {
		return nil, err
	}

	if filterField.Resource != currentResource {
		filterField.FieldPath = fmt.Sprintf("%v.%v", filterField.Resource, filterField.FieldPath)
	}

	switch filterField.FieldType {

	case SearchFieldTypeString:
		var text string
		if err := json.Unmarshal(b, &text); err != nil {
			return nil, err
		}
		switch filterField.MatchType {
		case FilterMatchTypeEqual:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{
				{"$regex", primitive.Regex{Pattern: fmt.Sprintf("^%v$", text), Options: "i"}}}}}, nil

		case FilterMatchTypeNotEqual:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$not", Value: bson.D{
				{"$regex", primitive.Regex{Pattern: fmt.Sprintf("^%v$", text), Options: "i"}}}}}}}, nil

		case FilterMatchTypeContains:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{
				{"$regex", primitive.Regex{Pattern: fmt.Sprintf("%v", text), Options: "i"}}}}}, nil

		case FilterMatchTypeNotContains:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$not", Value: bson.D{
				{"$regex", primitive.Regex{Pattern: fmt.Sprintf("%v", text), Options: "i"}}}}}}}, nil

		case FilterMatchTypeStartWith:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{
				{"$regex", primitive.Regex{Pattern: fmt.Sprintf("^%v", text), Options: "i"}}}}}, nil

		case FilterMatchTypeEndWith:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{
				{"$regex", primitive.Regex{Pattern: fmt.Sprintf("%v$", text), Options: "i"}}}}}, nil

		case FilterMatchTypeEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: nil}}, nil

		case FilterMatchTypeNotEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$ne", Value: nil}}}}, nil
		}

	case SearchFieldTypeNumber:
		var num float64
		if err := json.Unmarshal(b, &num); err != nil {
			return nil, err
		}
		switch filterField.MatchType {
		case FilterMatchTypeEqual:
			return bson.D{{Key: filterField.FieldPath, Value: num}}, nil

		case FilterMatchTypeNotEqual:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$ne", Value: num}}}}, nil

		case FilterMatchTypeEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: nil}}, nil

		case FilterMatchTypeNotEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$ne", Value: nil}}}}, nil

		case FilterMatchTypeGreaterThan:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$gt", Value: num}}}}, nil

		case FilterMatchTypeGreaterThanEqual:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$gte", Value: num}}}}, nil

		case FilterMatchTypeLessThan:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$lt", Value: num}}}}, nil

		case FilterMatchTypeLessThanEqual:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$lte", Value: num}}}}, nil
		}

	case SearchFieldTypeBoolean:
		var boolean bool
		if err := json.Unmarshal(b, &boolean); err != nil {
			return nil, err
		}
		switch filterField.MatchType {
		case FilterMatchTypeEqual:
			return bson.D{{Key: filterField.FieldPath, Value: boolean}}, nil

		case FilterMatchTypeNotEqual:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$ne", Value: boolean}}}}, nil

		case FilterMatchTypeEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: nil}}, nil

		case FilterMatchTypeNotEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$ne", Value: nil}}}}, nil
		}

	case SearchFieldTypeDatetime:
		var dt time.Time
		if err := json.Unmarshal(b, &dt); err != nil {
			return nil, err
		}
		switch filterField.MatchType {
		case FilterMatchTypeEqual:
			return bson.D{{Key: filterField.FieldPath, Value: dt}}, nil

		case FilterMatchTypeNotEqual:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$ne", Value: dt}}}}, nil

		case FilterMatchTypeEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: nil}}, nil

		case FilterMatchTypeNotEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$ne", Value: nil}}}}, nil

		case FilterMatchTypeGreaterThan:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$gt", Value: dt}}}}, nil

		case FilterMatchTypeGreaterThanEqual:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$gte", Value: dt}}}}, nil

		case FilterMatchTypeLessThan:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$lt", Value: dt}}}}, nil

		case FilterMatchTypeLessThanEqual:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$lte", Value: dt}}}}, nil
		}

	case SearchFieldTypeObjectIDList:
		var IDs []primitive.ObjectID
		if err := json.Unmarshal(b, &IDs); err != nil {
			return nil, err
		}
		switch filterField.MatchType {
		case FilterMatchTypeAnyOf:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$in", Value: IDs}}}}, nil

		case FilterMatchTypeNoneOf:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$nin", Value: IDs}}}}, nil

		case FilterMatchTypeEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: bson.A{}}}, nil

		case FilterMatchTypeNotEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$ne", Value: bson.A{}}}}}, nil
		}

	case SearchFieldTypeObjectID:
		var ID primitive.ObjectID
		if err := json.Unmarshal(b, &ID); err != nil {
			return nil, err
		}
		switch filterField.MatchType {
		case FilterMatchTypeEqual:
			return bson.D{{Key: filterField.FieldPath, Value: ID}}, nil

		case FilterMatchTypeNotEqual:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$ne", Value: ID}}}}, nil

		case FilterMatchTypeEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: nil}}, nil

		case FilterMatchTypeNotEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$ne", Value: nil}}}}, nil
		}

	case SearchFieldTypeStringList:
		var stringList []string
		if err := json.Unmarshal(b, &stringList); err != nil {
			return nil, err
		}
		switch filterField.MatchType {
		case FilterMatchTypeAnyOf:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$in", Value: stringList}}}}, nil

		case FilterMatchTypeNoneOf:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$nin", Value: stringList}}}}, nil

		case FilterMatchTypeEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: bson.A{}}}, nil

		case FilterMatchTypeNotEmpty:
			return bson.D{{Key: filterField.FieldPath, Value: bson.D{{Key: "$ne", Value: bson.A{}}}}}, nil
		}
	}

	return nil, ErrInvalidMatchType
}

func ValidMatchType(searchFieldType, matchType string) bool {
	switch searchFieldType {
	case SearchFieldTypeString:
		return matchType == FilterMatchTypeEqual ||
			matchType == FilterMatchTypeNotEqual ||
			matchType == FilterMatchTypeContains ||
			matchType == FilterMatchTypeNotContains ||
			matchType == FilterMatchTypeEmpty ||
			matchType == FilterMatchTypeNotEmpty ||
			matchType == FilterMatchTypeStartWith ||
			matchType == FilterMatchTypeEndWith

	case SearchFieldTypeNumber, SearchFieldTypeDatetime:
		return matchType == FilterMatchTypeEqual ||
			matchType == FilterMatchTypeNotEqual ||
			matchType == FilterMatchTypeEmpty ||
			matchType == FilterMatchTypeNotEmpty ||
			matchType == FilterMatchTypeGreaterThan ||
			matchType == FilterMatchTypeGreaterThanEqual ||
			matchType == FilterMatchTypeLessThan ||
			matchType == FilterMatchTypeLessThanEqual ||
			matchType == FilterMatchTypeEndWith

	case SearchFieldTypeBoolean, SearchFieldTypeObjectID:
		return matchType == FilterMatchTypeEqual ||
			matchType == FilterMatchTypeNotEqual ||
			matchType == FilterMatchTypeEmpty ||
			matchType == FilterMatchTypeNotEmpty

	case SearchFieldTypeObjectIDList, SearchFieldTypeStringList:
		return matchType == FilterMatchTypeAnyOf ||
			matchType == FilterMatchTypeNoneOf ||
			matchType == FilterMatchTypeEmpty ||
			matchType == FilterMatchTypeNotEmpty

	default:
		return false
	}
}
