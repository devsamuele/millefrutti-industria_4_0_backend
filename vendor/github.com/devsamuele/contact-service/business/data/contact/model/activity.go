package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Activity struct {
	ID              primitive.ObjectID   `json:"id" bson:"_id"`
	TenantID        primitive.ObjectID   `json:"tenant_id" bson:"tenant_id"`
	OrganizationIDs []primitive.ObjectID `json:"organization_ids" bson:"organization_ids"`
	PersonIDs       []primitive.ObjectID `json:"person_ids" bson:"person_ids"`
	//DealID          primitive.ObjectID     `json:"deal_id" bson:"deal_id"`
	//ProjectID       primitive.ObjectID     `json:"project_id" bson:"project_id"`
	CategoryID   primitive.ObjectID     `json:"category_id" bson:"category_id"`
	UserIDs      []primitive.ObjectID   `json:"user_ids" bson:"user_ids"`
	Title        string                 `json:"title" bson:"title"`
	Date         DatetimeRange          `json:"date" bson:"date"`
	Event        *Event                 `json:"event,omitempty" bson:"event,omitempty"`
	FieldValues  map[string]interface{} `json:"field_values" bson:"field_values"`
	Private      bool                   `json:"private" bson:"private"`
	Done         bool                   `json:"done" bson:"done"`
	Deleted      bool                   `json:"deleted,omitempty" bson:"deleted"`
	DateOfDelete *time.Time             `json:"date_of_delete,omitempty" bson:"date_of_delete"`
	Created      time.Time              `json:"created" bson:"created"`
	Updated      time.Time              `json:"updated" bson:"updated"`
}

type ActivityResponse struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	TenantID      primitive.ObjectID `json:"tenant_id" bson:"tenant_id"`
	Organizations []Organization     `json:"organizations" bson:"organizations"`
	Persons       []Person           `json:"persons" bson:"persons"`
	//Deal          primitive.ObjectID     `json:"deal" bson:"deal"`
	//Project       primitive.ObjectID     `json:"project" bson:"project"`
	Category     Category               `json:"category" bson:"category"`
	UserIDs      []primitive.ObjectID   `json:"user_ids" bson:"user_ids"`
	Title        string                 `json:"title" bson:"title"`
	Date         DatetimeRange          `json:"date" bson:"date"`
	Event        *Event                 `json:"event,omitempty" bson:"event,omitempty"`
	FieldValues  map[string]interface{} `json:"field_values" bson:"field_values"`
	Private      bool                   `json:"private" bson:"private"`
	Done         bool                   `json:"done" bson:"done"`
	Deleted      bool                   `json:"deleted,omitempty" bson:"deleted"`
	DateOfDelete *time.Time             `json:"date_of_delete,omitempty" bson:"date_of_delete"`
	Created      time.Time              `json:"created" bson:"created"`
	Updated      time.Time              `json:"updated" bson:"updated"`
}

type NewActivity struct {
	Title      string              `json:"title" bson:"title"` // required -> not ""
	Date       NewDatetimeRange    `json:"date" bson:"date"`   // Date.From required
	CategoryID *primitive.ObjectID `json:"category_id" bson:"category_id"`
	Private    bool                `json:"private" bson:"private"`
}

type DatetimeRange struct {
	Start      *time.Time `json:"start" bson:"start"`
	End        *time.Time `json:"end" bson:"end"`
	Expiration *time.Time `json:"expiration" bson:"expiration"`
	AllDay     bool       `json:"all_day" bson:"all_day"`
}

type NewDatetimeRange struct {
	Start      *time.Time `json:"start" bson:"start"`
	End        *time.Time `json:"end" bson:"end"`
	Expiration *time.Time `json:"expiration" bson:"expiration"`
	AllDay     *bool      `json:"all_day" bson:"all_day"`
}

type Category struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	TenantID primitive.ObjectID `json:"tenant_id" bson:"tenant_id"`
	Name     Translation        `json:"name" bson:"name"`
	Created  time.Time          `json:"created" bson:"created"`
	Updated  time.Time          `json:"updated" bson:"updated"`
}

type NewCategory struct {
	Name Translation `json:"name"`
}

// EVENT

type Event struct {
	Link        *string  `json:"link" bson:"link"`
	Position    *string  `json:"position" bson:"position"`
	Hosts       []string `json:"hosts" bson:"hosts"`
	Description *string  `json:"description" bson:"description"`
	Sent        *bool    `json:"sent" bson:"sent"`
}

// CHECKLIST
// max 100 checklistItem

type Checklist struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	TenantID   primitive.ObjectID `json:"tenant_id" bson:"tenant_id"`
	Resource   string             `json:"resource" bson:"resource"`
	ResourceID primitive.ObjectID `json:"resource_id" bson:"resource_id"`
	Items      []ChecklistItem    `json:"items" bson:"items"`
	Title      Translation        `json:"title" bson:"title"`
	Created    time.Time          `json:"created" bson:"created"`
	Updated    time.Time          `json:"updated" bson:"updated"`
}

type NewChecklistItem struct {
	Text *string `json:"text"` // required - not ""
}

type UpdateChecklistItem struct {
	Text     *string             `json:"text" json:"text"` // set limit of length -> max 250 char ?
	Done     *bool               `json:"done" bson:"done"`
	Datetime *time.Time          `json:"datetime" bson:"datetime"`
	UserID   *primitive.ObjectID `json:"user_id" bson:"user_id"`
}

type ChecklistItem struct {
	Text     string              `json:"text" bson:"text"`
	Done     bool                `json:"done" bson:"done"`
	Datetime *time.Time          `json:"datetime" bson:"datetime"`
	UserID   *primitive.ObjectID `json:"user_id" bson:"user_id"`
}

// REMINDER

type Reminder struct {
	Type     string               `json:"type" bson:"type"`         // notification, email ...
	UserIDs  []primitive.ObjectID `json:"user_ids" bson:"user_ids"` // who notify -> get email ??
	Email    []string             `json:"email" bson:"email"`       // email to notify ??
	Datetime time.Time            `json:"datetime" bson:"datetime"`
	Created  time.Time            `json:"created" bson:"created"`
	Updated  time.Time            `json:"updated" bson:"updated"`
}

// NOTE

type Note struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	TenantID primitive.ObjectID `json:"tenant_id" bson:"tenant_id"`
	Text     string             `json:"text" bson:"text"` // set limit of length -> max 1000 char ?
	Created  time.Time          `json:"created" bson:"created"`
	Updated  time.Time          `json:"updated" bson:"updated"`
}
