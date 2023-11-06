package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
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

func NewOrganizationController(s *store.Store, sessionStore *sessions.CookieStore) func(r chi.Router) {
	if oc == nil {
		oc = &OrganizationController{
			organizationService: service.NewOrganizationService(s),
		}
	}

	am := ConfigureMiddleware(sessionStore, "auth")

	return func(r chi.Router) {
		r.With(am.AuthenticateUser).Post("/", oc.CreateOrganization)
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

func (oc *OrganizationController) findOrganization(w http.ResponseWriter, r *http.Request, userId int, orgId string) {
	//org, err := oc.organizationService.GetOrganization(userId, orgId)
	//if err != nil {
	//	Error(w, r, http.StatusNotFound, err)
	//	return
	//}
	fmt.Println(orgId)

	Respond(w, r, http.StatusOK, nil)
}
