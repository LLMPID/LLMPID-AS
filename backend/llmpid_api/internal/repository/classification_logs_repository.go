package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"llm-promp-inj.api/internal/models"
)

type ClassificationLogsRepository struct {
	DB     *gorm.DB
	logger *logrus.Logger
}

func NewClassificationLogsRepository(db *gorm.DB, logger *logrus.Logger) *ClassificationLogsRepository {
	return &ClassificationLogsRepository{DB: db, logger: logger}
}

// InsertClassificationRequest inserts a classification  log (ClassificationLog) into the database.
func (r *ClassificationLogsRepository) InsertClassificationLog(classificationLog models.ClassificationLog) error {
	return r.DB.Create(&classificationLog).Error
}

// SelectClassificationLogByID returns a single database entry for a classification log based on ID.
func (r *ClassificationLogsRepository) SelectClassificationLogByID(id uint) (models.ClassificationLog, error) {
	var classificationLog models.ClassificationLog
	if err := r.DB.First(&classificationLog, id).Error; err != nil {
		return classificationLog, err
	}

	return classificationLog, nil
}

// SelectClassificationLogsByPage retrieves database entries of classification logs, based on page and limit for offsetting.
func (r *ClassificationLogsRepository) SelectClassificationLogsByPage(page int, limit int, orderBy string) ([]models.ClassificationLog, error) {
	var classificationLogs []models.ClassificationLog

	// Essentialy, `SELECT * FROM classification_logs ORDER BY id {desc || asc} LIMIT {limit} OFFSET {(page-1)*limit};`.
	if err := r.DB.Offset((page - 1) * limit).Limit(limit).Order(orderBy).Find(&classificationLogs).Error; err != nil {
		r.logger.Error("Failed to retrieve classification log from database.")
		return nil, err
	}

	return classificationLogs, nil
}
