package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"shop.go/config"
	"shop.go/models"
)

type CreateOrderRequest struct {
	RecipientName    string  `binding:"required"`
	RecipientPhone   string  `binding:"required"`
	RecipientEmail   string  `binding:"required"`
	RecipientAddress string  `binding:"required"`
	TotalAmount      float64 `binding:"required"`
	PaymentMethod    string  `binding:"required"`
}

type ListOrdersQuery struct {
	CurrentPage int `form:"currentPage" binding:"required"`
	PerPage     int `form:"perPage" binding:"required"`
}

type ListOrdersResponse struct {
	Message string
	List    []models.Order
	Total   int64
}

type UpdateOrderRequest struct {
	Status string `binding:"required"`
}

func CreateOrder(ctx *gin.Context) {
	tx := config.DB.Begin()

	// 建立訂單
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "userID not exist")
		return
	}

	req := CreateOrderRequest{}
	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newVal, err := strconv.ParseUint(userID.(string), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	order := models.Order{
		UserID:           uint(newVal),
		RecipientName:    req.RecipientName,
		RecipientPhone:   req.RecipientPhone,
		RecipientEmail:   req.RecipientEmail,
		RecipientAddress: req.RecipientAddress,
		TotalAmount:      req.TotalAmount,
		PaymentMethod:    req.PaymentMethod,
		Status:           models.OrderStatusPending,
	}
	err = tx.Create(&order).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// 拿購物車
	var cartItems []models.CartItem
	err = tx.Where("user_id = ?", userID).Find(&cartItems).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// 用購物車建立訂單細項
	var orderItems []models.OrderItem
	for _, cartItem := range cartItems {
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			UnitPrice: cartItem.UnitPrice,
		}
		orderItems = append(orderItems, orderItem)
	}
	err = tx.Create(&orderItems).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// 清空購物車
	err = tx.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, "建立訂單成功")
}

func GetOrder(ctx *gin.Context) {
	orderId := ctx.Param("orderId")
	order := models.Order{}
	err := config.DB.Preload("OrderItems.Product").First(&order, orderId).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func ListOrdersByCustomer(ctx *gin.Context) {
	var orders []models.Order
	var total int64
	var query ListOrdersQuery
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "userID not exist")
		return
	}

	// 自動綁定和驗證
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 建立查詢
	db := config.DB.Model(&models.Order{}).Preload("OrderItems")

	// 如果有分類，加入分類篩選
	if userID != 0 {
		db = db.Where("user_id = ?", userID)
	}

	// 計算總數
	db.Count(&total)

	// 加入排序
	db = db.Order("created_at DESC")

	// 只有當 CurrentPage 和 PerPage 都是 -1 時才返回全部，否則必須分頁
	if query.CurrentPage == -1 && query.PerPage == -1 {
		// 返回全部資料
		db.Find(&orders)
	} else {
		// 分頁查詢
		offset := (query.CurrentPage - 1) * query.PerPage
		db.Offset(offset).Limit(query.PerPage).Find(&orders)
	}

	ctx.JSON(http.StatusOK, ListOrdersResponse{
		Message: "success",
		List:    orders,
		Total:   total,
	})
}

func ListOrdersByAdmin(ctx *gin.Context) {
	var orders []models.Order
	var total int64
	var query ListOrdersQuery

	// 自動綁定和驗證
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 建立查詢
	db := config.DB.Model(&models.Order{}).Preload("OrderItems")

	// 計算總數
	db.Count(&total)

	// 加入排序
	db = db.Order("created_at DESC")

	// 只有當 CurrentPage 和 PerPage 都是 -1 時才返回全部，否則必須分頁
	if query.CurrentPage == -1 && query.PerPage == -1 {
		// 返回全部資料
		db.Find(&orders)
	} else {
		// 分頁查詢
		offset := (query.CurrentPage - 1) * query.PerPage
		db.Offset(offset).Limit(query.PerPage).Find(&orders)
	}

	ctx.JSON(http.StatusOK, ListOrdersResponse{
		List:  orders,
		Total: total,
	})
}

func UpdateOrder(ctx *gin.Context) {
	// 找商品
	orderId := ctx.Param("orderId")
	order := models.Order{}
	err := config.DB.First(&order, orderId).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	req := UpdateOrderRequest{}
	err = ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	order.Status = models.OrderStatus(req.Status)

	config.DB.Save(&order)

	ctx.JSON(http.StatusOK, "更新成功")
}
