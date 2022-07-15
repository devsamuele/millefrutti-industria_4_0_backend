package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/devsamuele/millefrutti-industria_4_0_backend/business/data/spindryer"
	"github.com/devsamuele/service-kit/web"
)

type SpindryerGroup struct {
	srv spindryer.Service
}

func NewSpindryerGroup(srv spindryer.Service) SpindryerGroup {
	return SpindryerGroup{
		srv: srv,
	}
}

func (g SpindryerGroup) QueryWork(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	work, err := g.srv.QueryWork(ctx)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, work, http.StatusOK)
}

func (g SpindryerGroup) GetOpcuaConnection(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	conn := g.srv.GetOpcuaConnection(ctx)

	return web.Respond(ctx, w, conn, http.StatusOK)
}

func (g SpindryerGroup) DeleteWork(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	id := web.URIParams(r)["id"]

	err := g.srv.DeleteWork(ctx, id)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (g SpindryerGroup) InsertWork(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	var nw spindryer.NewWork
	if err := web.Decode(r, &nw); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	work, err := g.srv.InsertWork(ctx, nw, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, work, http.StatusCreated)
}

func (g SpindryerGroup) CreatedDocument(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	ids := make([]spindryer.ID, 0)
	if err := web.Decode(r, &ids); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	err := g.srv.SetCreatedDocument(ctx, ids)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
