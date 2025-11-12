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
	// config.Migrate()

	// 創建 Gin 路由器
	router := gin.Default()
	// router.POST("/signup", handlers.Signup)

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
	// 不需 Customer 登入
	router.POST("/login", handlers.Login)
	router.GET("/categories", handlers.ListCategories)
	router.GET("/products", handlers.ListProducts)
	router.GET("/products/:productId", handlers.GetProduct)

	// 需要 Customer 登入
	router.GET("/me", middlewares.CustomerRequired, handlers.GetUser)
	router.POST("/logout", middlewares.CustomerRequired, handlers.Logout)
	router.POST("/cart/items", middlewares.CustomerRequired, handlers.AddCartItem)
	router.PUT("/cart/items/:cartItemId", middlewares.CustomerRequired, handlers.UpdateCartItemQuantity)
	router.DELETE("/cart/items/:cartItemId", middlewares.CustomerRequired, handlers.DeleteCartItem)
	router.POST("/order", middlewares.CustomerRequired, handlers.CreateOrder)
	router.GET("/orders", handlers.ListOrdersByCustomer)
	router.GET("/orders/:orderId", handlers.GetOrder)
}

func setUpAdminRoutes(router *gin.Engine) {
	adminGroup := router.Group("/admin")
	{
		// 不需 Admin 登入
		adminGroup.POST("/login", handlers.AdminLogin)

		// 需要 Admin 登入
		authorizedGroup := adminGroup.Group("")
		authorizedGroup.Use(middlewares.AdminRequired)
		{
			// Auth
			authorizedGroup.POST("/logout", handlers.Logout)
			authorizedGroup.POST("/signup", handlers.Signup)

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
		}
	}
}
