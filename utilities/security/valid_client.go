package security

import (
	"net/http"

	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/services"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	//"fmt"
)


type Security struct {
	db *gorm.DB
	oauthService services.IOauthClientService
}

func NewSecurity(db *gorm.DB, oauthService services.IOauthClientService) *Security{
	security := Security{}
	security.db = db
	security.oauthService = oauthService
	return &security
}

func (security Security) ValidClient(c *gin.Context)  {
	server, err := security.oauthService.GetOauth2Server("client_credentials")
	_, err = server.ValidationBearerToken(c.Request)
	   if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthoirzed client"})
		c.Abort()
		return
	   }
	c.Next()
}


func (security Security) ValidUser(c *gin.Context)  {
	server, err := security.oauthService.GetOauth2Server("password")
	_, err = server.ValidationBearerToken(c.Request)
	   if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthoirzed user"})
		c.Abort()
		return
	   }
	c.Next()
}