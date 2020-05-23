package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/services"
	"strconv"
	"errors"
)

type UserController struct {
	userService services.IUserService
}



func NewUserController(userService services.IUserService) *UserController {
	controller := UserController{}
	controller.userService = userService
	return &controller
}

func (controller UserController) Root(c *gin.Context) {
	c.JSON(http.StatusOK, true)
}

func (controller UserController) CreateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"User Created": true})
}
