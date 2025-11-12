package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"shop.go/config"
	"shop.go/models"
)

type CreateOrderRequest struct {
	RecipientName    string
	RecipientPhone   string
	RecipientEmail   string
	RecipientAddress string
	TotalAmount      float64
	PaymentMethod    string
}

func CreateOrder(ctx *gin.Context) {
	tx := config.DB.Begin()

	// 建立訂單
	session := sessions.Default(ctx)
	userID := session.Get("user_id").(uint)

	req := CreateOrderRequest{}
	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	order := models.Order{
		UserID:           userID,
		RecipientName:    req.RecipientName,
		RecipientPhone:   req.RecipientPhone,
		RecipientEmail:   req.RecipientEmail,
		RecipientAddress: req.RecipientAddress,
		TotalAmount:      req.TotalAmount,
		PaymentMethod:    req.PaymentMethod,
		Status:           "pending",
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
