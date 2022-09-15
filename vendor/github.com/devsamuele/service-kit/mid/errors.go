package mid

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/devsamuele/service-kit/web"
)

type errorResponse struct {
	Error error `json:"error"`
}

// Errors ...
func Errors(log *log.Logger) web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// TRACE ID
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			if err := handler(ctx, w, r); err != nil {

				// LOG ERROR
				log.Printf("\t%s : ERROR     : %v", v.TraceID, err)

				var errResp errorResponse
				var status int
				switch {
				case web.IsRequestError(err):
					reqErr := web.GetRequestError(err)
					errResp = errorResponse{
						Error: reqErr,
					}
					status = reqErr.Code
				default:
					errResp = errorResponse{
						Error: web.NewRequestError(errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError, web.ErrReasonInternalError, "", ""),
					}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, errResp, status); err != nil {
					return err
				}

				// SHUTDOWN SIGNAL
				if ok := web.IsShutdown(err); ok {
					// returning => shutdown
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
