package api

import (
	"cadastral-service/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/go-delve/delve/pkg/config"
)

func SetupRoutes(router *gin.Engine, handler *Handler, cfg *config.Config){
	//API group v1
	v1 := router.Group("/api/v1")

	//public endpoints
	v1.GET("/ping", handler.Ping)

	// protected endpoints with authorization if it turn on
	if cfg.Auth.Enabled {
		v1.POST("/login", handler.Login)
		v1.POST("/register", handler.Register)
		
		//use middleware auth
		authGroup := v1.Group("/")
		authGroup.Use(handler.AuthMiddleware())
		{
			authGroup.POST("/query", handler.CreateQuery)
			authGroup.GET("/history", handler.GetHistory)
			authGroup.GET("/history/:cadastral_number", handler.GetHistoryByCadastral)
		}
	} else {
		//without auth
		v1.POST("/query", handler.CreateQuery)
		v1.GET("/history", handler.GetHistory)
		v1.GET("/history/:cadastral_number", handler.GetHistoryByCadastral)
	}
	//endpoint for external server emulation
	router.POST("/api/result", handler.ProcessResult)

	//swagger doc (optional)
	if cfg.DocsEnabled {
		router.GET("/swagger/*any", handler.SwaggerHandler)
	}
}