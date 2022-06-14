package filter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/devsamuele/elit/field"
	"github.com/devsamuele/elit/resperr"
	"github.com/devsamuele/elit/utility/slice"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// var ErrInvalidMatchType = errors.New("invalid match_type")
// var ErrInvalidResource = errors.New("invalid resource")
// var ErrInvalidValue = errors.New("invalid value")
// var ErrFieldNotFound = errors.New("field not found")

// var ErrInvalidCondition = errors.New("invalid condition")

var (
	ErrInvalidSortField    = errors.New("invalid sort field")
	ErrInvalidSortOrder    = errors.New("invalid sort order")
	ErrSortFieldNotFound   = errors.New("sort field not found")
	ErrInvalidPageNumber   = errors.New("invalid page number")
	ErrInvalidItemsPerPage = errors.New("invalid items per page")
)

const typeForeignKey = "foreign_key"

const (
	MatchTypeEqual            = "eq"
	MatchTypeNotEqual         = "ne"
	MatchTypeContains         = "contains"
	MatchTypeNotContains      = "not_contains"
	MatchTypeStartWith        = "start_with"
	MatchTypeEndWith          = "end_with"
	MatchTypeEmpty            = "empty"
	MatchTypeNotEmpty         = "not_empty"
	MatchTypeGreaterThan      = "gt"
	MatchTypeGreaterThanEqual = "gte"
	MatchTypeLessThan         = "lt"
	MatchTypeLessThanEqual    = "lte"
	MatchTypeAnyOf            = "any_of"
	MatchTypeNoneOf           = "none_of"
)

type Sort struct {
	FieldID string `json:"field_id" bson:"field_id"`
	Order   int    `json:"order" bson:"order"`
}

type Filter struct {
	Sort    Sort       `json:"sort"`
	ORGroup []ANDGroup `json:"or_group"`
}

type Input struct {
	PageNumber   int           `json:"page_number"`
	ItemsPerPage int           `json:"items_per_page"`
	Sort         Sort          `json:"sort"`
	ORGroup      []NewANDGroup `json:"or_group"`
}

type ANDGroup struct {
	Condition []Condition `json:"and_group" bson:"and_group"`
}

type NewANDGroup struct {
	Condition []NewCondition `json:"and_group"`
}

type Condition struct {
	Resource  string      `json:"resource" bson:"resource"`
	FieldID   string      `json:"field_id" bson:"field_id"`
	Value     interface{} `json:"value" bson:"value"`
	MatchType string      `json:"match_type" bson:"matchType"`
}

type NewCondition struct {
	Resource  string      `json:"resource"`
	FieldID   string      `json:"field_id"`
	Value     interface{} `json:"value"`
	MatchType string      `json:"match_type"`
}

// type Resource struct {
// 	Name   string
// 	Fields []field.Fielder
// }

// type Builder struct {
// 	mainResource Resource
// 	fkResources  []Resource
// }

// func (b Builder) FkResourceNames() []string {
// 	return slice.Map(b.fkResources, func(r Resource) string {
// 		return r.Name
// 	})
// }

// func (b Builder) FkResources() []Resource {
// 	return b.fkResources
// }

// func (b Builder) MainResource() Resource {
// 	return b.mainResource
// }

// func NewBuilder(mainResource Resource, fkResources []Resource) Builder {
// 	b := Builder{
// 		mainResource: mainResource,
// 		fkResources:  fkResources,
// 	}
// 	return b
// }

func validatePageNumber(pageNumber int) (int, error) {
	var defaultPageNumber = 1 //* ascending order A-Z, 1..9
	if pageNumber == 0 {
		return defaultPageNumber, nil
	}

	if pageNumber < 0 {
		return 0, ErrInvalidPageNumber
	}

	return pageNumber, nil
}

func validateItemsPerPage(itemsPerPage int) (int, error) {
	var defaultItemsPerPage = 50
	if itemsPerPage == 0 {
		return defaultItemsPerPage, nil
	}

	if itemsPerPage < 0 || itemsPerPage > 1000 {
		return 0, ErrInvalidItemsPerPage
	}

	return itemsPerPage, nil
}

