package services

import (
	//"time"
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/models"
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/utilities"
	"github.com/jinzhu/gorm"
	"encoding/json"
)

type SessionService struct {
	db *gorm.DB
	configUtil     utilities.IConfigUtil
}

type Data struct {
	UserID string
	ClientID string
	Access string

}

func NewSessionService(configUtil utilities.IConfigUtil, db *gorm.DB) *SessionService {
	service := SessionService{}
	service.db = db
	service.configUtil = configUtil
	return &service
}

func (service SessionService) Find(accessToken string) (error, *Data) {
	var err error
	var session = &models.Oauth2Token{}
	err = service.db.Debug().Model(&models.Oauth2Token{}).Where("access = ?", accessToken).Take(&session).Error
	if err != nil {
		return err, nil
	}
	data := &Data{}
	err = json.Unmarshal([]byte(session.Data), data)
	if err != nil {
    	return err, nil
	}
	return nil, data
}
