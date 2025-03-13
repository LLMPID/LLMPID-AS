package database

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"llm-promp-inj.api/config"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Connect initializes a database connection
func Connect(cfg *config.Config, log *logrus.Logger) (*gorm.DB, error) {
	var err error
	// Initiate databse connection only once and then use it for the lifecycle of the application.
	once.Do(func() {
		// Format the DB connection string.
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.Database.Host,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.Port,
		)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatal("Database connection failed:", err)
		}
		log.Println("Database connected successfully")
	})
	return db, err
}
