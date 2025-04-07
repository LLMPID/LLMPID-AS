package handler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"llm-promp-inj.api/internal/dto"
	"llm-promp-inj.api/internal/middleware"
	"llm-promp-inj.api/internal/service"
)

type ExternalSystemHandler struct {
	ExternalSysService *service.ExternalSystemService
	AuthService        *service.AuthenticationService
	AuthMiddleware     *middleware.AuthMiddleware
}

func NewExternalSystemHandler(externalSysService *service.ExternalSystemService, authService *service.AuthenticationService, authMiddleware *middleware.AuthMiddleware) *ExternalSystemHandler {
	return &ExternalSystemHandler{
		ExternalSysService: externalSysService,
		AuthService:        authService,
		AuthMiddleware:     authMiddleware,
	}
}

func (h *ExternalSystemHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.With(h.AuthMiddleware.Authenticate([]string{"admin"})).Post("/", h.Create)
	r.With(h.AuthMiddleware.Authenticate([]string{"admin"})).Get("/", h.List)
	r.With(h.AuthMiddleware.Authenticate([]string{"admin"})).Delete("/{system_name}", h.Delete)
	r.With(h.AuthMiddleware.Authenticate([]string{"admin"})).Put("/{system_name}", h.Update)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/authenticate", h.Auth)
		r.With(h.AuthMiddleware.Authenticate([]string{"admin"})).Put("/deauthenticate/{system_name}", h.DeauthByName)
		r.With(h.AuthMiddleware.Authenticate([]string{"ext_sys"})).Put("/deauthenticate", h.Deauth)
	})
	return r
}

func (h *ExternalSystemHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var authServiceRequest dto.AuthExtSystemRequest

	if err := render.DecodeJSON(r.Body, &authServiceRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	// External systems are treated as an user with "ext_sys" role.
	// 500 days expiration of service auth token. Needs to be withdrawn when not in use.
	accessToken, err := h.AuthService.Authenticate(authServiceRequest.SystemName, authServiceRequest.AccessKey, 720000)
	if err != nil {
		response := dto.GenericResponse{
			Status:  "Failed to authenticate service",
			Message: err.Error(),
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response)
		return
	}

	if accessToken == "" {
		response := dto.GenericResponse{
			Status:  "Unauthorized",
			Message: "Wrong credetials",
		}

		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"status": "Success", "access_token": accessToken})
}

func (h *ExternalSystemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var response dto.GenericResponse
	var registerRequest dto.RegisterExtSystemRequest

	if err := render.DecodeJSON(r.Body, &registerRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	accessKey, err := h.ExternalSysService.Register(registerRequest.SystemName)
	if err != nil {
		response = dto.GenericResponse{
			Status:  "Fail",
			Message: err.Error(),
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"status": "Success", "access_key": accessKey})
}
func (h *ExternalSystemHandler) List(w http.ResponseWriter, r *http.Request) {
	servicesNames, err := h.ExternalSysService.List()
	if err != nil {
		response := dto.GenericResponse{
			Status:  "Fail",
			Message: err.Error(),
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, servicesNames)
}

func (h *ExternalSystemHandler) Update(w http.ResponseWriter, r *http.Request) {
	var updateExternalSystemRequest dto.UpdateExtSystemRequest // The update request is the same as the register one. We use the same DTO.

	if err := render.DecodeJSON(r.Body, &updateExternalSystemRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	err := h.ExternalSysService.Update(updateExternalSystemRequest.OldSystemName, updateExternalSystemRequest.NewSystemName)
	if err != nil {
		resp := dto.GenericResponse{Status: "Fail", Message: err.Error()}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp)
		return
	}

	err = h.AuthService.RevokeAllSessionsByUsername(updateExternalSystemRequest.OldSystemName)
	if err != nil {
		resp := dto.GenericResponse{Status: "Fail", Message: "failed to revoke sessions"}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"status": "Success", "data": ""})
}

func (h *ExternalSystemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	systemName := chi.URLParam(r, "system_name")

	err := h.ExternalSysService.DeleteBySysName(systemName)
	if err != nil {
		resp := dto.GenericResponse{Status: "Fail", Message: err.Error()}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp)
		return
	}

	err = h.AuthService.RevokeAllSessionsByUsername(systemName)
	if err != nil {
		resp := dto.GenericResponse{Status: "Fail", Message: "failed to revoke sessions"}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp)
		return
	}

	render.Status(r, http.StatusOK)
}

func (h *ExternalSystemHandler) Deauth(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	tokenString := parts[1]

	err := h.AuthService.RevokeAllSessionsByToken(tokenString)
	if err != nil {
		resp := dto.GenericResponse{Status: "Fail", Message: err.Error()}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp)
		return
	}

	render.Status(r, http.StatusOK)
}

func (h *ExternalSystemHandler) DeauthByName(w http.ResponseWriter, r *http.Request) {
	systemName := chi.URLParam(r, "system_name")

	// systemName == Username in the context of the authentication service.
	err := h.AuthService.RevokeAllSessionsByUsername(systemName)
	if err != nil {
		resp := dto.GenericResponse{Status: "Fail", Message: err.Error()}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp)
		return
	}

	render.Status(r, http.StatusOK)
}
