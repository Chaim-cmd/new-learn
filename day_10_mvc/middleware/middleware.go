package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("JWT-KEY")

// 自定义结构体
type MyClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 生成Token
func GenToken(username string) (string, error) {
	claims := MyClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)

}

// 鉴权中间件
func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "为提供Token",
			})
			ctx.Abort()
			return
		}
		parts := strings.SplitN(authHeader, "", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token 格式错误",
			})
			ctx.Abort()
			return
		}

		//解析验证
		claims := &MyClaims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token无效或者已经过期",
			})
			ctx.Abort()
			return
		}
		ctx.Set("username", claims.Username)
		ctx.Next()

	}
}
