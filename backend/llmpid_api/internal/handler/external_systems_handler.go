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

	r.Post("/auth", h.Auth)
	r.With(h.AuthMiddleware.Authenticate([]string{"user"})).Post("/register", h.Register)
	r.With(h.AuthMiddleware.Authenticate([]string{"user", "ext_sys"})).Put("/logout", h.Deauth)

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
	// 600 days expiration of service auth token. Needs to be withdrawn when not in use.
	accessToken, err := h.AuthService.AuthenticateUser(authServiceRequest.SystemName, authServiceRequest.AccessKey, 36000)
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

	response := dto.GenericResponse{
		Status:  "Success",
		Message: accessToken,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

func (h *ExternalSystemHandler) Register(w http.ResponseWriter, r *http.Request) {
	var authServiceRequest dto.RegisterExtSystemRequest

	if err := render.DecodeJSON(r.Body, &authServiceRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

}

func (h *ExternalSystemHandler) Deauth(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	tokenString := parts[1]

	err := h.AuthService.RevokeUserSession(tokenString)
	if err != nil {
		resp := dto.GenericResponse{Status: "Fail", Message: err.Error()}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp)
		return
	}

	render.Status(r, http.StatusOK)
}
