package main

import (
	"go_learning/two_week/day_10_mvc/config"
	"go_learning/two_week/day_10_mvc/routes"
)

func main() {
	db := config.InitDB()

	db.AutoMigrate()

	r := routes.SetupRouter(db)
	r.Run(":8080")
}
