package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/google/uuid"
)

type ctxKey int

// KeyValues ...
const KeyValues ctxKey = 1

// Values ...
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// Handler ...
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App ...

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.ctxMux.ServeHTTP(w, req)
}

// Router ...
type Router struct {
	ctxMux   *httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

// NewRouter ...
func NewRouter(shutdown chan os.Signal, mw ...Middleware) *Router {
	return &Router{
		ctxMux:   httptreemux.NewContextMux(),
		shutdown: shutdown,
		mw:       mw,
	}
}

// Group ...
type Group struct {
	ctxGroup *httptreemux.ContextGroup
	shutdown chan os.Signal
	mw       []Middleware
}

// Group ...
func (r *Router) Group(pattern string) *Group {
	return &Group{
		ctxGroup: r.ctxMux.NewContextGroup(pattern),
		shutdown: r.shutdown,
		mw:       r.mw,
	}
}

// SubGroup ...
func (g *Group) SubGroup(pattern string) *Group {
	return &Group{
		ctxGroup: g.ctxGroup.NewContextGroup(pattern),
		shutdown: g.shutdown,
		mw:       g.mw,
	}
}

// SignalShutdown ...
func (g *Group) SignalShutdown() {
	g.shutdown <- syscall.SIGTERM
}

func (g *Group) GetCtxGroup() *httptreemux.ContextGroup {
	return g.ctxGroup
}

// HandleFn ...
func (g *Group) HandleFn(method, pattern string, handler Handler, mw ...Middleware) {
	handler = WrapMiddleware(mw, handler)
	handler = WrapMiddleware(g.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {

		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}

		ctx := context.WithValue(r.Context(), KeyValues, &v)

		if err := handler(ctx, w, r); err != nil {
			g.SignalShutdown()
			return
		}
	}

	g.ctxGroup.Handle(method, pattern, h)
}
