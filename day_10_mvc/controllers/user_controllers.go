package controllers

import (
	"go_learning/two_week/day_10_mvc/middleware"
	"go_learning/two_week/day_10_mvc/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// Register 处理注册请求
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		if err := ctx.ShouldBindJSON(user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		bytPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "密码加密失败",
			})
			return
		}
		user.Password = string(bytPassword)
		if err := models.CreateUser(db, &user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "用户已存在或创建失败"})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "注册成功",
		})
	}
}

// 只负责处理HTTP请求和响应
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		user, err := models.GetUserByUserName(db, req.UserName)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "用户不存在",
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "密码不正确"})
			return
		}
		tokenStr, err := middleware.GenToken(req.UserName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "生成Token失败"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "登陆成功",
			"token":   tokenStr,
		})

	}
}

// GetProfile 获取当前用户信息（受保护接口）
func GetProfile(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{
		"message":      "访问受保护接口成功",
		"current_user": username,
	})
}
