package routes

import (
	"github.com/draco121/authenticationservice/controllers"
	"github.com/draco121/common/constants"
	"github.com/draco121/common/middlewares"
	"github.com/draco121/common/utils"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(controllers controllers.Controllers, router *gin.Engine) {
	utils.Logger.Info("Registering routes...")
	v1 := router.Group("/v1")
	v1.POST("/login", controllers.Login)
	v1.POST("/refresh", controllers.RefreshLogin)
	v1.POST("/logout", controllers.Logout)
	v1.POST("/user", controllers.CreateUser)
	v1.GET("/user", middlewares.AuthMiddleware(constants.Write), controllers.GetUserProfile)
	v1.PATCH("/user", middlewares.AuthMiddleware(constants.Write), controllers.UpdateUser)
	v1.DELETE("/user", middlewares.AuthMiddleware(constants.All), controllers.DeleteUser)
	utils.Logger.Info("Registered routes...")
}
