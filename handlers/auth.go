package handlers

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"shop.go/config"
	"shop.go/models"
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

func Signup(ctx *gin.Context) {
	// 從 Request Body 拿資料
	req := SignupRequest{}
	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// 新增使用者
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	err = config.DB.Create(&user).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 返回成功 Response
	ctx.JSON(200, user)
}

func Login(ctx *gin.Context) {
	// 從 Request Body 拿資料
	req := LoginRequest{}
	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	// 查詢 user
	user := models.User{Email: req.Email, Password: req.Password}
	err = config.DB.Where("email = ?", user.Email).First(&user).Error
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	// 驗證密碼
	if !user.CheckPassword(req.Password) {
		ctx.JSON(http.StatusUnauthorized, "password not correct")
		return
	}

	// 設置 session
	session := sessions.Default(ctx)
	session.Set("user_id", user.ID)
	session.Set("user_name", user.Name)
	session.Set("user_role", user.Role)
	session.Save()

	// 返回成功 Response
	ctx.JSON(http.StatusOK, user)
}

func AdminLogin(ctx *gin.Context) {
	// 從 Request Body 拿資料
	req := LoginRequest{}
	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	// 查詢 user
	user := models.User{Email: req.Email}
	err = config.DB.Where("email = ?", user.Email).First(&user).Error
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	// 檢查身份
	if user.Role != "admin" && user.Role != "staff" {
		ctx.JSON(http.StatusUnauthorized, "此帳號無權限")
		return
	}

	// 驗證密碼
	if !user.CheckPassword(req.Password) {
		ctx.JSON(http.StatusUnauthorized, "password not correct")
		return
	}

	// 設置 session
	session := sessions.Default(ctx)
	session.Set("user_id", user.ID)
	session.Set("user_name", user.Name)
	session.Set("user_role", user.Role)
	err = session.Save()
	if err != nil {
		log.Println(err)
	}

	// 返回成功 Response
	ctx.JSON(http.StatusOK, user)
}

func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
	ctx.JSON(http.StatusOK, "登出")
}
