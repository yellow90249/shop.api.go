package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
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

	// 大家都進來吧（開發用）
	if os.Getenv("GO_ENV") == "development" {
		router.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		}))
	}

	// 路由設定
	setUpPublicRoutes(router)
	setUpWebRoutes(router)
	setUpAdminRoutes(router)

	// 啟動服務
	router.Run(":" + os.Getenv("APP_PORT"))
}

func setUpPublicRoutes(router *gin.Engine) {
	router.Static("/api/uploads", "./uploads")
}

func setUpWebRoutes(router *gin.Engine) {
	customerGroup := router.Group("/api")
	{
		// ======================== 不需 Customer 登入 ========================
		customerGroup.POST("/login", handlers.Login([]string{"customer"}))
		customerGroup.GET("/categories", handlers.ListCategories)
		customerGroup.GET("/products", handlers.ListProducts)
		customerGroup.GET("/products/:productId", handlers.GetProduct)

		// ======================== 需要 Customer 登入 ========================
		// Auth
		customerGroup.GET("/me", middlewares.AuthRequire("customer"), handlers.GetUser)

		// 購物車
		customerGroup.POST("/cart/items", middlewares.AuthRequire("customer"), handlers.AddCartItem)
		customerGroup.PUT("/cart/items/:cartItemId", middlewares.AuthRequire("customer"), handlers.UpdateCartItemQuantity)
		customerGroup.DELETE("/cart/items/:cartItemId", middlewares.AuthRequire("customer"), handlers.DeleteCartItem)
		customerGroup.DELETE("/cart/items/all", middlewares.AuthRequire("customer"), handlers.DeleteAllCartItem)

		// 訂單
		customerGroup.POST("/order", middlewares.AuthRequire("customer"), handlers.CreateOrder)
		customerGroup.GET("/orders", middlewares.AuthRequire("customer"), handlers.ListOrdersByCustomer)
		customerGroup.GET("/orders/:orderId", middlewares.AuthRequire("customer"), handlers.GetOrder)
	}

}

func setUpAdminRoutes(router *gin.Engine) {
	adminGroup := router.Group("/api/admin")
	{
		// ======================== 不需 Admin 登入 ========================
		adminGroup.POST("/login", handlers.Login([]string{"admin", "staff"}))

		// ======================== 需要 Admin 登入 ========================
		// Auth
		adminGroup.POST("/signup", middlewares.AuthRequire("admin"), handlers.Signup)
		adminGroup.GET("/me", middlewares.AuthRequire("admin"), handlers.GetUser)

		// 種類
		adminGroup.POST("/categories", middlewares.AuthRequire("admin"), handlers.AddCategory)
		adminGroup.GET("/categories", middlewares.AuthRequire("admin"), handlers.ListCategories)
		adminGroup.PUT("/categories/:categoryId", middlewares.AuthRequire("admin"), handlers.UpdateCategory)
		adminGroup.DELETE("/categories/:categoryId", handlers.DeleteCategory)

		// 商品
		adminGroup.POST("/products", middlewares.AuthRequire("admin"), handlers.AddProduct)
		adminGroup.GET("/products", middlewares.AuthRequire("admin"), handlers.ListProducts)
		adminGroup.PUT("/products/:productId", middlewares.AuthRequire("admin"), handlers.UpdateProduct)
		adminGroup.PUT("/products/:productId/image", middlewares.AuthRequire("admin"), handlers.UpdateProductImage)
		adminGroup.DELETE("/products/:productId", middlewares.AuthRequire("admin"), handlers.DeleteProduct)

		// 用戶
			adminGroup.GET("/users", middlewares.AuthRequire("admin"), handlers.ListUsers)
			adminGroup.PUT("/users/:userId/image", middlewares.AuthRequire("admin"), handlers.UpdateUserImage)
	}
}
