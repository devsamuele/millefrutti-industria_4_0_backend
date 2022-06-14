package mid

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/devsamuele/service-kit/web"
)

// Logger ...
func Logger(log *log.Logger) web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			fmt.Println("")
			log.Printf("\t%s : started   : %s %s -> %s",
				v.TraceID,
				r.Method, r.URL.Path, r.RemoteAddr,
			)

			err := handler(ctx, w, r)

			log.Printf("\t%s : completed : %s %s -> %s (%d) (%s)",
				v.TraceID,
				r.Method, r.URL.Path, r.RemoteAddr,
				v.StatusCode, time.Since(v.Now),
			)

			return err
		}

		return h
	}

	return m
}
