package routes

import (
	"github.com/gin-gonic/gin"
	"shop.go/enum"
	"shop.go/handler"
	"shop.go/middleware"
)

func Setup(router *gin.Engine) {
	api := router.Group("/api")

	Auth := middleware.Auth
	RoleAdmin := enum.RoleAdmin
	RoleUser := enum.RoleUser

	// 權限
	api.POST("/user/signup", handler.Signup("user"))
	api.POST("/admin/signup", handler.Signup("admin"))
	api.POST("/user/login", handler.Login([]string{"user"}))
	api.POST("/admin/login", handler.Login([]string{"admin", "guest"}))

	// 用戶
	api.GET("/me", Auth(RoleAdmin, RoleUser), handler.GetUser)
	api.GET("/users", Auth(RoleAdmin), handler.ListUsers)
	api.PUT("/user/avatar", Auth(RoleAdmin, RoleUser), handler.UpdateUserImage)
	api.PUT("/user/:userId/password", Auth(RoleAdmin), handler.ResetUserPassword)

	// 種類
	api.GET("/categories", handler.ListCategories)
	api.POST("/category", Auth(RoleAdmin), handler.AddCategory)
	api.PUT("/category/:categoryId", Auth(RoleAdmin), handler.UpdateCategory)
	api.DELETE("/category/:categoryId", Auth(RoleAdmin), handler.DeleteCategory)

	// 商品
	api.GET("/products", handler.ListProducts)
	api.GET("/product/:productId", Auth(RoleAdmin, RoleUser), handler.GetProduct)
	api.POST("/product", Auth(RoleAdmin), handler.AddProduct)
	api.PUT("/product/:productId", Auth(RoleAdmin), handler.UpdateProduct)
	api.PUT("/product/:productId/image", Auth(RoleAdmin), handler.UpdateProductImage)
	api.DELETE("/product/:productId", Auth(RoleAdmin), handler.DeleteProduct)

	// 訂單
	api.GET("/order/:orderId", handler.GetOrder)
	api.GET("/user/me/orders", Auth(RoleUser), handler.ListOrdersByCustomer)
	api.GET("/orders", Auth(RoleAdmin), handler.ListOrdersByAdmin)
	api.POST("/order", Auth(RoleUser), handler.CreateOrder)
	api.PUT("/order/:orderId", Auth(RoleAdmin), handler.UpdateOrder)

	// 購物車
	api.POST("/cart/item", Auth(RoleUser), handler.AddCartItem)
	api.PUT("/cart/item/:cartItemId", Auth(RoleUser), handler.UpdateCartItemQuantity)
	api.DELETE("/cart/item/:cartItemId", Auth(RoleUser), handler.DeleteCartItem)
	api.DELETE("/cart/item/all", Auth(RoleUser), handler.DeleteAllCartItem)

	// 測試
	api.GET("/hello", func(ctx *gin.Context) { ctx.JSON(200, "cool") })
}
