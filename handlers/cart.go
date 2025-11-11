package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"shop.go/config"
	"shop.go/models"
)

type AddCartItemRequest struct {
	ProductID uint
	Quantity  uint
	UnitPrice float64
}

func AddCartItem(ctx *gin.Context) {
	// 找 user
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	user := models.User{}
	err := config.DB.First(&user, userID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// 從 body 拿資料
	req := AddCartItemRequest{}
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// 存記錄到 CartItem table
	cartItem := models.CartItem{
		UserID:    user.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
	}
	err = config.DB.Create(&cartItem).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, userID)
}
