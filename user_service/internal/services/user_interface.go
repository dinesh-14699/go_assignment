package services

import "user_service/models"

type UserServiceInterface interface {
    RegisterUser(models.User) error
    LoginUser(username, password string) (string, error) 
    GetUserByID(userID uint) (*models.User, error)
}
