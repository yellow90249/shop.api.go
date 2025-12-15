package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"shop.go/config"
	"shop.go/routes"
)

func main() {
	// 初始化設定
	config.LoadEnvFile()
	config.ConnectDB()
	config.ConnectStorage()

	// 創建 Gin 路由器
	router := gin.Default()

	// 大家都進來吧（開發用）
	// if os.Getenv("GO_ENV") == "development" {
	// 	router.Use(cors.New(cors.Config{
	// 		AllowAllOrigins:  true,
	// 		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	// 		AllowHeaders:     []string{"*"},
	// 		ExposeHeaders:    []string{"Content-Length"},
	// 		AllowCredentials: false,
	// 		MaxAge:           12 * time.Hour,
	// 	}))
	// }

	// 路由設定
	routes.Setup(router)

	// 啟動服務
	router.Run(":" + os.Getenv("APP_PORT"))
}
