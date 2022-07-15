package pasteurizer

import (
	"fmt"
	"time"
)

var OpcuaConnected bool

const (
	PROCESSING_STATUS_SENT  = "sent"
	PROCESSING_STATUS_WORK  = "work"
	PROCESSING_STATUS_ERROR = "error"
	PROCESSING_STATUS_DONE  = "done"
)

type Work struct {
	ID              int       `json:"id" db:"id"`
	CdLotto         string    `json:"cd_lotto" db:"cd_lotto"`
	CdAr            string    `json:"cd_ar" db:"cd_ar"`
	BasilAmount     int       `json:"basil_amount" db:"basil_amount"`
	Packages        int       `json:"packages" db:"packages"`
	Date            time.Time `json:"date" db:"date"`
	DocumentCreated bool      `json:"document_created" db:"document_created"`
	Status          string    `json:"status" db:"status"`
	Created         time.Time `json:"created" db:"created"`
}

type NewWork struct {
	CdLotto *string `json:"cd_lotto"`
	CdAr    *string `json:"cd_ar"`
}

func (nw NewWork) Validate() error {
	if nw.CdLotto == nil {
		return fmt.Errorf("cd_lotto is required")
	}

	if nw.CdAr == nil {
		return fmt.Errorf("cd_ar is required")
	}
	return nil
}

type ID struct {
	ID int `json:"id"`
}

type OpcuaConnection struct {
	Connected bool `json:"connected"`
}
