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
	var query ListQuery

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

	// 只有當 CurrentPage 和 PerPage 都是 -1 時才返回全部，否則必須分頁
	if query.CurrentPage == -1 && query.PerPage == -1 {
		// 返回全部資料
		db.Find(&categories)
	} else {
		// 分頁查詢
		offset := (query.CurrentPage - 1) * query.PerPage
		db.Offset(offset).Limit(query.PerPage).Find(&categories)
	}

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
