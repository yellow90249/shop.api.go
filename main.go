package main

import (
	"log"
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
	setUpPublicRoutes(router)
	setUpWebRoutes(router)
	setUpAdminRoutes(router)

	// 啟動服務
	router.Run(":7777")
}

func setUpPublicRoutes(router *gin.Engine) {
	router.Static("/uploads", "./uploads")
}

func setUpWebRoutes(router *gin.Engine) {
	router.GET("/categories", handlers.ListCategories)
	router.GET("/products", handlers.ListProducts)
	router.GET("/products/:productId", handlers.GetProduct)
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

			// 種類
			authorizedGroup.POST("/categories", handlers.AddCategory)
			authorizedGroup.GET("/categories", handlers.ListCategories)
			authorizedGroup.PUT("/categories/:categoryId", handlers.UpdateCategory)
			authorizedGroup.DELETE("/categories/:categoryId", handlers.DeleteCategory)

			// 商品
			authorizedGroup.POST("/products", handlers.AddProduct)
			authorizedGroup.GET("/products", handlers.ListProducts)
			authorizedGroup.PUT("/products/:productId", handlers.UpdateProduct)
			authorizedGroup.PUT("/products/:productId/image", handlers.UpdateProductImage)
			authorizedGroup.DELETE("/products/:productId", handlers.DeleteProduct)

			// 購物車
			authorizedGroup.POST("/carts", handlers.AddCart)
		}
	}
}
