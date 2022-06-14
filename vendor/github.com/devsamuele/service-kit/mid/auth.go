package mid

import (
	"context"
	"net/http"
	"strings"

	"github.com/devsamuele/service-kit/auth"
	"github.com/devsamuele/service-kit/web"
	"github.com/pkg/errors"
)

// Authenticate ...
func Authenticate(a *auth.Auth) web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			parts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: Bearer <token>")
				return web.NewRequestError(err, http.StatusUnauthorized, "", "", "")
			}

			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				return web.NewRequestError(err, http.StatusUnauthorized, "", "", "")
			}

			ctx = context.WithValue(ctx, auth.Key, claims)

			err = handler(ctx, w, r)
			if err != nil {
				return err
			}

			return nil
		}

		return h
	}

	return m
}

// Authorize ...
func Authorize(scopes ...string) web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// claims, ok := ctx.Value(auth.Key).(auth.Claims)
			// if !ok {
			// 	return errors.New("claims missing from context")
			// }

			resp, err := auth.Authorize(r, scopes...)
			if err != nil {
				return err
			}
			if !resp {
				return web.NewRequestError(
					errors.New("you are not authorized for that action"),
					http.StatusForbidden,
					"", "", "",
				)
			}

			return handler(ctx, w, r)

		}

		return h
	}

	return m
}
