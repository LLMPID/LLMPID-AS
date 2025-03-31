package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"llm-promp-inj.api/config"
	"llm-promp-inj.api/internal/database"
	"llm-promp-inj.api/internal/handler"
	"llm-promp-inj.api/internal/log"
	"llm-promp-inj.api/internal/middleware"
	"llm-promp-inj.api/internal/pkg"
	"llm-promp-inj.api/internal/repository"
	"llm-promp-inj.api/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	// Instantiate logger and point it to log to a logfile
	log, err := log.NewLogger(cfg.Host.LogsDirPath, cfg.Host.Environment)
	if err != nil {
		fmt.Println("System logger setup failed.")
		panic(err)
	}
	log.Info("Inititated system logger.")

	// Open a database connection
	db, err := database.Connect(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	log.Info("Instantiate database connection.")

	// Dynamically generate a 32-byte server secret for token generation.
	serverSecretKey := generateSecureServerKey(32)

	// Instantiate repositories
	classificationLogsRepo := repository.NewClassificationLogsRepository(db, log)
	internalClassifierRepo := repository.NewInternalClassifierAPIRepository(cfg.Classifier.ClassifierAPIPath, log)
	userRepo := repository.NewUserRepository(db, log)
	tokenRepo := repository.NewTokenRepository(serverSecretKey, db, log)
	cryptoRepo := repository.NewCryptoRepository(log)
	sessionRepo := repository.NewSessionRepository(db, log)
	log.Info("Instantiate repositories.")

	// Instantiate services
	classficationService := service.NewClassificationService(classificationLogsRepo, internalClassifierRepo)
	tokenService := service.NewTokenService(tokenRepo)
	authService := service.NewAuthenticationService(userRepo, tokenRepo, cryptoRepo, sessionRepo)
	extSystemService := service.NewExternalSystemService(cryptoRepo, userRepo, log)

	log.Info("Instantiate services.")

	// Insatntiate middlewares
	authMiddleware := middleware.NewAuthMiddleware(tokenService, authService)

	// Instantiate handlers
	classificationHandler := handler.NewClassificationHandler(classficationService, authMiddleware)
	userHandler := handler.NewUserHandler(authService, authMiddleware)
	extSysHandler := handler.NewExternalSystemHandler(extSystemService, authService, authMiddleware)

	// Map handlers to routes
	// {handler_route}:{handler}
	handlers := map[string]pkg.Handler{
		"classification":  classificationHandler,
		"user":            userHandler,
		"system/external": extSysHandler,
		// Add more handlers
	}
	router := pkg.NewRouter(handlers, log)
	log.Info("Initiated handlers and router.")

	// Start server
	log.Printf("Running %s envrionment on port %s...\n", cfg.Host.Environment, cfg.Host.Port)
	hostStr := fmt.Sprintf(":%s", cfg.Host.Port)
	http.ListenAndServe(hostStr, router)
}

func generateSecureServerKey(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(bytes)
}
