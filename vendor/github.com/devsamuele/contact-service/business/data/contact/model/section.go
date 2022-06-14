package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Section ...
type Section struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	TenantID    primitive.ObjectID `json:"tenant_id" bson:"tenant_id"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Title       Translation        `json:"title" bson:"title"`
	Description *Translation       `json:"description" bson:"description"`
	Custom      bool               `json:"custom" bson:"custom"`
	Created     time.Time          `json:"created" bson:"created"`
	Updated     time.Time          `json:"updated" bson:"updated"`
}

type NewSection struct {
	Title       *Translation `json:"title"`
	Description *Translation `json:"description"`
}

type UpdateSection struct {
	Title       *Translation `json:"title"`
	Description *Translation `json:"description"`
}

func (ns NewSection) Validate() error {

	if ns.Title == nil || (ns.Title != nil && (*ns.Title).IsEmpty()) {
		return Error{
			Message:      "title is required",
			Reason:       ErrReasonRequired,
			LocationType: "argument",
			Location:     "title",
		}
	}

	if !(*ns.Title).IsValidLength(200) {
		return Error{
			Message:      "max length exceeded",
			Reason:       ErrReasonInvalidArgument,
			LocationType: "argument",
			Location:     "title",
		}
	}

	return nil
}

func (us UpdateSection) Validate(oldSection Section) error {
	if us.Title != nil {
		if !oldSection.Custom {
			return Error{
				Message:      "unable to update private title",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "title",
			}
		}

		if (*us.Title).IsEmpty() {
			return Error{
				Message:      "title is required",
				Reason:       ErrReasonRequired,
				LocationType: "argument",
				Location:     "title",
			}
		}

		if !(*us.Title).IsValidLength(200) {
			return Error{
				Message:      "max length exceeded",
				Reason:       ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "title",
			}
		}
	}

	if us.Description != nil {
		if !oldSection.Custom {
			return Error{
				Message:      "title already exist",
				Reason:       ErrReasonConflict,
				LocationType: "argument",
				Location:     "title",
			}
		}
	}

	return nil
}
