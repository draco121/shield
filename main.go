package main

import (
	"github.com/draco121/horizon/utils"
	"os"

	"shield/controllers"
	"shield/core"
	"shield/repository"
	"shield/routes"

	"github.com/draco121/horizon/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func RunApp() {
	utils.Logger.Info("Starting authentication service")
	client := database.NewMongoDatabase(os.Getenv("MONGODB_URI"))
	db := client.Database("authentication-service")
	authRepo := repository.NewAuthenticationRepository(db)
	userRepo := repository.NewUserRepository(db)
	userService := core.NewUserService(client, userRepo)
	authService := core.NewAuthenticationService(client, authRepo, userRepo)
	controller := controllers.NewControllers(authService, userService)
	router := gin.New()
	router.Use(gin.LoggerWithWriter(utils.Logger.Out))
	routes.RegisterRoutes(controller, router)
	err := router.Run()
	utils.Logger.Info("authentication service started successfully")
	if err != nil {
		utils.Logger.Fatal(err)
		return
	}
}
func main() {
	_ = godotenv.Load()
	RunApp()
}
