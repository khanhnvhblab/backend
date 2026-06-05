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

	todoRepo := repository.NewTodoRepository()
	todoSvc := service.NewTodoService(todoRepo)

	categoryRepo := repository.NewCategoryRepository()
	categorySvc := service.NewCategoryService(categoryRepo)

	authH := handler.NewAuthHandler(authSvc, userSvc)
	userH := handler.NewUserHandler(userSvc)
	todoH := handler.NewTodoHandler(todoSvc)
	categoryH := handler.NewCategoryHandler(categorySvc)

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

		authorized.GET("/todos", todoH.List)
		authorized.POST("/todos", todoH.Create)
		authorized.GET("/todos/:id", todoH.GetByID)
		authorized.PUT("/todos/:id", todoH.Update)
		authorized.PATCH("/todos/:id/status", todoH.UpdateStatus)
		authorized.DELETE("/todos/:id", todoH.Delete)

		authorized.GET("/categories", categoryH.List)
		authorized.POST("/categories", categoryH.Create)
		authorized.GET("/categories/:id", categoryH.GetByID)
		authorized.PUT("/categories/:id", categoryH.Update)
		authorized.DELETE("/categories/:id", categoryH.Delete)
	}

	return r
}
