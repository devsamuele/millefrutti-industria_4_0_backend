package handler

import (
	"context"
	"github.com/devsamuele/contact-service/business/data/contact/service"
	"github.com/devsamuele/service-kit/web"
	"net/http"
	"time"
)

type adminGroup struct {
	s service.Admin
}

func newAdminGroup(s service.Admin) adminGroup {
	return adminGroup{
		s: s,
	}
}

func (ag adminGroup) buildTenant(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	err := ag.s.InitTenant(ctx, TENANT_ID, time.Now())
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
