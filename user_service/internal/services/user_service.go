package services

import (
	"errors"
	"time"
	"user_service/config"
	"user_service/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
    DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{DB: db}
}

func (s *UserService) RegisterUser(user models.User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user.Password = string(hashedPassword)
    user.CreatedAt = time.Now()

    return s.DB.Create(&user).Error
}

func (s *UserService) LoginUser(username, password string) (string, error) {
    var user models.User
    if err := s.DB.Where("email = ?", username).First(&user).Error; err != nil {
        return "", errors.New("invalid username or password")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", errors.New("invalid username or password")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "username": user.Username,
        "email": user.Email,
        "exp":     time.Now().Add(time.Minute * 2).Unix(), 
    })

    tokenString, err := token.SignedString([]byte(config.JWTSecret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
    var user models.User
    if err := s.DB.First(&user, userID).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
