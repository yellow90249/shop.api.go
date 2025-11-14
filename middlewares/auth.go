package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func AuthRequire(userRole string) gin.HandlerFunc {
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

		if userRole != claims["user_role"].(string) {
			ctx.JSON(http.StatusUnauthorized, "身份錯誤")
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims["user_id"].(string))
		ctx.Set("user_role", claims["user_role"].(string))

		ctx.Next()
	}
}
