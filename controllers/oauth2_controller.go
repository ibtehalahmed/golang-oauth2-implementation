package controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/viewModels"
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/services"
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/models"

	"gopkg.in/oauth2.v3/store"
	"gopkg.in/oauth2.v3/server"
	"golang.org/x/oauth2"
	"context"
	"net/http"
)

type OAuth2Controller struct {
	config *oauth2.Config
	cs *store.ClientStore
	server *server.Server
	userService services.IUserService
	oauthClientService services.IOauthClientService
}

func NewOAuth2Controller(userService services.IUserService, oauthClientService services.IOauthClientService) *OAuth2Controller {
	controller := OAuth2Controller{}
	controller.userService = userService
	controller.oauthClientService = oauthClientService
	return &controller
}
type Client struct{
	ClientSecret string `form:"client_secret"`
	GrantType string `form:"grant_type"`
	ClientId string `form:"client_id"`
	Username string `form:"username"`
	Password string `form:"password"`
}

func (controller OAuth2Controller) GetClientToken(c *gin.Context) {
	var cl Client
	err := c.Bind(&cl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	server, err := controller.oauthClientService.GetOauth2Server(cl.GrantType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	controller.server = server
	err = controller.server.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (controller OAuth2Controller) Login(c *gin.Context) {
	var loginViewModel viewModels.LoginViewModel
	err := c.Bind(&loginViewModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.config = services.Config
	token, err := controller.config.PasswordCredentialsToken(context.Background(), loginViewModel.Username, loginViewModel.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, token)
}

func (controller OAuth2Controller) Register(c *gin.Context) {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	credentials := user
	controller.userService.Prepare(&user)
	err = controller.userService.Validate(&user, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = controller.userService.BeforeSave(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = controller.userService.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	controller.config = services.Config
	token, err := controller.config.PasswordCredentialsToken(context.Background(), credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, token)
}