func validateSort(sort Sort) (Sort, error) {
	if sort.FieldID == "" {
		return Sort{}, ErrInvalidSortField
	}

	order := 1 //* default: ascending order A-Z, 1..9
	if sort.Order != -1 && sort.Order != 1 {
		return Sort{}, ErrInvalidSortOrder
	}
	order = sort.Order

	return Sort{
		FieldID: sort.FieldID,
		Order:   order,
	}, nil
}

func (i *Input) NewFromInput(ctx context.Context, tenantID primitive.ObjectID, MainResource string, fieldResourceMap map[string][]field.Fielder) (bson.D, error) {

	sort, err := validateSort(i.Sort)
	if err != nil {
		return nil, err
	}
	i.Sort = sort

	itemsPerPage, err := validateItemsPerPage(i.ItemsPerPage)
	if err != nil {
		return nil, err
	}
	i.ItemsPerPage = itemsPerPage

	pageNumber, err := validatePageNumber(i.PageNumber)
	if err != nil {
		return nil, err
	}
	i.PageNumber = pageNumber

	mainResourceFields := fieldResourceMap[MainResource]
	_, found := slice.Find(mainResourceFields, func(f field.Fielder) bool {
		return f.GetBase().ID.Hex() == i.Sort.FieldID
	})
	if !found {
		return nil, ErrSortFieldNotFound
	}

	if sort.FieldID != "created" && sort.FieldID != "updated" && sort.FieldID != "_id" {
		sort.FieldID = fmt.Sprintf("field_values.%v", sort.FieldID)
	}

	var ORGroup bson.A
	for _, andGroup := range i.ORGroup {
		var ANDGroup bson.A
		for _, condition := range andGroup.Condition {

			var fieldPath string

			// *Check resource and fieldID
			currentResource := condition.Resource
			var currentFieldType string
			var currentFieldSearchable bool
			var currentFieldID string

			fields, foundResource := fieldResourceMap[currentResource]
			if !foundResource {
				return nil, resperr.Error{Message: "resource not found"}
			}

			f, foundField := slice.Find(fields, func(f field.Fielder) bool {
				return f.GetBase().ID.Hex() == condition.FieldID
			})

			if foundField {
				currentFieldType = f.GetBase().Type
				currentFieldID = f.GetBase().ID.Hex()
				currentFieldSearchable = f.GetBase().Searchable

				fieldPath = fmt.Sprintf("field_values.%v", currentFieldID)
				if MainResource != currentResource {
					fieldPath = fmt.Sprintf("%v.%v", currentResource, fieldPath)
				}

			} else {
				if condition.FieldID == "created" || condition.FieldID == "updated" {
					currentFieldType = field.TypeDatetime
					currentFieldID = condition.FieldID
					currentFieldSearchable = true
				} else if condition.FieldID == "_id" {
					currentFieldType = typeForeignKey
					currentFieldID = condition.FieldID
					currentFieldSearchable = true
				} else {
					return nil, resperr.Error{}
				}

				fieldPath = currentFieldID
				if MainResource != currentResource {
					fieldPath = fmt.Sprintf("%v.%v", currentResource, fieldPath)
				}
			}

			if !currentFieldSearchable {
				return nil, resperr.Error{}
			}

			builtCondition, err := buildBsonCondition(currentFieldType, fieldPath, condition.MatchType, condition.Value)
			if err != nil {
				return nil, err
			}
			ANDGroup = append(ANDGroup, builtCondition)
		}
		ORGroup = append(ORGroup, bson.D{{Key: "$and", Value: ANDGroup}})
	}

	result := bson.D{{Key: "$or", Value: ORGroup}}
	return result, nil
}

