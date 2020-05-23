package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/utilities"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

var errList = make(map[string]string)

func (server *Server) Initialize(configUtil utilities.IConfigUtil) {
	var err error

	// If using mysql
	if configUtil.GetConfig("dbDriver") == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", 
		configUtil.GetConfig("dbUser"),
		configUtil.GetConfig("dbPassword"), configUtil.GetConfig("dbHost"), configUtil.GetConfig("dbPort"), configUtil.GetConfig("dbName"))
		server.DB, err = gorm.Open(configUtil.GetConfig("dbDriver"), DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", configUtil.GetConfig("dbDriver"))
			log.Fatal(err)
		} else {
			fmt.Printf("We are connected to the %s database", configUtil.GetConfig("dbDriver"))
		}
	} else if configUtil.GetConfig("dbDriver") == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", configUtil.GetConfig("dbHost"), configUtil.GetConfig("dbPort"), configUtil.GetConfig("dbUser"), configUtil.GetConfig("dbName"), configUtil.GetConfig("dbPassword"))
		server.DB, err = gorm.Open(configUtil.GetConfig("dbDriver"), DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", configUtil.GetConfig("dbDriver"))
			log.Fatal(err)
		} else {
			fmt.Printf("We are connected to the %s database", configUtil.GetConfig("dbDriver"))
		}
	} else {
		fmt.Println("Unknown Driver")
	}

	server.Router = gin.Default()

	// Apply the middleware to the router (works with groups too)
	server.Router.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, PUT, POST, DELETE",
		RequestHeaders: "Origin, Authorization, Content-Type",
		ExposedHeaders: "",
		MaxAge: 300 * time.Second,
		Credentials: false,
		ValidateHeaders: false,
	}))

	server.initializeRoutes(configUtil, server.DB)

}

func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
