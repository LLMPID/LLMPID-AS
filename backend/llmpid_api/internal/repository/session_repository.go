package repository

import (
	"errors"
	"time"

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

func (r *SessionRepository) CreateSession(sessionID string, sub string, expiresAt int64) error {
	session := models.Session{SessionID: sessionID, Sub: sub, ExpiresAt: expiresAt}

	err := r.DB.Create(&session).Error
	if err != nil {
		r.logger.Error("Unable to create new user session. ERR: ", err)
		return errors.New("unable to create session")
	}

	return nil
}

func (r *SessionRepository) IsValidSession(sessionID string, sub string) bool {
	session, _ := r.SelectSessionBySIDAndSub(sessionID, sub)
	if (models.Session{} == session) {
		return false
	}

	if time.Now().Unix() > session.ExpiresAt {
		r.DeleteSessionBySub(sessionID)
		return false
	}

	return true
}

func (r *SessionRepository) DeleteSessionBySID(sid string) {
	r.DB.Where("session_id = ?", sid).Delete(&models.Session{})
}

func (r *SessionRepository) DeleteSessionBySub(sub string) {
	r.DB.Where("sub = ?", sub).Delete(&models.Session{})
}

func (r *SessionRepository) SelectSessionBySub(sub string) (models.Session, error) {
	var session models.Session

	err := r.DB.Where("sub = ?", sub).First(&session).Error
	if err != nil {
		return models.Session{}, err
	}

	return session, nil
}

func (r *SessionRepository) SelectSessionBySIDAndSub(sid string, sub string) (models.Session, error) {
	var session models.Session

	err := r.DB.Where("session_id = ?", sid).Where("sub = ?", sub).First(&session).Error
	if err != nil {
		r.logger.Error("Unable to find session by session_id and sub. ERR: ", err)
		return session, err
	}

	return session, nil
}
