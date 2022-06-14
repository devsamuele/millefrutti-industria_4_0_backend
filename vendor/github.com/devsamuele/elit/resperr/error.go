package resperr

var (
	ErrReasonRequired         = "required"
	ErrReasonInvalidArgument  = "invalidArgument"  // body
	ErrReasonInvalidParameter = "invalidParameter" //url
	ErrReasonConflict         = "conflict"
	ErrReasonNotFound         = "notFound"
)

type Error struct {
	Message      string `json:"message"`
	Reason       string `json:"reason"`
	LocationType string `json:"locationType"` // argument, parameter, ...
	Location     string `json:"location"`     // argumentName, argument parameter...
}

func (e Error) Error() string {
	return e.Message
}
