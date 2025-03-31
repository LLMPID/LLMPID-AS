package handler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"llm-promp-inj.api/config"
	"llm-promp-inj.api/internal/dto"
	"llm-promp-inj.api/internal/middleware"
	"llm-promp-inj.api/internal/service"
)

type UserHandler struct {
	UserService    *service.UserService
	AuthService    *service.AuthenticationService
	AuthMiddleware *middleware.AuthMiddleware

	Config *config.Config
}

func NewUserHandler(authService *service.AuthenticationService, authMiddleware *middleware.AuthMiddleware) *UserHandler {
	return &UserHandler{
		AuthService:    authService,
		AuthMiddleware: authMiddleware,
	}
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.With(h.AuthMiddleware.Authenticate([]string{"user"})).Put("/logout", h.Logout)
		//r.With(h.AuthMiddleware.Authenticate([]string{"user"})).Post("/register", h.Register)
		r.With(h.AuthMiddleware.Authenticate([]string{"user"})).Post("/credentials/change", h.ChangePassword)
	})
	return r
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.AuthUserRequest

	if err := render.DecodeJSON(r.Body, &loginRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	jwt, err := h.AuthService.Authenticate(loginRequest.Usernames, loginRequest.Password, 60)
	if err != nil {

		response := dto.GenericResponse{
			Status:  "Failed to authenticate user",
			Message: err.Error(),
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response)
		return
	}

	if jwt == "" {
		response := dto.GenericResponse{
			Status:  "Unauthorized",
			Message: "Wrong credetials",
		}

		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, response)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"status": "Success", "access_token": jwt})
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest dto.AuthUserRequest

	if err := render.DecodeJSON(r.Body, &registerRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	user, err := h.UserService.CreateUser(registerRequest.Usernames, registerRequest.Usernames, "user")
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Failed to create user"})
		return
	}

	userDTO := dto.UserResponse{
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, userDTO)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var changePasswordRequest dto.ChangeCredentialsRequest

	// Get auth header
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	tokenString := parts[1]

	if err := render.DecodeJSON(r.Body, &changePasswordRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	// Change password and retrieve the new access token for the new session.
	newJWT, err := h.AuthService.ChangePassword(
		changePasswordRequest.Username,
		changePasswordRequest.OldPassword,
		changePasswordRequest.NewPassword,
		tokenString,
	)
	if err != nil {
		resp := dto.GenericResponse{Status: "Fail", Message: err.Error()}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp)
		return
	}

	if newJWT == "" {
		resp := dto.GenericResponse{Status: "Fail", Message: "Credentials issue."}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"status": "Success", "access_token": newJWT})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	tokenString := parts[1]

	revokeAllTrigger := r.URL.Query().Get("all")

	var err error

	if revokeAllTrigger != "" {
		err = h.AuthService.RevokeAllSessions(tokenString)
	} else {
		err = h.AuthService.RevokeSession(tokenString)
	}

	if err != nil {
		resp := dto.GenericResponse{Status: "Fail", Message: err.Error()}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp)
		return
	}

	render.Status(r, http.StatusOK)
}
