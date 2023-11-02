package controller

import (
	"encoding/json"
	"net/http"
	"relay-backend/internal/model"
	"relay-backend/internal/service"
	"relay-backend/internal/store"
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

func OrganizationHandleFunc(s *store.Store) func(w http.ResponseWriter, r *http.Request) {
	if oc == nil {
		oc = &OrganizationController{
			organizationService: service.NewOrganizationService(s),
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			oc.createOrganization(w, r)
		case http.MethodGet:
		case http.MethodPut:
		case http.MethodDelete:

		}
	}
}

func OrganizationListHandleFunc(s *store.Store) func(w http.ResponseWriter, r *http.Request) {
	if oc == nil {
		oc = &OrganizationController{
			organizationService: service.NewOrganizationService(s),
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (oc *OrganizationController) createOrganization(w http.ResponseWriter, r *http.Request) {
	orgData := &organizationData{}

	if err := json.NewDecoder(r.Body).Decode(orgData); err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	organization := &model.Organization{
		OwnerId:     r.Context().Value(CtxKeyUser).(int),
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
