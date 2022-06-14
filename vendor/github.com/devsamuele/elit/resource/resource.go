package resource

import (
	"fmt"
	"time"

	"github.com/devsamuele/elit/field"
	"github.com/devsamuele/elit/resperr"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Resource struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	TenantID    primitive.ObjectID `json:"tenant_id" bson:"tenant_id"`
	FieldValues field.Values       `json:"field_values" bson:"field_values"`
	Created     time.Time          `json:"created" bson:"created"`
	Updated     time.Time          `json:"updated" bson:"updated"`
}

type Input struct {
	Values field.Values `json:"field_values"`
}

// func New(name string, fs *field.Store, ss *section.Store, cs *category.Store) {
// 	return Resource{

// 	}
// }

// type Info struct {
// 	ID          primitive.ObjectID `json:"id" bson:"_id"`
// 	TenantID    primitive.ObjectID `json:"tenant_id" bson:"tenant_id"`
// 	Values Values        `json:"field_values" bson:"field_values"`
// 	Created     time.Time          `json:"created" bson:"created"`
// 	Updated     time.Time          `json:"updated" bson:"updated"`
// }

// type Input struct {
// 	Values Values `json:"field_values"`
// }

// type Update struct {
// 	Values Values `json:"field_values"`
// }

// type Response bson.Raw

// func (r Response) Decode(v any) error {
// 	return bson.Unmarshal(r, v)
// }

// type FKInfo struct {
// 	Name            string
// 	OneToManyRel    bool
// 	OwnForeignField bool
// }

// type fkInfo struct {
// 	Name            string
// 	OneToManyRel    bool
// 	OwnForeignField bool
// 	Fields          []Fielder
// }

// type mainInfo struct {
// 	Name   string
// 	Fields []Fielder
// }

// type Resource struct {
// 	db              *mongo.Database
// 	log             *log.Logger
// 	withSection     bool
// 	withCategory    bool
// 	resourceStore   *resourceStore
// 	sectionStore    *sectionStore
// 	fieldStore      *fieldStore
// 	categoryStore   *categoryStore
// 	validFieldTypes []string
// 	name            string
// 	fkResourceMap   map[string]FKResource
// }

// const (
// 	FKRelOneToOne  = "one_to_one"
// 	FKRelOneToMany = "one_to_many"
// )

// type FKResource struct {
// 	Resource
// 	RelOneToMany      bool
// 	LocalForeignField bool
// }

// func New(db *mongo.Database, log *log.Logger, name string, withSection bool, withCategory bool, validFieldTypes []string, fkResourceList []FKResource) Resource {

// 	resourceStore := newResourceStore(db, log, name)
// 	fieldStore := newFieldStore(db, log, name)

// 	fkResourceMap := make(map[string]FKResource)
// 	for _, r := range fkResourceList {
// 		fkResourceMap[r.name] = r
// 	}

// 	b := Resource{
// 		db:              db,
// 		log:             log,
// 		resourceStore:   &resourceStore,
// 		fieldStore:      &fieldStore,
// 		withSection:     withSection,
// 		withCategory:    withCategory,
// 		validFieldTypes: validFieldTypes,
// 		name:            name,
// 		fkResourceMap:   fkResourceMap,
// 	}

// 	if withSection {
// 		sectionStore := newSectionStore(db, log, name)
// 		b.sectionStore = &sectionStore
// 	}

// 	if withCategory {
// 		categoryStore := newCategoryStore(db, log, name)
// 		b.categoryStore = &categoryStore
// 	}

// 	return b
// }

// func (r Resource) StartDBSession() (mongo.Session, error) {
// 	return r.db.Client().StartSession()
// }

// func (r Resource) Name() string {
// 	return r.name
// }

// func (r Resource) FKResourceNames() []string {
// 	return maps.Keys(r.fkResourceMap)
// }

// func (r Resource) FKResourceList() []Resource {
// 	var resourceList []Resource
// 	for _, r := range r.fkResourceMap {
// 		resourceList = append(resourceList, r.Resource)
// 	}
// 	return resourceList
// }

