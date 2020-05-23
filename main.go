package main

import (
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/utilities"
	"gitlab.com/ibtehalahmed/golang-oauth2-implementation/server"
	"log"
	"fmt"
)

func GetPort() string {
	confiUtil := utilities.NewConfigUtil()
	port := confiUtil.GetConfig("port")
	if port == "" {
		port = "2020"
		log.Println("[-] No PORT environment variable detected. Setting to ", port)
	}
	return ":" + port
}

func main() {

	server := server.Server{}
	configUtil := utilities.NewConfigUtil()
	server.Initialize(configUtil)

	apiPort :=  GetPort()
	fmt.Printf("Listening to port %s", apiPort)
	server.Run(apiPort)
}