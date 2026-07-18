package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string `json:"code"`
	Price uint   `json:"price"`
}

var db *gorm.DB

// 初始化
func init() {
	//初始化数据库连接
	dsn := "root:qwq512369@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接错误：" + err.Error())
	}
	db.AutoMigrate(&Product{})

}
func main() {
	r := gin.Default()

	//获取数据
	r.GET("/product", func(c *gin.Context) {
		var products []Product
		db.Find(&products)
		c.JSON(http.StatusOK, gin.H{
			"message": "数据获取成功",
			"user":    products,
		})
	})
	//新增商品
	r.POST("/products", func(c *gin.Context) {
		var newProduct Product
		if err := c.ShouldBindJSON(&newProduct); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errors": err,
			})
			return
		}
		db.Create(&newProduct)
		c.JSON(http.StatusCreated, gin.H{
			"message": "添加商品成功",
			"product": newProduct,
		})

	})
	r.Run(":8080")
}
