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
		r.logger.Error("Unable to insert user objeect into the database. ERR: ", err.Error())
		return user, errors.New("unable to insert object")
	}

	return user, nil
}

func (r *UserRepository) UpdatePasswordHashByUserID(id uint, passwordHash string) error {
	err := r.DB.Model(&models.User{}).Where("id = ?", id).Update("password_hash", passwordHash).Error
	if err != nil {
		r.logger.Error("Unable to update user's password hash. ERR: ", err.Error())
		return errors.New("unalbe to update password")
	}

	return nil
}

func (r *UserRepository) UpdateUsername(oldUsername string, newUsername string) (models.User, error) {
	var updatedUser models.User
	updateEvent := r.DB.Model(&models.User{}).Where("username = ?", oldUsername).Update("username", newUsername)

	if updateEvent.Error != nil {
		r.logger.Error("Unable to update user's username. ERR: ", updateEvent.Error.Error())
		return models.User{}, errors.New("unalbe to update name")
	}

	updateEvent.Scan(updatedUser)
	return updatedUser, nil
}
func (r *UserRepository) SelectUserByUsername(username string) (models.User, error) {
	var foundUser models.User
	if err := r.DB.Where("username = ?", username).First(&foundUser).Error; err != nil {
		return foundUser, err
	}

	return foundUser, nil
}
func (r *UserRepository) SelectUserByRole(role string) ([]models.User, error) {
	var users []models.User

	if err := r.DB.Where("role = ?", role).Find(&users).Error; err != nil {
		r.logger.Error("Failed to retrieve user by role. ERR: ", err.Error())
		return users, err
	}

	return users, nil

}

func (r *UserRepository) DeleteByUsername(username string) error {
	if err := r.DB.Where("username = ?", username).Delete(&models.User{}).Error; err != nil {
		r.logger.Error("Failed to delete user from database. ERR: ", err.Error())
		return errors.New("unable to delete object")
	}

	return nil
}
