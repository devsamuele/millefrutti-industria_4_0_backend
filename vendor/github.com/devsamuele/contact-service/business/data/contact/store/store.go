package store

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

const (
	organizationCollection        = "organization"
	organizationFieldCollection   = "organization_field"
	organizationSectionCollection = "organization_section"
	personCollection              = "person"
	personFieldCollection         = "person_field"
	personSectionCollection       = "person_section"
	officeCollection              = "office"
	officeFieldCollection         = "office_field"
	personRoleCollection          = "person_role"
)
