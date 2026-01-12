package services

import (
	"PicSearch/app/db"
	"PicSearch/app/db/models"
	"errors"
)

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := db.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user by their unique ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string "error": "User not found"
// @Failure 500 {object} map[string]string "error": "Internal server error"
// @Router /users/{id} [get]
func GetUserByID(id int) (models.User, error) {
	var user models.User
	err := db.DB.Find(&user, id).Error
	if err != nil {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func CreateUser(user models.User) (models.User, error) {
	err := db.DB.Create(&user).Error
	if err != nil {
		return models.User{}, errors.New("could not create user")
	}
	return user, nil
}

func UpdateUser(id int, updatedData models.User) (models.User, error) {
	var user models.User
	err := db.DB.First(&user, id).Error
	if err != nil {
		return models.User{}, errors.New("user not found")
	}

	user.Name = updatedData.Name
	user.Email = updatedData.Email
	user.Password = updatedData.Password

	err = db.DB.Save(&user).Error
	if err != nil {
		return models.User{}, errors.New("could not update user")
	}
	return user, nil
}

func DeleteUser(id int) error {
	err := db.DB.Delete(&models.User{}, id).Error
	if err != nil {
		return errors.New("could not delete user")
	}
	return nil
}
