package handler

import (
	"errors"
	"net/http"

	"github.com/devsamuele/contact-service/business/data/contact/model"
	"github.com/devsamuele/elit/resperr"
	"github.com/devsamuele/service-kit/web"
)

func errHandler(err error) error {
	var restErr resperr.Error

	switch {
	case errors.As(err, &restErr):
		switch restErr.Reason {
		case resperr.ErrReasonNotFound:
			return web.NewRequestError(restErr, http.StatusNotFound, restErr.Reason, restErr.LocationType, restErr.Location)
		case resperr.ErrReasonRequired, model.ErrReasonInvalidParameter, model.ErrReasonInvalidArgument:
			return web.NewRequestError(restErr, http.StatusBadRequest, restErr.Reason, restErr.LocationType, restErr.Location)
		case resperr.ErrReasonConflict:
			return web.NewRequestError(restErr, http.StatusConflict, restErr.Reason, restErr.LocationType, restErr.Location)
		default:
			return err
		}
	default:
		return err
	}
}
