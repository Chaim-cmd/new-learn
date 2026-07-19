package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

var db *gorm.DB

func init() {
	dsn := "root:qwq512369@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接错误：" + err.Error())
	}
	db.AutoMigrate()

}

func main() {
	r := gin.Default()

	//自定义一个中间件
	//再每一个请求开始前，打印请求方法和路径
	r.Use(func(ctx *gin.Context) {
		start := time.Now()

		println("请求开始：", ctx.Request.Method, ctx.Request.URL.Path)

		ctx.Next()

		print("请求结束，耗时：", time.Since(start))
	})

	//获取所有商品
	r.GET("/products", func(ctx *gin.Context) {
		var products []Product
		result := db.Find(&products)
		if result.Error != nil {
			println("数据库查询错误：", result.Error.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "服务器内部错误，请联系管理员",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "获取成功",
			"data":    products,
		})

	})

	//新增商品
	r.POST("/product", func(c *gin.Context) {
		var NewProduct Product
		if err := c.ShouldBindJSON(&NewProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "请求参数错误" + err.Error(),
			})
			return
		}
		if result := db.Create(&NewProduct); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "新增数据失败",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "写入成功",
			"data":    NewProduct,
		})
	})
	r.Run(":8080")
}
