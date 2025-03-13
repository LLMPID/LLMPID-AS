package pkg

import (
	"net/http"
	"time"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"llm-promp-inj.api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	Routes() chi.Router
}

func NewRouter(handlers map[string]Handler, logger *logrus.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.LogrusLogger(logger))         // In-depth Logrus logging for each request.
	router.Use(chiMiddleware.Recoverer)                 // Prevents crashes on panics.
	router.Use(chiMiddleware.Timeout(60 * time.Second)) // Prevents slow requests from blocking the API.
	router.Use(middleware.XSSHandler)                   // Makes sure that the "request_text" field in returned classification logs does not contain valid HTML and JS.

	router.Route("/api", func(api chi.Router) {
		for path, handler := range handlers {
			api.Mount("/"+path, handler.Routes())
		}
	})

	// Health check route
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	return router
}
