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

type ListCategoryQuery struct {
	CurrentPage int    `form:"currentPage" binding:"required,min=1"`
	PerPage     int    `form:"perPage" binding:"required,min=1"`
	Name        string `form:"name"`
}

type ListCategoryResponse struct {
	List  []models.Category
	Total int64
}

type DeleteCategoryRequest struct {
	CategoryID uint `binding:"required"`
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

func ListCategories(ctx *gin.Context) {
	var categories []models.Category
	var total int64
	var query ListCategoryQuery

	// 自動綁定和驗證
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 建立查詢
	db := config.DB.Model(&models.Category{})

	// 如果有搜尋名稱，加入模糊搜尋
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}

	// 計算總數
	db.Count(&total)

	// 分頁查詢
	offset := (query.CurrentPage - 1) * query.PerPage
	db.Offset(offset).Limit(query.PerPage).Find(&categories)

	ctx.JSON(http.StatusOK, ListCategoryResponse{
		List:  categories,
		Total: total,
	})
}

func DeleteCategory(ctx *gin.Context) {
	categoryId := ctx.Param("categoryId")

	err := config.DB.Unscoped().Delete(&models.Category{}, categoryId).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "已刪除")
}
