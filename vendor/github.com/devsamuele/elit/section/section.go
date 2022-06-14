package section

import (
	"time"

	"github.com/devsamuele/elit/resperr"
	"github.com/devsamuele/elit/translation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Section ...
type Section struct {
	ID          primitive.ObjectID       `json:"id" bson:"_id"`
	TenantID    primitive.ObjectID       `json:"tenant_id" bson:"tenant_id"`
	Name        string                   `json:"name,omitempty" bson:"name,omitempty"`
	Title       translation.Translation  `json:"title" bson:"title"`
	Description *translation.Translation `json:"description" bson:"description"`
	Custom      bool                     `json:"custom" bson:"custom"`
	Created     time.Time                `json:"created" bson:"created"`
	Updated     time.Time                `json:"updated" bson:"updated"`
}

type NewSection struct {
	Title       *translation.Translation `json:"title"`
	Description *translation.Translation `json:"description"`
}

type UpdateSection struct {
	Title       *translation.Translation `json:"title"`
	Description *translation.Translation `json:"description"`
}

func (ns NewSection) Validate() error {

	if ns.Title == nil || (ns.Title != nil && (*ns.Title).IsEmpty()) {
		return resperr.Error{
			Message:      "title is required",
			Reason:       resperr.ErrReasonRequired,
			LocationType: "argument",
			Location:     "title",
		}
	}

	if !(*ns.Title).GreatherThan(200) {
		return resperr.Error{
			Message:      "max length exceeded",
			Reason:       resperr.ErrReasonInvalidArgument,
			LocationType: "argument",
			Location:     "title",
		}
	}

	return nil
}

func (us UpdateSection) Validate(oldSection Section) error {
	if us.Title != nil {
		if !oldSection.Custom {
			return resperr.Error{
				Message:      "unable to update private title",
				Reason:       resperr.ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "title",
			}
		}

		if (*us.Title).IsEmpty() {
			return resperr.Error{
				Message:      "title is required",
				Reason:       resperr.ErrReasonRequired,
				LocationType: "argument",
				Location:     "title",
			}
		}

		if !(*us.Title).GreatherThan(200) {
			return resperr.Error{
				Message:      "max length exceeded",
				Reason:       resperr.ErrReasonInvalidArgument,
				LocationType: "argument",
				Location:     "title",
			}
		}
	}

	if us.Description != nil {
		if !oldSection.Custom {
			return resperr.Error{
				Message:      "title already exist",
				Reason:       resperr.ErrReasonConflict,
				LocationType: "argument",
				Location:     "title",
			}
		}
	}

	return nil
}
