package routes

import (
	"go_learning/two_week/day_10_mvc/controllers"
	"go_learning/two_week/day_10_mvc/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.POST("/register", controllers.Register(db))
	r.POST("/login", controllers.Login(db))

	// 受保护接口组
	authGroup := r.Group("/api")
	authGroup.Use(middleware.JWTAuth())
	{
		authGroup.GET("/profile", controllers.GetProfile)
	}

	return r
}
