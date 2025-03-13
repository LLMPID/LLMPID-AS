package main

import (
	"fmt"
	"net/http"

	"llm-promp-inj.api/config"
	"llm-promp-inj.api/internal/database"
	"llm-promp-inj.api/internal/handler"
	"llm-promp-inj.api/internal/log"
	"llm-promp-inj.api/internal/pkg"
	"llm-promp-inj.api/internal/repository"
	"llm-promp-inj.api/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	log, err := log.NewLogger(cfg.Host.LogsDirPath, cfg.Host.Environment)
	if err != nil {
		fmt.Println("System logger setup failed.")
		panic(err)
	}
	log.Info("Inititated system logger.")

	db, err := database.Connect(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	log.Info("Initiated database connection.")

	// Instantiate repositories
	classificationLogsRepo := repository.NewClassificationLogRepository(db, log)
	internalClassifierRepo := repository.NewInternalClassifierAPIRepository(cfg.Classifier.ClassifierAPIPath, log)
	log.Info("Initiated repositories.")

	// Instantiate services
	classificationReqService := service.NewClassificationRequestService(classificationLogsRepo, internalClassifierRepo)
	log.Info("Initiated services.")

	// Instantiate handlers
	classificationReqHandler := handler.NewClassificationRequestHandler(classificationReqService)

	// Map handlers to routes
	// {handler_route}:{handler}
	handlers := map[string]pkg.Handler{
		"classification": classificationReqHandler,
		// Add more handlers
	}
	router := pkg.NewRouter(handlers, log)
	log.Info("Initiated handlers and router.")

	// Start server
	log.Printf("Running %s envrionment on port %s...\n", cfg.Host.Environment, cfg.Host.Port)
	hostStr := fmt.Sprintf(":%s", cfg.Host.Port)
	http.ListenAndServe(hostStr, router)
}
