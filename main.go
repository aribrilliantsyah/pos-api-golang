package main

import (
	"log"
	"pos-api/app/server"
	"pos-api/util/config"
	"pos-api/util/dbmigrate"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	if err := dbmigrate.RunMigrations(config); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	appServer := server.NewServer(config)
	err = appServer.Run()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
