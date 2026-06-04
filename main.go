// @title           TodoList API
// @version         1.0
// @description     REST API for TodoList — manage todos, categories, and user accounts.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Enter: Bearer {access_token}

package main

import (
	"log"
	_ "todolist/backend/docs"
	"todolist/backend/config"
	"todolist/backend/internal/db"
	"todolist/backend/internal/router"
)

func main() {
	config.Load()

	db.Connect()
	defer db.Disconnect()

	r := router.New()

	addr := ":" + config.App.AppPort
	log.Printf("Server running on %s (env: %s)", addr, config.App.AppEnv)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
