package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"llm-promp-inj.api/internal/dto"
	"llm-promp-inj.api/internal/middleware"
	"llm-promp-inj.api/internal/models"
	"llm-promp-inj.api/internal/service"
)

type ClassificationHandler struct {
	ClssService    *service.ClassificationService
	AuthMiddleware *middleware.AuthMiddleware
}

func NewClassificationHandler(service *service.ClassificationService, authMiddleware *middleware.AuthMiddleware) *ClassificationHandler {
	return &ClassificationHandler{ClssService: service, AuthMiddleware: authMiddleware}
}

// Define the routes of the controller
func (h *ClassificationHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.With(h.AuthMiddleware.Authorize([]string{"admin", "ext_sys"})).Post("/", h.CreateClassificationRequest)
	r.With(h.AuthMiddleware.Authorize([]string{"admin", "ext_sys"})).Get("/logs/{id}", h.GetClassificationRequestByID)
	r.With(h.AuthMiddleware.Authorize([]string{"admin", "ext_sys"})).Get("/logs", h.GetClassificationRequestsByPage)

	return r
}

func (h *ClassificationHandler) CreateClassificationRequest(w http.ResponseWriter, r *http.Request) {
	var classificationRequest dto.ClassificationRequest
	var usernameClaim string

	if err := render.DecodeJSON(r.Body, &classificationRequest); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	userClaimsCtx, ok := r.Context().Value("userClaims").(*models.AccessTokenClaims)
	if !ok {
		usernameClaim = ""
	} else {
		usernameClaim = userClaimsCtx.Data["username"]
	}

	classificationRequestResult, err := h.ClssService.ClassifyText(classificationRequest, usernameClaim)
	if err != nil {

		response := dto.GenericResponse{
			Status:  "Failed",
			Message: err.Error(),
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, classificationRequestResult)
}

func (h *ClassificationHandler) GetClassificationRequestByID(w http.ResponseWriter, r *http.Request) {
	idURLParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idURLParam, 10, 32)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	// Convert the id int parameter to unsigned int.
	idUint := uint(id)

	foundClssRequest, err := h.ClssService.GetClassificationRequestLogByID(idUint)
	if err != nil {
		errResponse := dto.GenericResponse{
			Status:  "Failed for ID",
			Message: err.Error(),
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, foundClssRequest)
}

func (h *ClassificationHandler) GetClassificationRequestsByPage(w http.ResponseWriter, r *http.Request) {
	pageURLParam := r.URL.Query().Get("page")
	limitURLParam := r.URL.Query().Get("limit")
	orderByURLParam := r.URL.Query().Get("sortBy")

	// Default values in case the request does not contain them.
	pageNum := 1
	limit := 15
	var err error

	// Convert the query parameters from string to int and handle error.
	if pageURLParam != "" {
		pageNum, err = strconv.Atoi(pageURLParam)

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}
	}

	if limitURLParam != "" {
		limit, err = strconv.Atoi(limitURLParam)

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}
	}

	allClassificationReqs, err := h.ClssService.GetClassificationLogsByPage(pageNum, limit, orderByURLParam)
	if err != nil {
		errResponse := dto.GenericResponse{
			Status:  "Failed for page",
			Message: err.Error(),
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, errResponse)

		return

	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, allClassificationReqs)
}
