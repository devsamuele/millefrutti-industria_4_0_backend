package handler

import (
	"database/sql"
	"log"
	"os"

	"github.com/devsamuele/service-kit/mid"
	"github.com/devsamuele/service-kit/web"
)

var TENANT_ID = "507f1f77bcf86cd799439011"

func API(build string, db *sql.DB, shutdown chan os.Signal, log *log.Logger) *web.Router {
	router := web.NewRouter(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panic(log))
	contactGroup := router.Group("/v1/contact")

	// organizationRouter(contactGroup, ch, &organization)

	return router
}
