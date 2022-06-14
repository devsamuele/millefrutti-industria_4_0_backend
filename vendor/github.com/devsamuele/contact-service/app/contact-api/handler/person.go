package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/devsamuele/contact-service/business/data/person"
	"github.com/devsamuele/elit/field"
	"github.com/devsamuele/elit/filter"
	"github.com/devsamuele/elit/section"

	"github.com/devsamuele/service-kit/web"
)

type personGroup struct {
	s person.Service
}

func newPersonGroup(s person.Service) personGroup {
	return personGroup{
		s: s,
	}
}

func (pg personGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }

	// claims := auth.Claims{TenantID: TENANT_ID}

	persID := web.URIParams(r)["id"]

	pers, err := pg.s.QueryByID(ctx, TENANT_ID, persID)
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, pers, http.StatusOK)
}

func (pg personGroup) search(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }

	// claims := auth.Claims{TenantID: TENANT_ID}

	var nf filter.Input
	if err := web.Decode(r, &nf); err != nil {
		return err
	}

	pers, err := pg.s.Search(ctx, TENANT_ID, nf)
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, pers, http.StatusOK)
}

func (pg personGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }

	// claims := auth.Claims{TenantID: TENANT_ID}

	persID := web.URIParams(r)["id"]

	var uo person.Update
	if err := web.Decode(r, &uo); err != nil {
		return err
	}

	if err := pg.s.Update(ctx, TENANT_ID, persID, uo, time.Now()); err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (pg personGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}
	var no person.Input
	if err := web.Decode(r, &no); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	pers, err := pg.s.Create(ctx, TENANT_ID, no, time.Now())
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, pers, http.StatusOK)
}

func (pg personGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

	err := pg.s.Delete(ctx, TENANT_ID, id)
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// func (pg personGroup) trash(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	_, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }
// 	// claims := auth.Claims{TenantID: TENANT_ID}

// 	id := web.URIParams(r)["id"]

// 	err := pg.s.Trash(ctx, TENANT_ID, id, time.Now())
// 	if err != nil {
// 		return errHandler(err)
// 	}

// 	return web.Respond(ctx, w, nil, http.StatusNoContent)
// }

// func (pg personGroup) restore(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	_, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }
// 	// claims := auth.Claims{TenantID: TENANT_ID}

// 	id := web.URIParams(r)["id"]

// 	err := pg.s.Restore(ctx, TENANT_ID, id)
// 	if err != nil {
// 		return errHandler(err)
// 	}

// 	return web.Respond(ctx, w, nil, http.StatusNoContent)
// }

func (pg personGroup) querySectionByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	sectionID := web.URIParams(r)["id"]

	section, err := pg.s.QuerySectionByID(ctx, TENANT_ID, sectionID)
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, section, http.StatusOK)
}

func (pg personGroup) querySection(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	sections, err := pg.s.QuerySection(ctx, TENANT_ID)
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, sections, http.StatusOK)
}

func (pg personGroup) createSection(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	var ns section.NewSection
	if err := web.Decode(r, &ns); err != nil {
		return err
	}

	section, err := pg.s.CreateSection(ctx, TENANT_ID, ns, time.Now())
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, section, http.StatusCreated)
}

func (pg personGroup) updateSection(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	sectionID := web.URIParams(r)["id"]

	var us section.UpdateSection
	if err := web.Decode(r, &us); err != nil {
		return err
	}

	err := pg.s.UpdateSection(ctx, TENANT_ID, sectionID, us, time.Now())
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (pg personGroup) deleteSection(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	sectionID := web.URIParams(r)["id"]
	err := pg.s.DeleteSectionByID(ctx, TENANT_ID, sectionID, time.Now())
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (pg personGroup) createField(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	var nf field.Input
	if err := web.Decode(r, &nf); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	field, err := pg.s.InsertField(ctx, TENANT_ID, nf, time.Now())
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, field, http.StatusCreated)
}

func (pg personGroup) updateField(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

	var uf field.Update
	if err := web.Decode(r, &uf); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	err := pg.s.UpdateField(ctx, TENANT_ID, fieldID, uf, time.Now())
	if err != nil {
		if err != nil {
			return errHandler(err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (pg personGroup) queryField(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}
	fields, err := pg.s.QueryField(ctx, TENANT_ID)
	if err != nil {
		if err != nil {
			return errHandler(err)
		}
	}

	return web.Respond(ctx, w, fields, http.StatusOK)
}

func (pg personGroup) queryFieldByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

	field, err := pg.s.QueryFieldByID(ctx, TENANT_ID, fieldID)
	if err != nil {
		if err != nil {
			return errHandler(err)
		}
	}

	return web.Respond(ctx, w, field, http.StatusOK)
}

func (pg personGroup) deleteField(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

	err := pg.s.DeleteFieldByID(ctx, TENANT_ID, fieldID, time.Now())
	if err != nil {
		if err != nil {
			return errHandler(err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
