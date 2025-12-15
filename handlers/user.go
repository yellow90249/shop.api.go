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

type ListUsersQuery struct {
	CurrentPage int    `binding:"required"`
	PerPage     int    `binding:"required"`
	Role        string `binding:"required"`
	Name        string
}

type ListUsersResponse struct {
	List  []models.User
	Total int64
}

func ListUsers(ctx *gin.Context) {
	var users []models.User
	var total int64
	var query ListUsersQuery

	// 自動綁定和驗證
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 建立查詢
	db := config.DB.Model(&models.User{})

	// 角色篩選
	db = db.Where("role = ?", query.Role)

	// 如果有搜尋名稱，加入模糊搜尋
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}

	// 計算總數
	db.Count(&total)

	// 只有當 CurrentPage 和 PerPage 都是 -1 時才返回全部，否則必須分頁
	if query.CurrentPage == -1 && query.PerPage == -1 {
		// 返回全部資料
		db.Find(&users)
	} else {
		// 分頁查詢
		offset := (query.CurrentPage - 1) * query.PerPage
		db.Offset(offset).Limit(query.PerPage).Find(&users)
	}

	ctx.JSON(http.StatusOK, ListUsersResponse{
		List:  users,
		Total: total,
	})
}

func UpdateUserImage(ctx *gin.Context) {
	// 找商品
	userId := ctx.Param("userId")
	user := models.User{}
	err := config.DB.First(&user, userId).Error
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

	err = os.Remove(user.Avatar)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "刪除失敗")
		return
	}

	// 存 DB
	user.Avatar = dst
	config.DB.Save(&user)

	ctx.JSON(http.StatusOK, "使用者圖片更新成功")
}
