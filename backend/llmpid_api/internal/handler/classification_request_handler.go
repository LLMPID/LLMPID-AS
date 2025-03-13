package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"llm-promp-inj.api/internal/dto"
	"llm-promp-inj.api/internal/service"
)

type ClassificationRequestHandler struct {
	Service *service.ClassificationRequestService
}

func NewClassificationRequestHandler(service *service.ClassificationRequestService) *ClassificationRequestHandler {
	return &ClassificationRequestHandler{Service: service}
}

// Define the routes of the controller
func (h *ClassificationRequestHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.CreateClassificationRequest)
	r.Get("/logs/{id}", h.GetClassificationRequestByID)
	r.Get("/logs", h.GetClassificationRequestsByPage)

	return r
}

func (h *ClassificationRequestHandler) CreateClassificationRequest(w http.ResponseWriter, r *http.Request) {
	var classificationRequest dto.ClassificationRequest

	if err := render.DecodeJSON(r.Body, &classificationRequest); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	classificationRequestResult, err := h.Service.ClassifyText(classificationRequest)
	if err != nil {

		response := dto.Error{
			Error:   "Failed to classify data.",
			Message: err.Error(),
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, classificationRequestResult)
}

func (h *ClassificationRequestHandler) GetClassificationRequestByID(w http.ResponseWriter, r *http.Request) {
	idURLParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idURLParam, 10, 32)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	// Convert the id int parameter to unsigned int.
	idUint := uint(id)

	foundClssRequest, err := h.Service.GetClassificationRequestLogByID(idUint)
	if err != nil {
		errResponse := dto.Error{
			Error:   "Unable to retrieve log with such id.",
			Message: err.Error(),
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, foundClssRequest)
}

func (h *ClassificationRequestHandler) GetClassificationRequestsByPage(w http.ResponseWriter, r *http.Request) {
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

	allClassificationReqs, err := h.Service.GetClassificationLogsByPage(pageNum, limit, orderByURLParam)
	if err != nil {
		errResponse := dto.Error{
			Error:   "Unable to retrieve logs for page.",
			Message: err.Error(),
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, errResponse)

		return

	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, allClassificationReqs)
}
