package services

import (
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/models"
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/utilities"
	"strings"
	"time"
	"errors"
	"html"
	"unicode"
	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/gorm"
)

type IUserService interface {
	CurrentUser(accessToken string) (error, *models.User)
	SaveUser(user *models.User) (error)
	BeforeSave(u *models.User) error
	Prepare(user *models.User)
	FindByUsername(username string) (*models.User, error)
	FindByID(id string) (*models.User, error)
}

type UserService struct {
	db *gorm.DB
	configUtil     utilities.IConfigUtil
	sessionService *SessionService
}

func NewUserService(configUtil utilities.IConfigUtil, db *gorm.DB) *UserService {
	service := UserService{}
	service.db = db
	service.configUtil = configUtil
	service.sessionService = NewSessionService(configUtil, db)

	return &service
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (service UserService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (service UserService) BeforeSave(u *models.User)error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (service UserService) Prepare(user *models.User) {
	user.ID = 0
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

func (service UserService) SaveUser(user *models.User) (error) {
	var err error
	err = service.db.Debug().Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}
func (service UserService) CurrentUser(accessToken string) (error, *models.User) {

	err , userSession:= service.sessionService.Find(accessToken)
	if err != nil {
		return err, nil
	}
	user, err := service.FindByID(userSession.UserID)
	return err, user

}

func (service UserService) FindByID(id string) (*models.User, error) {
	var err error
	var user = &models.User{}
	err = service.db.Debug().Model(&models.User{}).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (service UserService) FindByUsername(username string) (*models.User, error) {
	var err error
	var user = &models.User{}
	err = service.db.Debug().Model(&models.User{}).Where("username = ?", username).Take(&user).Error
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}
