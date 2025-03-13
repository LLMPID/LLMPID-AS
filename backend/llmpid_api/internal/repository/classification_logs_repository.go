package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"llm-promp-inj.api/internal/models"
)

type ClassificationLogRepository struct {
	DB     *gorm.DB
	logger *logrus.Logger
}

func NewClassificationLogRepository(db *gorm.DB, logger *logrus.Logger) *ClassificationLogRepository {
	return &ClassificationLogRepository{DB: db, logger: logger}
}

// InsertClassificationRequest inserts a classification request log (ClassificationLog) into the database.
func (r *ClassificationLogRepository) InsertClassificationLog(ClassificationLog *models.ClassificationLog) error {
	return r.DB.Create(ClassificationLog).Error
}

// SelectClassificationLogByID returns a single database entry for a classification request based on ID.
func (r *ClassificationLogRepository) SelectClassificationLogByID(id uint) (*models.ClassificationLog, error) {
	var ClassificationLog models.ClassificationLog
	if err := r.DB.First(&ClassificationLog, id).Error; err != nil {
		return nil, err
	}

	return &ClassificationLog, nil
}

// SelectClassificationLogsByPage retrieves database entries of classification requests, based on page and limit for offsetting.
func (r *ClassificationLogRepository) SelectClassificationLogsByPage(page int, limit int, orderBy string) (*[]models.ClassificationLog, error) {
	var ClassificationLogs *[]models.ClassificationLog

	// Essentialy, `SELECT * FROM classification_logs ORDER BY id {desc || asc} LIMIT {limit} OFFSET {(page-1)*limit};`.
	if err := r.DB.Offset((page - 1) * limit).Limit(limit).Order(orderBy).Find(&ClassificationLogs).Error; err != nil {
		r.logger.Error("Failed to retrieve classification requests from database.")
		return nil, err
	}

	return ClassificationLogs, nil
}
