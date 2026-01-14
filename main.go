package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"shop.go/boot"
	"shop.go/middleware"
	"shop.go/routes"
)

func main() {
	// 初始化設定
	boot.LoadEnvFile()
	boot.ConnectDB()
	boot.ConnectStorage()

	// 創建 Gin 路由器
	router := gin.Default()

	// CORS 設定
	router.Use(middleware.CORS())

	// 路由設定
	routes.Setup(router)


	// 啟動服務
	router.Run(":" + os.Getenv("APP_PORT"))
}
