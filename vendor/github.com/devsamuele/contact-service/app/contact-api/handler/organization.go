package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/devsamuele/contact-service/business/data/contact/model"
	"github.com/devsamuele/contact-service/business/data/contact/service"

	"github.com/devsamuele/service-kit/web"
)

type organizationGroup struct {
	s service.Organization
}

func newOrganizationGroup(s service.Organization) organizationGroup {
	return organizationGroup{
		s: s,
	}
}

// query ... Get all organization
// @Summary Get all organization
// @Description get all users
// @Tags organization
// @Success 200 {array} organization.organization
// @Failure 404 {object} object
// @Router /organization [get]
// func (og organizationGroup) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }

// 	claims := auth.Claims{TenantID: TENANT_ID}

// 	var queryOpts organization.QueryOpts
// 	if err := web.Decode(r, &queryOpts); err != nil {
// 		return err
// 	}

// 	organizations, err := og.store.Query(ctx, v.TraceID, claims, queryOpts)
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	return web.Respond(ctx, w, organizations, http.StatusOK)
// }

// func (og organizationGroup) search(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	if !ok {
// 		return web.NewShutdownError("claims missing from context")
// 	}

// 	URIparams := web.URIParams(r)
// 	pageNumber, err := strconv.Atoi(URIparams["page"])
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	rowsPerPage, err := strconv.Atoi(URIparams["rows"])
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	trashedParam := web.QueryParams(r)["trashed"]
// 	if trashedParam != "true" {
// 		trashedParam = "false"
// 	}

// 	sortParam := web.QueryParams(r)["sort"]
// 	if sortParam != "DESC" {
// 		sortParam = "ASC"
// 	}

// 	var searchFilter organization.NewSearchFilter
// 	if err := web.Decode(r, &searchFilter); err != nil {
// 		return web.NewRequestError(err, http.StatusBadRequest)
// 	}

// 	opts := organization.QueryOptions{
// 		PageNumber:     pageNumber,
// 		organizationStatus: organization.STATUS_UNTRASHED,
// 		organizationLimit:  100,
// 		RowsPerPage:    rowsPerPage,
// 		Sort: organization.Sort{
// 			FieldEnName: "name",
// 			Order:       "asc",
// 		},
// 	}

// 	organizations, err := og.store.Search(ctx, v.TraceID, claims, nil, searchFilter, opts, nil, "name")
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	return web.Respond(ctx, w, organizations, http.StatusOK)
// }

// func (og organizationGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	_, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }

// 	claims := auth.Claims{TenantID: TENANT_ID}

// 	URIparams := web.URIParams(r)
// 	orgID := URIparams["id"]

// 	organization, err := og.svc(ctx, TENANT_ID, orgID)
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusNotFound)
// 	}

// 	return web.Respond(ctx, w, organization, http.StatusOK)
// }

// func (og organizationGroup) checkByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	_, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }

// 	// claims := auth.Claims{TenantID: TENANT_ID}

// 	URIparams := web.URIParams(r)
// 	orgID := URIparams["id"]

// 	err := og.svc.CheckID(ctx, TENANT_ID, orgID)
// 	if err != nil {
// 		return errHandler()(err)
// 	}

// 	return web.Respond(ctx, w, nil, http.StatusOK)
// }

// func (og organizationGroup) queryOrganization(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	if !ok {
// 		return web.NewShutdownError("claims missing from context")
// 	}

// 	params := web.URIParams(r)
// 	orgStrID := params["id"]
// 	orgID, err := database.ParseID(orgStrID)
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	trashedParam := web.QueryParams(r)["trashed"]
// 	if trashedParam != "true" {
// 		trashedParam = "false"
// 	}

// 	trashed, err := strconv.ParseBool(trashedParam)
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	org, err := og.store.QueryByID(ctx, v.TraceID, claims, orgID, trashed)
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusNotFound)
// 	}

// 	return web.Respond(ctx, w, org, http.StatusOK)
// }

// func (og organizationGroup) queryPreview(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	if !ok {
// 		return web.NewShutdownError("claims missing from context")
// 	}

// 	params := web.URIParams(r)
// 	pageNumber, err := strconv.Atoi(params["page"])
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusBadRequest)
// 	}

// 	rowsPerPage, err := strconv.Atoi(params["rows"])
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusBadRequest)
// 	}

// 	trashedParam := web.QueryParams(r)["trashed"]
// 	if trashedParam != "true" {
// 		trashedParam = "false"
// 	}

// 	sortParam := web.QueryParams(r)["sort"]
// 	if strings.ToUpper(sortParam) != "DESC" {
// 		sortParam = "ASC"
// 	}

// 	opts := organization.QueryOptions{
// 		organizationStatus: organization.STATUS_TRASHED,
// 		organizationLimit:  100,
// 		PageNumber:     pageNumber,
// 		RowsPerPage:    rowsPerPage,
// 		Sort: organization.Sort{
// 			FieldEnName: "name",
// 			Order:       "asc",
// 		},
// 	}
// 	previews, err := og.store.QueryPreview(ctx, v.TraceID, claims, nil, opts, "name")
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	return web.Respond(ctx, w, previews, http.StatusOK)
// }

