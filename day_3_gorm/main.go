package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 定义一个Product 结构体
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	dsn := "root:qwq512369@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"

	//使用gorm.Open 连接mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败")
	}

	//自动迁移表
	db.AutoMigrate(&Product{})

	//插入数据
	db.Create(&Product{Code: "A001", Price: 100})
	db.Create(&Product{Code: "A002", Price: 40})
	db.Create(&Product{Code: "A003", Price: 120})

	//查询
	var Products []Product
	db.Where("price > ?", 50).Find(&Products)

	fmt.Println(Products)
}
