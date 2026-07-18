package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var users = []User{
	{ID: 1, Name: "张三"},
	{ID: 2, Name: "李四"},
}

func main() {

	r := gin.Default()

	r.GET("/users", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    users,
		})
	})
	r.GET("/users/:id", func(ctx *gin.Context) {

		for _, user := range users {
			if ctx.Param("id") == string(rune(user.ID)) {
				ctx.JSON(http.StatusOK, user)
			}
		}
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
	})
	r.POST("/users", func(ctx *gin.Context) {
		var newUser User
		if err := ctx.ShouldBindJSON(&newUser); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		Users := append(users, newUser)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "添加成功",
			"data":    Users,
		})
	})

	r.Run(":8080")
}