// func (og organizationGroup) queryPreviewByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	if !ok {
// 		return web.NewShutdownError("claims missing from context")
// 	}

// 	uriParams := web.URIParams(r)
// 	orgID, err := entity.ParseID(uriParams["id"])
// 	if err != nil {
// 		return err
// 	}

// 	qParams := web.QueryParams(r)["trashed"]
// 	if qParams != "true" {
// 		qParams = "false"
// 	}

// 	preview, err := og.store.QueryPreviewByID(ctx, v.TraceID, claims, nil, orgID, organization.STATUS_UNTRASHED, "name")
// 	if err != nil {
// 		return err
// 	}

// 	return web.Respond(ctx, w, preview, http.StatusOK)
// }

func (og organizationGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

	org, err := og.s.QueryByID(ctx, TENANT_ID, orgID)
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, org, http.StatusOK)
}

func (og organizationGroup) search(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }

	// claims := auth.Claims{TenantID: TENANT_ID}

	var nf model.NewFilter
	if err := web.Decode(r, &nf); err != nil {
		return err
	}

	org, err := og.s.Search(ctx, TENANT_ID, nf)
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, org, http.StatusOK)
}

func (og organizationGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

	var uo model.UpdateOrganization
	if err := web.Decode(r, &uo); err != nil {
		return err
	}

	if err := og.s.Update(ctx, TENANT_ID, orgID, uo, time.Now()); err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// func (og organizationGroup) updateMany(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }
// 	claims := auth.Claims{TenantID: TENANT_ID}

// 	var urs organization.Updates
// 	if err := web.Decode(r, &urs, og.validate); err != nil {
// 		return err
// 	}

// 	if err := og.store.UpdateMany(ctx, v.TraceID, claims, urs); err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	return nil
// }

// func (og organizationGroup) createMany(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }
// 	claims := auth.Claims{TenantID: TENANT_ID}

// 	var nfs organization.News
// 	if err := web.Decode(r, &nfs, og.validate); err != nil {
// 		return err
// 	}

// 	orgs, err := og.store.CreateMany(ctx, v.TraceID, claims, nfs, v.Now)
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	return web.Respond(ctx, w, orgs, http.StatusCreated)
// }

func (og organizationGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}
	var no model.NewOrganization
	if err := web.Decode(r, &no); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	org, err := og.s.Create(ctx, TENANT_ID, no, time.Now())
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, org, http.StatusOK)
}

// func (og organizationGroup) Search(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }
// 	claims := auth.Claims{TenantID: TENANT_ID}

// 	URIparams := web.URIParams(r)

// 	pageNumber, err := strconv.Atoi(URIparams["page"])
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	itemsPerPage, err := strconv.Atoi(URIparams["rows"])
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	trashed := organization.STATUS_TRASHED
// 	trashedParam := web.QueryParams(r)["trashed"]
// 	if trashedParam != "true" {
// 		trashedParam = "false"
// 		trashed = organization.STATUS_UNTRASHED
// 	}

// 	var sf organization.NewSearchFilter
// 	if err := web.Decode(r, &sf, og.validate); err != nil {
// 		return err
// 	}

// 	sort := make([]organization.Sort, 0)
// 	sortBy := web.QueryParams(r)["sort_by"]
// 	if sortBy != "" {
// 		sortArray := strings.Split(sortBy, ",")
// 		for _, s := range sortArray {
// 			left := strings.Index(s, "[")
// 			right := strings.Index(s, "]")
// 			pathToValue := s
// 			order := -1
// 			if left != -1 && right != -1 {
// 				pathToValue = s[:left]
// 				orderStr := s[left+1 : right]
// 				if strings.ToLower(orderStr) != "desc" {
// 					order = 1
// 				}
// 			}

// 			sort = append(sort, organization.Sort{
// 				PathToValue: pathToValue,
// 				Order:       order,
// 			})
// 		}
// 	}

// 	opts := organization.QueryOpts{
// 		OrganizationStatus: trashed,
// 		OrganizationLimit:  100,
// 		PageNumber:         pageNumber,
// 		ItemsPerPage:       itemsPerPage,
// 		Sort:               sort,
// 	}

// 	organizations, err := og.store.Search(ctx, v.TraceID, claims, sf, opts)
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	return web.Respond(ctx, w, organizations, http.StatusOK)
// }

// func (og organizationGroup) deleteMany(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }
// 	claims := auth.Claims{TenantID: TENANT_ID}

// 	var ids organization.ObjectIDs
// 	if err := web.Decode(r, &ids); err != nil {
// 		return err
// 	}

// 	err := og.store.DeleteByIDs(ctx, v.TraceID, claims, ids)
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	return web.Respond(ctx, w, nil, http.StatusNoContent)
// }

