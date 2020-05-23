package main

import (
	"fmt"
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/utilities"
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/models"
)

var users = []models.User{
	models.User{
		Username: "ibtehal",
		Name: "ibtehal",
		Surname: "ibrahim",
		Email:    "ibtehal@mail.com",
		Password: "$2a$10$OfRrzXjoQOQ2/1S68q3dLOKmKv6t09bhwgk.xlI1UJ1fhf8RZgaMO",
		Phone: "1234567891",
	},
}

type Oauth2Token struct {}

var oauth2_clients = []models.Oauth2Client{
	models.Oauth2Client{
		ID:1,
		Name: "web-password",
		Redirect: "http://localhost",
		PersonalAccessClient: "0",
		PasswordClient: "1",
		Revoked: "0",
		Domain: "",
		Secret: "",
	},
	models.Oauth2Client{
		ID:2,
		Name: "web-client",
		Redirect: "",
		PersonalAccessClient: "0",
		PasswordClient: "0",
		Revoked: "0",
		Domain: "",
		Secret: "",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Post{}, 
		&models.User{},
		&Oauth2Token{},
		&models.Oauth2Client{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(
		&models.User{},
		&models.Oauth2Client{},).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}


	for i, _ := range oauth2_clients {
		err = db.Debug().Model(&models.Oauth2Client{}).Create(&oauth2_clients[i]).Error
		if err != nil {
			log.Fatalf("cannot seed oauth2_clients table: %v", err)
		}
	}


	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

}

func main() {
	var err error
	configUtil := utilities.NewConfigUtil()

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", configUtil.GetConfig("dbHost"), configUtil.GetConfig("dbPort"), configUtil.GetConfig("dbUser"), configUtil.GetConfig("dbName"), configUtil.GetConfig("dbPassword"))
	DB, err := gorm.Open(configUtil.GetConfig("dbDriver"), DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", configUtil.GetConfig("dbDriver"))
		log.Fatal(err)
	} else {
		fmt.Printf("We are connected to the %s database", configUtil.GetConfig("dbDriver"))
	}
	//database seed
	Load(DB);
}