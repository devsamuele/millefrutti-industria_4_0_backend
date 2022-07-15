package handler

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/devsamuele/millefrutti-industria_4_0_backend/business/data/pasteurizer"
	"github.com/devsamuele/millefrutti-industria_4_0_backend/business/data/spindryer"
	"github.com/devsamuele/service-kit/mid"
	"github.com/devsamuele/service-kit/web"
	"github.com/devsamuele/service-kit/ws"
	"github.com/gopcua/opcua"
)

func API(build string, db *sql.DB, spindryerClient *opcua.Client, pasteurizerClient *opcua.Client, io *ws.EventEmitter, shutdown chan os.Signal, log *log.Logger) *web.Router {

	handler := io.OnConnection(func(r *http.Request, socket *ws.Socket) {})
	router := web.NewRouter(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panic(log))

	v1 := router.Group("/v1")
	v1.HandleFn(http.MethodGet, "/ws", handler)

	spindryerRouter := v1.SubGroup("/spindryer")
	spindryerGroup := NewSpindryerGroup(spindryer.NewService(spindryer.NewStore(db, log), spindryerClient, io))
	spindryerRouter.HandleFn(http.MethodPost, "/createdDocuments", spindryerGroup.CreatedDocument)
	spindryerRouter.HandleFn(http.MethodPost, "/work", spindryerGroup.InsertWork)
	spindryerRouter.HandleFn(http.MethodGet, "/work", spindryerGroup.QueryWork)
	spindryerRouter.HandleFn(http.MethodGet, "/opcuaConnection", spindryerGroup.GetOpcuaConnection)
	spindryerRouter.HandleFn(http.MethodDelete, "/work/:id", spindryerGroup.DeleteWork)

	pasteurizerRouter := v1.SubGroup("/pasteurizer")
	pasteurizerGroup := NewPasteurizerGroup(pasteurizer.NewService(pasteurizer.NewStore(db, log), pasteurizerClient, io))
	pasteurizerRouter.HandleFn(http.MethodPost, "/createdDocuments", pasteurizerGroup.CreatedDocument)
	pasteurizerRouter.HandleFn(http.MethodPost, "/work", pasteurizerGroup.InsertWork)
	pasteurizerRouter.HandleFn(http.MethodGet, "/work", pasteurizerGroup.QueryWork)
	pasteurizerRouter.HandleFn(http.MethodGet, "/opcuaConnection", pasteurizerGroup.GetOpcuaConnection)
	pasteurizerRouter.HandleFn(http.MethodDelete, "/work/:id", pasteurizerGroup.DeleteWork)

	return router
}

/*
REST API
- insert new lotto -> check if exisit first (add record in xTABLE)
- remove processing record
- setDocumentCreated (passing id of xTable record
*/