func (og organizationGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

// func (og organizationGroup) trashMany(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }
// 	claims := auth.Claims{TenantID: TENANT_ID}

// 	var ids organization.ObjectIDs
// 	if err := web.Decode(r, &ids); err != nil {
// 		return err
// 	}

// 	err := og.store.TrashByIDs(ctx, v.TraceID, claims, ids, time.Now())
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	return web.Respond(ctx, w, nil, http.StatusNoContent)
// }

// func (og organizationGroup) trash(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

// 	err := og.s.Trash(ctx, TENANT_ID, id, time.Now())
// 	if err != nil {
// 		return errHandler(err)
// 	}

// 	return web.Respond(ctx, w, nil, http.StatusNoContent)
// }

// func (og organizationGroup) restoreMany(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }
// 	claims := auth.Claims{TenantID: TENANT_ID}

// 	var ids organization.ObjectIDs
// 	if err := web.Decode(r, &ids); err != nil {
// 		return err
// 	}

// 	err := og.store.RestoreByIDs(ctx, v.TraceID, claims, ids)
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	return web.Respond(ctx, w, nil, http.StatusNoContent)
// }

// func (og organizationGroup) restore(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

// 	err := og.s.Restore(ctx, TENANT_ID, id)
// 	if err != nil {
// 		return errHandler(err)
// 	}

// 	return web.Respond(ctx, w, nil, http.StatusNoContent)
// }

func (og organizationGroup) querySectionByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

	section, err := og.s.QuerySectionByID(ctx, TENANT_ID, sectionID)
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, section, http.StatusOK)
}

func (og organizationGroup) querySection(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	sections, err := og.s.QuerySection(ctx, TENANT_ID)
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, sections, http.StatusOK)
}

// func (og organizationGroup) createManySection(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }
// 	claims := auth.Claims{TenantID: TENANT_ID}

// 	var newSections organization.NewSections
// 	if err := web.Decode(r, &newSections, og.validate); err != nil {
// 		return err
// 	}

// 	sections, err := og.store.CreateManySection(ctx, v.TraceID, claims, newSections, time.Now())
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	return web.Respond(ctx, w, sections, http.StatusCreated)
// }

func (og organizationGroup) createSection(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	var ns model.NewSection
	if err := web.Decode(r, &ns); err != nil {
		return err
	}

	section, err := og.s.CreateSection(ctx, TENANT_ID, ns, time.Now())
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, section, http.StatusCreated)
}

// func (og organizationGroup) updateManySection(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }
// 	claims := auth.Claims{TenantID: TENANT_ID}

// 	var updateSections organization.UpdateSections
// 	if err := web.Decode(r, &updateSections, og.validate); err != nil {
// 		return err
// 	}

// 	err := og.store.UpdateManySection(ctx, v.TraceID, claims, updateSections, time.Now())
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	return web.Respond(ctx, w, nil, http.StatusNoContent)
// }

func (og organizationGroup) updateSection(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

	var us model.UpdateSection
	if err := web.Decode(r, &us); err != nil {
		return err
	}

	err := og.s.UpdateSection(ctx, TENANT_ID, sectionID, us, time.Now())
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// func (og organizationGroup) deleteManySection(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

// 	v, ok := ctx.Value(web.KeyValues).(*web.Values)
// 	if !ok {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
// 	// if !ok {
// 	// 	return web.NewShutdownError("claims missing from context")
// 	// }
// 	claims := auth.Claims{TenantID: TENANT_ID}

// 	var sectionIDs organization.ObjectIDs
// 	if err := web.Decode(r, &sectionIDs); err != nil {
// 		return err
// 	}

// 	err := og.store.DeleteSectionByIDs(ctx, v.TraceID, claims, sectionIDs)
// 	if err != nil {
// 		return web.NewRequestError(err, http.StatusInternalServerError)
// 	}

// 	return web.Respond(ctx, w, nil, http.StatusNoContent)
// }

func (og organizationGroup) deleteSection(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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
	err := og.s.DeleteSectionByID(ctx, TENANT_ID, sectionID, time.Now())
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (og organizationGroup) createField(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return web.NewShutdownError("claims missing from context")
	// }
	// claims := auth.Claims{TenantID: TENANT_ID}

	var nf model.NewFieldWithSection
	if err := web.Decode(r, &nf); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	field, err := og.s.InsertField(ctx, TENANT_ID, nf, time.Now())
	if err != nil {
		return errHandler(err)
	}

	return web.Respond(ctx, w, field, http.StatusCreated)
}

func (og organizationGroup) updateField(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

	var uf model.UpdateFieldWithSection
	if err := web.Decode(r, &uf); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	err := og.s.UpdateField(ctx, TENANT_ID, fieldID, uf, time.Now())
	if err != nil {
		if err != nil {
			return errHandler(err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (og organizationGroup) queryField(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

func (og organizationGroup) queryFieldByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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

func (og organizationGroup) deleteField(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

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
