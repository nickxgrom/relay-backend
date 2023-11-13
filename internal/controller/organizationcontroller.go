package controller

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"relay-backend/internal/model"
	"relay-backend/internal/service"
	"relay-backend/internal/store"
	"strconv"
)

type OrganizationController struct {
	organizationService *service.OrganizationService
}

type organizationData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Email       string `json:"email"`
}

var (
	oc *OrganizationController
)

func NewOrganizationController(s *store.Store, middleware *AuthMiddleware) func(r chi.Router) {
	if oc == nil {
		oc = &OrganizationController{
			organizationService: service.NewOrganizationService(s),
		}
	}

	return func(r chi.Router) {
		r.Use(middleware.AuthenticateUser)
		r.Post("/", oc.CreateOrganization)
		r.Get("/{orgId}", oc.findOrganization)
		r.Get("/{page}/{pageSize}", oc.getOrganizationList)

		r.Post("/{orgId}/employees", oc.addEmployees)
	}
}

func (oc *OrganizationController) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(CtxKeyUser).(int)

	orgData := &organizationData{}

	if err := json.NewDecoder(r.Body).Decode(orgData); err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	organization := &model.Organization{
		OwnerId:     userId,
		Name:        orgData.Name,
		Description: orgData.Description,
		Address:     orgData.Address,
		Email:       orgData.Email,
	}

	err := oc.organizationService.CreateOrganization(organization)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	Respond(w, r, http.StatusCreated, organization)
}

func (oc *OrganizationController) findOrganization(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(CtxKeyUser).(int)
	orgId, err := strconv.Atoi(chi.URLParam(r, "orgId"))

	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	org, err := oc.organizationService.GetOrganization(userId, orgId)
	if err != nil {
		Error(w, r, http.StatusNotFound, err)
		return
	}

	Respond(w, r, http.StatusOK, org)
}

func (oc *OrganizationController) getOrganizationList(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(CtxKeyUser).(int)

	//TODO: consider move to method
	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	pageSize, err := strconv.Atoi(chi.URLParam(r, "pageSize"))
	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	if pageSize == 0 {
		Error(w, r, http.StatusBadRequest, errors.New("page-size-must-be-greater-than-zero"))
		return
	}

	orgList, err := oc.organizationService.GetOrganizationList(userId, page, pageSize)

	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	Respond(w, r, http.StatusOK, orgList)
}

func (oc *OrganizationController) addEmployees(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(CtxKeyUser).(int)
	orgId, err := strconv.Atoi(chi.URLParam(r, "orgId"))
	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	var employeeIds *[]int

	if err := json.NewDecoder(r.Body).Decode(&employeeIds); err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	err = oc.organizationService.AddOrganizationEmployees(userId, orgId, *employeeIds)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	Respond(w, r, http.StatusOK, nil)
}
