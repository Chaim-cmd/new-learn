package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 定义key
var jwtSecret = []byte("my_super_key")
var db *gorm.DB

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserName string `json:"string"`
	Password string `json:"-"`
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// 自定义
type MyClaims struct {
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

// 生成token
func GenToken(username string) (string, error) {
	claims := MyClaims{
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// JWT 鉴权中间件
func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "未提供token",
			})
			ctx.Abort() //阻止后续
			return

		}
		//提取token字符串
		parts := strings.SplitN(authHeader, "", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token格式错误",
			})
			ctx.Abort()
			return

		}
		//解析验证token
		token, err := jwt.ParseWithClaims(parts[1], &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil && !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "token无效或已过期",
			})
			ctx.Abort()
			return

		}
		claims, ok := token.Claims.(*MyClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "token解析失败",
			})
			ctx.Abort()
			return
		}
		ctx.Set("Username", claims.UserName)
		ctx.Next()

	}
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败:" + err.Error())
	}
	db.AutoMigrate(&User{})

	r := gin.Default()
	r.POST("/register", func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		hsPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		db.Create(&User{UserName: user.UserName, Password: string(hsPassword)})
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "注册成功",
		})
	})
	r.POST("login", func(ctx *gin.Context) {
		var req LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		var user User
		if db.Where("user_name = ?", req.UserName).First(&user).Error != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "用户不存在",
			})
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "密码不正确",
			})
		}
		tokenStr, _ := GenToken(req.UserName)
		ctx.JSON(http.StatusOK, gin.H{"token": tokenStr})
	})

	//受保护的组
	authGroup := r.Group("/api")
	authGroup.Use(JWTAuth())
	{
		authGroup.GET("/profile", func(c *gin.Context) {
			// 从上下文中获取中间件存入的 username
			username, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{
				"message":      "这是受保护的接口",
				"current_user": username,
			})
		})
	}
	r.Run(":8080")
}
