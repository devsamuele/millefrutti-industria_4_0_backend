package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/devsamuele/contact-service/business/data/contact/service"
	"github.com/devsamuele/contact-service/business/data/contact/store"
	"github.com/devsamuele/contact-service/business/data/person"
	"github.com/devsamuele/contact-service/business/sys/database"
	"github.com/devsamuele/elit/field"
	"github.com/devsamuele/elit/resource"
	"github.com/devsamuele/elit/section"

	"github.com/streadway/amqp"

	"github.com/devsamuele/service-kit/auth"
	"github.com/devsamuele/service-kit/mid"
	"github.com/devsamuele/service-kit/web"

	"go.mongodb.org/mongo-driver/mongo"
)

var TENANT_ID = "507f1f77bcf86cd799439011"

func API(build string, ch *amqp.Channel, db *mongo.Database, shutdown chan os.Signal, log *log.Logger, a *auth.Auth) *web.Router {
	router := web.NewRouter(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panic(log))
	contactGroup := router.Group("/v1/contact")

	organizationRouter(contactGroup, ch, db, log)
	personRouter(contactGroup, ch, db, log)
	officeRouter(contactGroup, ch, db, log)
	adminRouter(contactGroup, ch, db, log)

	return router
}

func organizationRouter(group *web.Group, ch *amqp.Channel, db *mongo.Database, log *log.Logger) {
	orgSrv := service.NewOrganization(database.NewSession(db), store.NewOrganization(db, log), store.NewOffice(db, log))
	og := newOrganizationGroup(orgSrv)

	group.HandleFn(http.MethodGet, "/organization/search", og.search) // ✅
	group.HandleFn(http.MethodGet, "/organization/:id", og.queryByID) // ✅
	group.HandleFn(http.MethodPost, "/organization", og.create)       // ✅
	group.HandleFn(http.MethodPut, "/organization/:id", og.update)    // ✅
	// group.HandleFn(http.MethodPost, "/organization/:id/trash", og.trash)     // ✅
	// group.HandleFn(http.MethodPost, "/organization/:id/restore", og.restore) // ✅
	group.HandleFn(http.MethodDelete, "/organization/:id", og.delete)

	group.HandleFn(http.MethodGet, "/organizationField", og.queryField)         // ✅
	group.HandleFn(http.MethodGet, "/organizationField/:id", og.queryFieldByID) // ✅
	group.HandleFn(http.MethodPost, "/organizationField", og.createField)       // ✅
	group.HandleFn(http.MethodPut, "/organizationField/:id", og.updateField)    // ✅
	group.HandleFn(http.MethodDelete, "/organizationField/:id", og.deleteField) // ✅

	group.HandleFn(http.MethodGet, "/organizationSection", og.querySection)         // ✅
	group.HandleFn(http.MethodGet, "/organizationSection/:id", og.querySectionByID) // ✅
	group.HandleFn(http.MethodPost, "/organizationSection", og.createSection)       // ✅
	group.HandleFn(http.MethodPut, "/organizationSection/:id", og.updateSection)    // ✅
	group.HandleFn(http.MethodDelete, "/organizationSection/:id", og.deleteSection) // ✅
}

func personRouter(group *web.Group, ch *amqp.Channel, db *mongo.Database, log *log.Logger) {
	persSrv := person.NewService(database.NewSession(db), resource.NewStore(db, log, "person"),
		field.NewStore(db, log, "person"), section.NewStore(db, log, "person"), field.NewStore(db, log, "organization"))
	pg := newPersonGroup(persSrv)

	group.HandleFn(http.MethodPost, "/person/search", pg.search) // ✅
	group.HandleFn(http.MethodGet, "/person/:id", pg.queryByID)  // ✅
	group.HandleFn(http.MethodPost, "/person", pg.create)        // ✅
	group.HandleFn(http.MethodPut, "/person/:id", pg.update)     // ✅
	// group.HandleFn(http.MethodPost, "/person/:id/trash", pg.trash)     // ✅
	// group.HandleFn(http.MethodPost, "/person/:id/restore", pg.restore) // ✅
	group.HandleFn(http.MethodDelete, "/person/:id", pg.delete)

	group.HandleFn(http.MethodGet, "/personField", pg.queryField)         // ✅
	group.HandleFn(http.MethodGet, "/personField/:id", pg.queryFieldByID) // ✅
	group.HandleFn(http.MethodPost, "/personField", pg.createField)       // ✅
	group.HandleFn(http.MethodPut, "/personField/:id", pg.updateField)    // ✅
	group.HandleFn(http.MethodDelete, "/personField/:id", pg.deleteField) // ✅n

	group.HandleFn(http.MethodGet, "/personSection", pg.querySection)         // ✅
	group.HandleFn(http.MethodGet, "/personSection/:id", pg.querySectionByID) // ✅
	group.HandleFn(http.MethodPost, "/personSection", pg.createSection)       // ✅
	group.HandleFn(http.MethodPut, "/personSection/:id", pg.updateSection)    // ✅
	group.HandleFn(http.MethodDelete, "/personSection/:id", pg.deleteSection) // ✅
}

func officeRouter(group *web.Group, ch *amqp.Channel, db *mongo.Database, log *log.Logger) {
	offSrv := service.NewOffice(database.NewSession(db), store.NewOrganization(db, log), store.NewOffice(db, log))
	og := newOfficeGroup(offSrv)

	//group.HandleFn(http.MethodGet, "/office/:id", og.QueryByID) // ✅
	group.HandleFn(http.MethodPost, "/office", og.create)    // ✅
	group.HandleFn(http.MethodPut, "/office/:id", og.update) // ✅
	group.HandleFn(http.MethodDelete, "/office/:id", og.delete)

	group.HandleFn(http.MethodGet, "/officeField", og.queryField)         // ✅
	group.HandleFn(http.MethodGet, "/officeField/:id", og.queryFieldByID) // ✅
	group.HandleFn(http.MethodPost, "/officeField", og.createField)       // ✅
	group.HandleFn(http.MethodPut, "/officeField/:id", og.updateField)    // ✅
	group.HandleFn(http.MethodDelete, "/officeField/:id", og.deleteField) // ✅

}

func adminRouter(group *web.Group, ch *amqp.Channel, db *mongo.Database, log *log.Logger) {
	adminService := service.NewAdmin(database.NewSession(db), resource.NewStore(db, log, "person"), field.NewStore(db, log, "person"), section.NewStore(db, log, "person"))
	og := newAdminGroup(adminService)

	//group.HandleFn(http.MethodGet, "/office/:id", og.QueryByID) // ✅
	group.HandleFn(http.MethodPost, "/admin/buildTenant", og.buildTenant) // ✅
}
