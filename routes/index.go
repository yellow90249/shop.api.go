package routes

import (
	"github.com/gin-gonic/gin"
	"shop.go/handlers"
	"shop.go/middlewares"
)

func Setup(router *gin.Engine) {
	api := router.Group("/api")

	// Auth
	api.POST("/user/signup", handlers.Signup("user"))
	api.POST("/admin/signup", handlers.Signup("admin"))
	api.POST("/user/login", handlers.Login([]string{"user"}))
	api.POST("/admin/login", handlers.Login([]string{"admin", "guest"}))
	api.GET("/me", handlers.GetUser)

	// 用戶
	api.GET("/users", middlewares.AuthRequire("admin"), handlers.ListUsers)
	api.PUT("/users/:userId/image", middlewares.AuthRequire("admin"), handlers.UpdateUserImage)

	// 種類
	api.GET("/categories", handlers.ListCategories)
	api.POST("/categories", middlewares.AuthRequire("admin"), handlers.AddCategory)
	api.PUT("/categories/:categoryId", middlewares.AuthRequire("admin"), handlers.UpdateCategory)
	api.DELETE("/categories/:categoryId", middlewares.AuthRequire("admin"), handlers.DeleteCategory)

	// 商品
	api.GET("/products", handlers.ListProducts)
	api.GET("/products/:productId", handlers.GetProduct)
	api.POST("/products", middlewares.AuthRequire("admin"), handlers.AddProduct)
	api.PUT("/products/:productId", middlewares.AuthRequire("admin"), handlers.UpdateProduct)
	api.PUT("/products/:productId/image", middlewares.AuthRequire("admin"), handlers.UpdateProductImage)
	api.DELETE("/products/:productId", middlewares.AuthRequire("admin"), handlers.DeleteProduct)

	// 訂單
	api.POST("/order", handlers.CreateOrder)
	api.GET("/user/orders", handlers.ListOrdersByCustomer)
	api.GET("/orders/:orderId", handlers.GetOrder)
	api.GET("/admin/orders", middlewares.AuthRequire("admin"), handlers.ListOrdersByAdmin)
	api.PUT("/orders/:orderId", middlewares.AuthRequire("admin"), handlers.UpdateOrder)

	// 購物車
	api.POST("/cart/items", middlewares.AuthRequire("customer"), handlers.AddCartItem)
	api.PUT("/cart/items/:cartItemId", middlewares.AuthRequire("customer"), handlers.UpdateCartItemQuantity)
	api.DELETE("/cart/items/:cartItemId", middlewares.AuthRequire("customer"), handlers.DeleteCartItem)
	api.DELETE("/cart/items/all", middlewares.AuthRequire("customer"), handlers.DeleteAllCartItem)
}