func buildBsonCondition(currentFieldType, fieldPath, currentMatchType string, currentValue interface{}) (bson.D, error) {
	conditionValueB, err := json.Marshal(currentValue)
	if err != nil {
		return nil, err
	}
	switch currentFieldType {

	// *String
	case field.TypeText, field.TypeAddress, field.TypeRichText, field.TypeURL, field.TypePhone:
		var text string
		if err := json.Unmarshal(conditionValueB, &text); err != nil {
			return nil, err
		}

		switch currentMatchType {
		case MatchTypeEqual:
			return bson.D{{Key: fieldPath, Value: bson.D{
				{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("^%v$", text), Options: "i"}}}}}, nil

		case MatchTypeNotEqual:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$not", Value: bson.D{
				{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("^%v$", text), Options: "i"}}}}}}}, nil

		case MatchTypeContains:
			return bson.D{{Key: fieldPath, Value: bson.D{
				{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("%v", text), Options: "i"}}}}}, nil

		case MatchTypeNotContains:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$not", Value: bson.D{
				{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("%v", text), Options: "i"}}}}}}}, nil

		case MatchTypeStartWith:
			return bson.D{{Key: fieldPath, Value: bson.D{
				{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("^%v", text), Options: "i"}}}}}, nil

		case MatchTypeEndWith:
			return bson.D{{Key: fieldPath, Value: bson.D{
				{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("%v$", text), Options: "i"}}}}}, nil

		case MatchTypeEmpty:
			return bson.D{{Key: fieldPath, Value: nil}}, nil

		case MatchTypeNotEmpty:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$ne", Value: nil}}}}, nil

		default:
			return nil, resperr.Error{}
		}

	// *Email
	case field.TypeEmail:
		var text string
		if err := json.Unmarshal(conditionValueB, &text); err != nil {
			return nil, err
		}
		switch currentMatchType {
		case MatchTypeEqual:
			return bson.D{{Key: "$or", Value: bson.A{
				bson.D{{Key: fmt.Sprintf("%v.primary", fieldPath), Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("^%v$", text), Options: "i"}}}}},
				bson.D{{Key: fmt.Sprintf("%v.secondary", fieldPath), Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("^%v$", text), Options: "i"}}}}},
			}}}, nil

		case MatchTypeNotEqual:
			return bson.D{{Key: "$or", Value: bson.A{
				bson.D{{Key: fmt.Sprintf("%v.primary", fieldPath), Value: bson.D{{Key: "$not", Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("^%v$", text), Options: "i"}}}}}}},
				bson.D{{Key: fmt.Sprintf("%v.secondary", fieldPath), Value: bson.D{{Key: "$not", Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("^%v$", text), Options: "i"}}}}}}},
			}}}, nil

		case MatchTypeContains:
			return bson.D{{Key: "$or", Value: bson.A{
				bson.D{{Key: fmt.Sprintf("%v.primary", fieldPath), Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("%v", text), Options: "i"}}}}},
				bson.D{{Key: fmt.Sprintf("%v.secondary", fieldPath), Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("%v", text), Options: "i"}}}}},
			}}}, nil

		case MatchTypeNotContains:
			return bson.D{{Key: "$or", Value: bson.A{
				bson.D{{Key: fmt.Sprintf("%v.primary", fieldPath), Value: bson.D{{Key: "$not", Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("%v", text), Options: "i"}}}}}}},
				bson.D{{Key: fmt.Sprintf("%v.secondary", fieldPath), Value: bson.D{{Key: "$not", Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("%v", text), Options: "i"}}}}}}},
			}}}, nil

		case MatchTypeStartWith:
			return bson.D{{Key: "$or", Value: bson.A{
				bson.D{{Key: fmt.Sprintf("%v.primary", fieldPath), Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("^%v", text), Options: "i"}}}}},
				bson.D{{Key: fmt.Sprintf("%v.secondary", fieldPath), Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("^%v", text), Options: "i"}}}}},
			}}}, nil

		case MatchTypeEndWith:
			return bson.D{{Key: "$or", Value: bson.A{
				bson.D{{Key: fmt.Sprintf("%v.primary", fieldPath), Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("%v$", text), Options: "i"}}}}},
				bson.D{{Key: fmt.Sprintf("%v.secondary", fieldPath), Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: fmt.Sprintf("%v$", text), Options: "i"}}}}},
			}}}, nil

		case MatchTypeEmpty:
			return bson.D{{Key: "$or", Value: bson.A{
				bson.D{{Key: fmt.Sprintf("%v.primary", fieldPath), Value: nil}},
				bson.D{{Key: fmt.Sprintf("%v.secondary", fieldPath), Value: nil}},
			}}}, nil

		case MatchTypeNotEmpty:
			return bson.D{{Key: "$or", Value: bson.A{
				bson.D{{Key: fmt.Sprintf("%v.primary", fieldPath), Value: bson.D{{Key: "$ne", Value: nil}}}},
				bson.D{{Key: fmt.Sprintf("%v.secondary", fieldPath), Value: bson.D{{Key: "$ne", Value: nil}}}},
			}}}, nil

		default:
			return nil, resperr.Error{}
		}

	// *Boolean
	case field.TypeBoolean:
		var boolean bool
		if err := json.Unmarshal(conditionValueB, &boolean); err != nil {
			return nil, err
		}
		switch currentMatchType {
		case MatchTypeEqual:
			return bson.D{{Key: fieldPath, Value: boolean}}, nil

		case MatchTypeNotEqual:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$ne", Value: boolean}}}}, nil

		case MatchTypeEmpty:
			return bson.D{{Key: fieldPath, Value: nil}}, nil

		case MatchTypeNotEmpty:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$ne", Value: nil}}}}, nil

		default:
			return nil, resperr.Error{}
		}

	// *Number
	case field.TypeNumber, field.TypeCurrency, field.TypeDuration, field.TypePercent, field.TypeRating:
		var num float64
		if err := json.Unmarshal(conditionValueB, &num); err != nil {
			return nil, err
		}
		switch currentMatchType {
		case MatchTypeEqual:
			return bson.D{{Key: fieldPath, Value: num}}, nil

		case MatchTypeNotEqual:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$ne", Value: num}}}}, nil

		case MatchTypeEmpty:
			return bson.D{{Key: fieldPath, Value: nil}}, nil

		case MatchTypeNotEmpty:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$ne", Value: nil}}}}, nil

		case MatchTypeGreaterThan:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$gt", Value: num}}}}, nil

		case MatchTypeGreaterThanEqual:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$gte", Value: num}}}}, nil

		case MatchTypeLessThan:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$lt", Value: num}}}}, nil

		case MatchTypeLessThanEqual:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$lte", Value: num}}}}, nil

		default:
			return nil, resperr.Error{}
		}

	// *Datetime
	case field.TypeDatetime:
		var dt time.Time
		if err := json.Unmarshal(conditionValueB, &dt); err != nil {
			return nil, err
		}
		switch currentMatchType {
		case MatchTypeEqual:
			return bson.D{{Key: fieldPath, Value: dt}}, nil

		case MatchTypeNotEqual:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$ne", Value: dt}}}}, nil

		case MatchTypeEmpty:
			return bson.D{{Key: fieldPath, Value: nil}}, nil

		case MatchTypeNotEmpty:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$ne", Value: nil}}}}, nil

		case MatchTypeGreaterThan:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$gt", Value: dt}}}}, nil

		case MatchTypeGreaterThanEqual:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$gte", Value: dt}}}}, nil

		case MatchTypeLessThan:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$lt", Value: dt}}}}, nil

		case MatchTypeLessThanEqual:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$lte", Value: dt}}}}, nil

		default:
			return nil, resperr.Error{}
		}

	// *
	case field.TypeDatetimeRange:
		// datetime

	// * ForeignKey
	case typeForeignKey:
		var fk primitive.ObjectID
		if err := json.Unmarshal(conditionValueB, &fk); err != nil {
			return nil, err
		}

		switch currentMatchType {
		case MatchTypeEqual:
			return bson.D{{Key: fieldPath, Value: fk}}, nil

		case MatchTypeNotEqual:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$ne", Value: fk}}}}, nil

		case MatchTypeEmpty:
			return bson.D{{Key: fieldPath, Value: nil}}, nil

		case MatchTypeNotEmpty:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$ne", Value: nil}}}}, nil

		default:
			return nil, resperr.Error{}
		}

	// *[]ObjectID
	case field.TypeSingleSelect, field.TypeMultipleSelect, field.TypeUser:
		var IDs []primitive.ObjectID
		if err := json.Unmarshal(conditionValueB, &IDs); err != nil {
			return nil, err
		}
		switch currentMatchType {
		case MatchTypeAnyOf:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$in", Value: IDs}}}}, nil

		case MatchTypeNoneOf:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$nin", Value: IDs}}}}, nil

		case MatchTypeEmpty:
			return bson.D{{Key: fieldPath, Value: bson.A{}}}, nil

		case MatchTypeNotEmpty:
			return bson.D{{Key: fieldPath, Value: bson.D{{Key: "$ne", Value: bson.A{}}}}}, nil

		default:
			return nil, resperr.Error{}
		}
	}

	return nil, errors.New("invalid field type")
}
