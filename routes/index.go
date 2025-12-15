package routes

import (
	"github.com/gin-gonic/gin"
	"shop.go/handlers"
)

func Setup(router *gin.Engine) {
	api := router.Group("/api")

	api.POST("/:role/signup", handlers.Signup)
}
