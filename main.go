package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"shop.go/config"
	"shop.go/handlers"
	"shop.go/middlewares"
)

func main() {
	// 初始化設定
	config.LoadEnvFile()
	config.ConnectDB()

	// 創建 Gin 路由器
	router := gin.Default()

	// 使用 Session
	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", "", []byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		log.Println("Redis Store 初始化失敗", err)
		return
	}
	router.Use(sessions.Sessions("user_session", store))

	// 路由設定
	setUpWebRoutes(router)
	setUpAdminRoutes(router)

	// 啟動服務
	router.Run(":7777")
}

func setUpWebRoutes(router *gin.Engine) {
	router.GET("/set-cookie", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		session.Set("user_id", "123")
		session.Set("user_xx", "456")
		session.Set("user_dd", "789")
		session.Save()
		ctx.JSON(http.StatusOK, "pong")
	})

	router.GET("/read-cookie", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		userID := session.Get("user_id")
		ctx.JSON(http.StatusOK, userID)
	})

	router.GET("/update-cookie", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		// session.Options(sessions.Options{MaxAge: -1})
		session.Save()
		ctx.JSON(http.StatusOK, "cool")
	})
}

func setUpAdminRoutes(router *gin.Engine) {
	adminGroup := router.Group("/admin")
	{
		// 不需驗證
		adminGroup.POST("/login", handlers.AdminLogin)

		// 需要驗證
		authorizedGroup := adminGroup.Group("")
		authorizedGroup.Use(middlewares.AdminRequired)
		{
			authorizedGroup.POST("/logout", handlers.Logout)
			authorizedGroup.POST("/categories", handlers.AddCategory)
			authorizedGroup.GET("/categories", handlers.ListCategories)
		}
	}
}
