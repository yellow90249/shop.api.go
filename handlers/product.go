package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shop.go/config"
	"shop.go/models"
)

type AddProductRequest struct {
	Name          string  `form:"Name" binding:"required"`
	CategoryID    uint    `form:"CategoryID" binding:"required"`
	Price         float64 `form:"Price" binding:"required"`
	StockQuantity uint    `form:"StockQuantity" binding:"required"`
	Description   string  `form:"Description" binding:"required"`
}

type ListProductsResponse struct {
	List  []models.Product
	Total int64
}

func AddProduct(ctx *gin.Context) {
	req := AddProductRequest{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	file, err := ctx.FormFile("UploadedFile")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 儲存檔案
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	dst := filepath.Join("uploads", filename)

	log.Println(dst)

	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, "儲存失敗")
		return
	}

	// DB 存紀錄
	product := models.Product{
		CategoryID:    req.CategoryID,
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		ImageURL:      dst,
	}

	err = config.DB.Create(&product).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "cool")
}

func UpdateProduct(ctx *gin.Context) {
	// 找商品
	productId := ctx.Param("productId")
	product := models.Product{}
	err := config.DB.First(&product, productId).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	req := AddProductRequest{}
	err = ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	product.Name = req.Name
	product.CategoryID = req.CategoryID
	product.Price = req.Price
	product.StockQuantity = req.StockQuantity
	product.Description = req.Description

	config.DB.Save(&product)

	ctx.JSON(http.StatusOK, "更新成功")
}

func UpdateProductImage(ctx *gin.Context) {
	// 找商品
	productId := ctx.Param("productId")
	product := models.Product{}
	err := config.DB.First(&product, productId).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// 存 FS
	file, err := ctx.FormFile("UploadedFile")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	dst := filepath.Join("uploads", filename)
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, "儲存失敗")
		return
	}

	err = os.Remove(product.ImageURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "刪除失敗")
		return
	}

	// 存 DB
	product.ImageURL = dst
	config.DB.Save(&product)

	ctx.JSON(http.StatusOK, "商品圖片更新成功")
}

func ListProducts(ctx *gin.Context) {
	var products []models.Product
	var total int64
	var query ListQuery

	// 自動綁定和驗證
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 建立查詢
	db := config.DB.Model(&models.Product{}).Preload("Category")

	// 如果有搜尋名稱，加入模糊搜尋
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}

	// 計算總數
	db.Count(&total)

	// 只有當 CurrentPage 和 PerPage 都是 -1 時才返回全部，否則必須分頁
	if query.CurrentPage == -1 && query.PerPage == -1 {
		// 返回全部資料
		db.Find(&products)
	} else {
		// 分頁查詢
		offset := (query.CurrentPage - 1) * query.PerPage
		db.Offset(offset).Limit(query.PerPage).Find(&products)
	}

	ctx.JSON(http.StatusOK, ListProductsResponse{
		List:  products,
		Total: total,
	})
}

func DeleteProduct(ctx *gin.Context) {
	productId := ctx.Param("productId")

	err := config.DB.Unscoped().Delete(&models.Product{}, productId).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "已刪除")
}
