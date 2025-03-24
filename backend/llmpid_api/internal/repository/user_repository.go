package repository

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"llm-promp-inj.api/internal/models"
)

type UserRepository struct {
	DB     *gorm.DB
	logger *logrus.Logger
}

func NewUserRepository(db *gorm.DB, logger *logrus.Logger) *UserRepository {
	return &UserRepository{DB: db, logger: logger}
}

func (r *UserRepository) InsertUser(username string, passwordHash string, role string) (models.User, error) {
	user := models.User{Username: username, PasswordHash: passwordHash, Role: role}

	err := r.DB.Create(&user).Error
	if err != nil {
		r.logger.Error("Unable to insert user objeect into the database. ERR: ", err)
		return user, errors.New("unable to insert object")
	}

	return user, nil
}

func (r *UserRepository) SelectUserByUsername(username string) (models.User, error) {
	var foundUser models.User
	if err := r.DB.Where("username = ?", username).First(&foundUser).Error; err != nil {
		return foundUser, err
	}

	return foundUser, nil
}
