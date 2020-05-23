package services

import (
	dbModels "gitlab.com/ibtehalahmed/golang-oauth2-implementation/models"
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/utilities"
	"log"
	"fmt"
	"github.com/jinzhu/gorm"

	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"golang.org/x/oauth2"
	"gopkg.in/oauth2.v3/models"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/oauth2.v3/generates"
	"strconv"
	"gopkg.in/oauth2.v3/errors"
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	pg "github.com/vgarvardt/go-oauth2-pg"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
)
var Config *oauth2.Config
type IOauthClientService interface {
	GetOauth2Server(grantType string) (*server.Server, error)
}

type OauthClientService struct {
	db *gorm.DB
	config *oauth2.Config
	clientStore *store.ClientStore
	manager *manage.Manager
	server *server.Server
	userService IUserService
	configUtil utilities.IConfigUtil
}

func NewOauthClientService( db *gorm.DB, userService IUserService, configUtil utilities.IConfigUtil ) *OauthClientService {
	service := OauthClientService{}
	service.userService = userService
	service.configUtil = configUtil
	service.clientStore = store.NewClientStore()
	service.manager = manage.NewDefaultManager()

	service.manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	DBURL := configUtil.GetConfig("PG_URI")
	pgxConn, err := pgx.Connect(context.TODO(), DBURL)
	if (err != nil) {
		fmt.Println("=============>", err)
	}

	// use PostgreSQL token store with pgx.Connection adapter
	adapter := pgx4adapter.NewConn(pgxConn)
	tokenStore, _ := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	defer tokenStore.Close()
	//service.clientStore, _ = pg.NewClientStore(adapter)
	service.manager.MapTokenStorage(tokenStore)
	service.manager.MapClientStorage(service.clientStore)


	//service.manager.MustTokenStorage(tokenStore)
	service.manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("00000000"), jwt.SigningMethodHS512))
	service.server = server.NewServer(server.NewConfig(), service.manager)
	service.server.SetAllowGetAccessRequest(true)
	service.server.SetClientInfoHandler(server.ClientFormHandler)
	service.manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)
	service.server.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	 })
	 service.server.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	 })
	service.config = &oauth2.Config{}
	service.db = db
	return &service
}

const (
	authServerURL = "http://localhost:8080"
)

func (service OauthClientService) GetOauth2Server(grantType string) (*server.Server, error) {
	var err error
	var client = &dbModels.Oauth2Client{}
	if (grantType == "password") {
		err = service.db.Debug().Model(&dbModels.Oauth2Client{}).Where("password_client = ?", "1").Take(&client).Error
	} else {
		err = service.db.Debug().Model(&dbModels.Oauth2Client{}).Where("password_client = ?", "0").Take(&client).Error
	}
	if err != nil {
		return service.server, err
	}
	Config = &oauth2.Config{
		ClientID:  strconv.FormatUint(client.ID, 10),
		ClientSecret: client.Secret,
		Scopes:       []string{""},
		RedirectURL:  "localhost:8080",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/authorize",
			TokenURL: authServerURL + "/api/v1/oauth/token",
		},
	}
	err = service.clientStore.Set(strconv.FormatUint(client.ID, 10), &models.Client{
		ID:     strconv.FormatUint(client.ID, 10),
		Secret: client.Secret,
		Domain: "http://localhost:8080",
	})
	if(err!=nil){
		return service.server, err
	}

	service.manager.MapClientStorage(service.clientStore)
	service.server.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		err, userID = service.passwordAuthorizationHandler(username, password)
		if(err != nil) {
			return "", err
		}
		return userID, err
	})
	return service.server, err
}

func (service OauthClientService) passwordAuthorizationHandler(username string, password string) (error, string) {
	var (
		user = &dbModels.User {
			Username: username,
			Password: password,
		}
	)
	err := service.userService.Validate(user, "login")
	if err != nil {
		return err, ""
	}
	userID, err := service.SignIn(username, password)
	if err != nil {
		return err, ""
	}
	return nil, userID
}

func (service OauthClientService) SignIn(username, password string) (string, error) {
	var err error
	user := dbModels.User{}
	err = service.db.Debug().Model(dbModels.User{}).Where("username = ?", username).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = service.userService.VerifyPassword(user.Password, password)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(user.ID, 10), nil
}
