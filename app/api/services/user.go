package services

import (
	"PicSearch/app/db"
	"PicSearch/app/db/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type LoginUserResponse struct {
	models.User
	Type  string `json:"type"`
	Token string `json:"token"`
}

func LoginUser(email, password string) (LoginUserResponse, error) {
	var user models.User
	err := db.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return LoginUserResponse{}, errors.New("No user found with the provided email")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return LoginUserResponse{}, errors.New("Incorrect password")
	}

	signedToken, err := GenerateToken(user.ID)

	return LoginUserResponse{
		User:  user,
		Type:  "Bearer",
		Token: signedToken,
	}, nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := db.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

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
