package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shop.go/config"
	"shop.go/models"
)

type AddCategoryRequest struct {
	Name        string `binding:"required"`
	Description string `binding:"required"`
}

func GetCategories(ctx *gin.Context) {
	categories := []models.Category{}
	err := config.DB.Find(&categories).Error
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func GetCategory(ctx *gin.Context) {

}

func AddCategory(ctx *gin.Context) {
	req := AddCategoryRequest{}
	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	err = config.DB.Create(&category).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func DeleteCategory(ctx *gin.Context) {

}
