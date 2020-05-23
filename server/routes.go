package server

import (
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/services"
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/controllers"
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/utilities"
	"github.com/jinzhu/gorm"
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/utilities/security"
	"net/http"
	"github.com/gin-gonic/gin"
)

const (
	authServerURL = "http://localhost:8080/api/v1"
)

func (s *Server) initializeRoutes(configUtil utilities.IConfigUtil,db *gorm.DB) {

	v1 := s.Router.Group("/api/v1")
	{
		v1.GET("/", func(c *gin.Context) {
            http.ServeFile(c.Writer, c.Request, "index.html")
		})
		
		// User Routes
		userService := services.NewUserService(configUtil, db)
		usersController := controllers.NewUserController(userService)

		// OAuth Routes
		oauthService := services.NewOauthClientService(db, userService, configUtil)
		security := security.NewSecurity(db, oauthService)
		oauthController := controllers.NewOAuth2Controller(userService, oauthService)
		// v1.GET("/oauth/credentials",  oauthController.GetClientCridentials)
		v1.POST("/oauth/token",  oauthController.GetClientToken)
		v1.POST("/login",security.ValidClient, oauthController.Login)
		v1.POST("/register", security.ValidClient, oauthController.Register)
	}
}