func (i Input) Build(fields []field.Fielder) (field.Values, error) {

	fieldValues := make(field.Values)

	for _, f := range fields {
		fieldValue, ok := i.Values[f.GetBase().ID.Hex()]
		if !ok {
			if f.GetBase().Required {
				return nil, resperr.Error{
					Message:      fmt.Sprintf("field[%v] is required", f.GetBase().ID.Hex()),
					Reason:       resperr.ErrReasonRequired,
					LocationType: "argument",
					Location:     fmt.Sprintf("field[%v]", f.GetBase().ID.Hex()),
				}
			}
			fieldValues[f.GetBase().ID.Hex()] = f.GetDefaultValue()
		} else {
			decodedValue, err := f.DecodeAndValidateValue(fieldValue)
			if err != nil {
				return nil, err
			}
			fieldValues[f.GetBase().ID.Hex()] = decodedValue
		}
	}
	return fieldValues, nil
}

func (i Input) Update(resource *Resource, fields []field.Fielder) error {

	for _, f := range fields {
		fieldValue, ok := i.Values[f.GetBase().ID.Hex()]
		if ok {
			decodedValue, err := f.DecodeAndValidateValue(fieldValue)
			if err != nil {
				return err
			}
			resource.FieldValues[f.GetBase().ID.Hex()] = decodedValue
		}
	}
	return nil
}

// func (b Builder) MainInfo(ctx context.Context, tenantID primitive.ObjectID) (mainInfo, error) {
// 	fields, err := b.fieldStore.Query(ctx, tenantID)
// 	if err != nil {
// 		return mainInfo{}, err
// 	}
// 	return mainInfo{
// 		Name:   b.mainResource,
// 		Fields: fields,
// 	}, nil
// }

// func (b Builder) FKInfo(ctx context.Context, tenantID primitive.ObjectID) ([]fkInfo, error) {
// 	var resources []fkInfo
// 	for _, r := range b.fkResources {
// 		fields, err := newFieldStore(b.db, b.log, r.Name).Query(ctx, tenantID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		resources = append(resources, fkInfo{
// 			Name:            r.Name,
// 			Fields:          fields,
// 			OneToManyRel:    r.OneToManyRel,
// 			OwnForeignField: r.OwnForeignField,
// 		})
// 	}
// 	return resources, nil
// }

// func (r Resource) FieldStore() *fieldStore {
// 	return r.fieldStore
// }
// func (r Resource) SectionStore() *sectionStore {
// 	if r.sectionStore == nil {
// 		r.log.Fatalln("resource not configured with section")
// 	}
// 	return r.sectionStore
// }
// func (r Resource) CategoryStore() *categoryStore {
// 	if r.categoryStore == nil {
// 		r.log.Fatalln("resource not configured with category")
// 	}
// 	return r.categoryStore
// }
// func (r Resource) ResourceStore() *resourceStore {
// 	return r.resourceStore
// }

// func (r Resource) SetUpIndex(ctx context.Context, fkIDs []primitive.ObjectID) error {

// 	im := []mongo.IndexModel{
// 		{Keys: bson.D{{Key: "tenant_id", Value: 1}},
// 			Options: nil},
// 	}

// 	for _, fkID := range fkIDs {
// 		im = append(im, mongo.IndexModel{
// 			Keys: bson.D{{Key: fmt.Sprintf("field_values.%v", fkID.Hex()), Value: 1}},
// 		})
// 	}

// 	_, err := r.db.Collection(r.name).Indexes().CreateMany(ctx, im)
// 	if err != nil {
// 		return fmt.Errorf("creating %v index: %w", r.name, err)
// 	}

// 	_, err = r.db.Collection(fmt.Sprintf("%v_field", r.name)).Indexes().CreateOne(ctx, mongo.IndexModel{
// 		Keys:    bson.D{{Key: "tenant_id", Value: 1}},
// 		Options: nil,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("creating %v_field index: %w", r.name, err)
// 	}

// 	if r.withSection {
// 		_, err = r.db.Collection(fmt.Sprintf("%v_section", r.name)).Indexes().CreateOne(ctx, mongo.IndexModel{
// 			Keys:    bson.D{{Key: "tenant_id", Value: 1}},
// 			Options: nil,
// 		})
// 		if err != nil {
// 			return fmt.Errorf("creating %v_section index: %w", r.name, err)
// 		}
// 	}

// 	if r.withCategory {
// 		_, err = r.db.Collection(fmt.Sprintf("%v_category", r.name)).Indexes().CreateOne(ctx, mongo.IndexModel{
// 			Keys:    bson.D{{Key: "tenant_id", Value: 1}},
// 			Options: nil,
// 		})
// 		if err != nil {
// 			return fmt.Errorf("creating %v_category index: %w", r.name, err)
// 		}
// 	}

// 	return nil

// }
