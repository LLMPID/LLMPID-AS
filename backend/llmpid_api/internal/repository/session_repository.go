package repository

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"llm-promp-inj.api/internal/models"
)

type SessionRepository struct {
	DB     *gorm.DB
	logger *logrus.Logger
}

func NewSessionRepository(db *gorm.DB, logger *logrus.Logger) *SessionRepository {
	return &SessionRepository{DB: db, logger: logger}
}

func (r *SessionRepository) CreateSession(sessionSlug string, sub string) error {
	if r.IsValidSession(sessionSlug, sub) {
		return errors.New("user already has a session")
	}

	session := models.Session{Slug: sessionSlug, Sub: sub}

	err := r.DB.Create(&session).Error
	if err != nil {
		r.logger.Error("Unable to create new user session. ERR: ", err)
		return errors.New("unable to create session")
	}

	return nil
}

func (r *SessionRepository) IsValidSession(sessionSlug string, sub string) bool {
	session, _ := r.SelectSessionBySlugAndSub(sessionSlug, sub)

	return models.Session{} == session
}

func (r *SessionRepository) DeleteSessionBySub(sub string) error {
	err := r.DB.Where("sub = ?", sub).Delete(&models.Session{})
	if err != nil {
		return errors.New("failed to remove session")
	}
	return nil
}

func (r *SessionRepository) SelectSessionBySub(sub string) (models.Session, error) {
	session := models.Session{}

	err := r.DB.Where("slug = ?", sub).First(&session).Error
	if err != nil {
		return session, err
	}

	return session, nil
}

func (r *SessionRepository) SelectSessionBySlugAndSub(slug string, sub string) (models.Session, error) {
	session := models.Session{}

	err := r.DB.Where("slug = ? AND sub = ?", slug, sub).First(&session).Error
	if err != nil {
		return session, err
	}

	return session, nil
}
