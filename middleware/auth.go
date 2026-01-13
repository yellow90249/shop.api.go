package middleware

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"shop.go/enum"
	"shop.go/utils"
)

func getTokenStringFromAuthorizationHeader(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("你的 Bearer 前綴哩?")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}

func Auth(userRoleList ...enum.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := getTokenStringFromAuthorizationHeader(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, err.Error())
			ctx.Abort()
			return
		}

		token, err := utils.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, err.Error())
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, "Token is not valid!")
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, "Claim is not valid!")
			ctx.Abort()
			return
		}

		userRole := enum.UserRole(claims["user_role"].(string))
		if !slices.Contains(userRoleList, userRole) {
			ctx.JSON(http.StatusUnauthorized, "身份錯誤")
			ctx.Abort()
			return
		}

		if claims["user_id"] == nil {
			ctx.JSON(http.StatusUnauthorized, "user_id not exists")
			ctx.Abort()
			return
		}
		ctx.Set("user_id", claims["user_id"].(string))

		if claims["user_role"] == nil {
			ctx.JSON(http.StatusUnauthorized, "user_role not exists")
			ctx.Abort()
			return
		}
		ctx.Set("user_role", claims["user_role"].(string))

		ctx.Next()
	}
}
