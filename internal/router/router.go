package router

import (
	"todolist/backend/internal/handler"
	"todolist/backend/internal/middleware"
	"todolist/backend/internal/repository"
	"todolist/backend/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userRepo := repository.NewUserRepository()
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(userRepo)

	authH := handler.NewAuthHandler(authSvc, userSvc)
	userH := handler.NewUserHandler(userSvc)

	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/refresh", authH.Refresh)
	}

	authorized := v1.Group("/")
	authorized.Use(middleware.Auth())
	{
		authorized.POST("/auth/logout", authH.Logout)

		authorized.GET("/users/me", userH.GetMe)
		authorized.PUT("/users/me", userH.UpdateMe)
		authorized.DELETE("/users/me", userH.DeleteMe)
	}

	return r
}
