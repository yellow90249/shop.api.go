package handler

import (
	"log"
	"net/http"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"shop.go/boot"
	"shop.go/model"
	"shop.go/utils"
)

type SignupRequest struct {
	Name     string `binding:"required"`
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

type LoginRequest struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func Signup(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := SignupRequest{}

		validRoles := []string{"admin", "guest", "user"}
		if !slices.Contains(validRoles, role) {
			ctx.JSON(http.StatusBadRequest, "Role is not valid")
			return
		}

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
		file.Filename = uuid.New().String() + ext
		log.Println(file.Filename)

		err = boot.UploadFile(ctx, file)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		// DB 存紀錄
		user := model.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
			Role:     role,
			Avatar:   file.Filename,
		}

		err = user.HashPassword()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		err = boot.DB.Create(&user).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}

func Login(userRoleList []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 從 Request Body 拿資料
		req := LoginRequest{}
		err := ctx.ShouldBindBodyWithJSON(&req)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		// 查詢 user
		user := model.User{Email: req.Email}
		err = boot.DB.Where("email = ?", user.Email).First(&user).Error
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, err.Error())
			return
		}

		// 檢查身份（得是 userRoleList 裡的身份才能登入）
		if !slices.Contains(userRoleList, user.Role) {
			ctx.JSON(http.StatusUnauthorized, "role invalid")
			return
		}

		// 驗證密碼
		if !user.CheckPassword(req.Password) {
			ctx.JSON(http.StatusUnauthorized, "password not correct")
			return
		}

		// 產生 token
		userId := strconv.FormatUint(uint64(user.ID), 10)
		token, err := utils.GenerateToken(userId, user.Role, user.Name)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		// 返回成功 Response
		ctx.JSON(http.StatusOK, token)
	}
}

func GetUser(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "userID not exist")
		return
	}

	user := model.User{}
	err := boot.DB.
		Preload("CartItems", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("CartItems.Product").First(&user, userID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}
