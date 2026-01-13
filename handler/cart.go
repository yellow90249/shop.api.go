package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shop.go/boot"
	"shop.go/model"
)

type AddCartItemRequest struct {
	ProductID uint
	Quantity  uint
	UnitPrice float64
}

type UpdateCartItemRequest struct {
	Quantity uint
}

func AddCartItem(ctx *gin.Context) {
	// 找 user
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "userID not exist")
		return
	}
	user := model.User{}
	err := boot.DB.First(&user, userID).Error
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
	cartItem := model.CartItem{
		UserID:    user.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
	}
	err = boot.DB.Create(&cartItem).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, "success")
}

func UpdateCartItemQuantity(ctx *gin.Context) {
	// Get Data
	cartId := ctx.Param("cartItemId")
	req := UpdateCartItemRequest{}
	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// 更新
	cart := model.CartItem{}
	err = boot.DB.First(&cart, cartId).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	cart.Quantity = req.Quantity
	boot.DB.Save(&cart)

	ctx.JSON(http.StatusOK, "success")
}

func DeleteCartItem(ctx *gin.Context) {
	cartItemId := ctx.Param("cartItemId")

	err := boot.DB.Unscoped().Delete(&model.CartItem{}, cartItemId).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "已刪除")
}

func DeleteAllCartItem(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "userID not exist")
		return
	}
	err := boot.DB.Where("user_id = ?", userID).Delete(&model.CartItem{}).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "已刪除")
}
