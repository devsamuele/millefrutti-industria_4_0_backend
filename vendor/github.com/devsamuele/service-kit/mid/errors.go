package mid

import (
	"context"
	"log"
	"net/http"

	"github.com/devsamuele/service-kit/web"
)

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

				var errResp web.ErrorResponse
				var status int
				switch {
				case web.IsRequestError(err):
					reqErr := web.GetRequestError(err)
					errResp = web.ErrorResponse{
						Error: *reqErr,
					}
					status = reqErr.Code
				default:
					errResp = web.ErrorResponse{
						Error: web.RequestError{
							Code:    http.StatusInternalServerError,
							Message: http.StatusText(http.StatusInternalServerError),
						},
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
