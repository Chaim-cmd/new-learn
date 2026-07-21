package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 1. 定义 JWT 密钥
var jwtSecret = []byte("my_super_key")

// 2. 全局数据库变量
var db *gorm.DB

// 3. 用户模型
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserName string `json:"username"`
	Password string `json:"-"` // 返回前端时自动隐藏密码
}

// 4. 登录请求体
type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// 5. 生成 Token 的辅助函数
func generateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24小时后过期
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func main() {
	// 6. 初始化数据库连接
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	db.AutoMigrate(&User{})

	r := gin.Default()

	// 7. 注册接口 (POST /register)
	r.POST("/register", func(ctx *gin.Context) {
		var user User
		// 注意：必须传指针 &user
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 密码哈希加密 (使用默认成本 10)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "加密失败"})
			return
		}

		// 存入数据库
		newUser := User{
			UserName: user.UserName,
			Password: string(hashedPassword),
		}
		if result := db.Create(&newUser); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "用户已存在或创建失败"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "注册成功"})
	})

	// 8. 登录接口 (POST /login)
	r.POST("/login", func(ctx *gin.Context) {
		var req LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 根据用户名查找用户
		var user User
		result := db.Where("user_name = ?", req.UserName).First(&user)
		if result.Error != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "用户不存在"})
			return
		}

		// 比对密码：数据库里的哈希密码 vs 用户输入的明文密码
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "密码不正确"})
			return
		}

		// 生成 Token
		tokenStr, err := generateToken(req.UserName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "生成Token失败"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "登录成功",
			"token":   tokenStr,
		})
	})

	r.Run(":8080")
}
