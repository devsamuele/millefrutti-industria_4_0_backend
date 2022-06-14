package handler

import (
	"context"
	"fmt"
	"github.com/devsamuele/contact-service/business/data/contact/model"
	"github.com/devsamuele/contact-service/business/data/contact/service"
	"net/http"
	"time"

	"github.com/devsamuele/service-kit/web"
)

type officeGroup struct {
	s service.Office
}

func newOfficeGroup(s service.Office) officeGroup {
	return officeGroup{
		s: s,
	}
}

func (og officeGroup) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }

	// claims := auth.Claims{TenantID: TENANT_ID}

	officeID := web.URIParams(r)["id"]

	var uo model.UpdateOffice
	if err := web.Decode(r, &uo); err != nil {
		return err
	}

	o, err := og.s.QueryByID(ctx, TENANT_ID, officeID)
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, o, http.StatusOK)
}

func (og officeGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }

	// claims := auth.Claims{TenantID: TENANT_ID}

	orgID := web.URIParams(r)["id"]

	var uo model.UpdateOffice
	if err := web.Decode(r, &uo); err != nil {
		return err
	}

	if err := og.s.Update(ctx, TENANT_ID, orgID, uo, time.Now()); err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (og officeGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}
	var no model.NewOffice
	if err := web.Decode(r, &no); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	org, err := og.s.Create(ctx, TENANT_ID, no, time.Now())
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, org, http.StatusCreated)
}

func (og officeGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	id := web.URIParams(r)["id"]

	err := og.s.Delete(ctx, TENANT_ID, id)
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (og officeGroup) createField(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}
	var nf model.NewField
	if err := web.Decode(r, &nf); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	field, err := og.s.InsertField(ctx, TENANT_ID, nf, time.Now())
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, field, http.StatusCreated)
}

func (og officeGroup) updateField(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	fieldID := web.URIParams(r)["id"]

	var ul model.UpdateField
	if err := web.Decode(r, &ul); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	err := og.s.UpdateField(ctx, TENANT_ID, fieldID, ul, time.Now())
	if err != nil {
		if err != nil {
			return errHandler(err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (og officeGroup) queryField(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}
	fields, err := og.s.QueryField(ctx, TENANT_ID)
	if err != nil {
		if err != nil {
			return errHandler(err)
		}
	}

	return web.Respond(ctx, w, fields, http.StatusOK)
}

func (og officeGroup) queryFieldByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	fieldID := web.URIParams(r)["id"]

	field, err := og.s.QueryFieldByID(ctx, TENANT_ID, fieldID)
	if err != nil {
		if err != nil {
			return errHandler(err)
		}
	}

	return web.Respond(ctx, w, field, http.StatusOK)
}

func (og officeGroup) deleteField(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	fieldID := web.URIParams(r)["id"]

	err := og.s.DeleteFieldByID(ctx, TENANT_ID, fieldID, time.Now())
	if err != nil {
		if err != nil {
			return errHandler(err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
