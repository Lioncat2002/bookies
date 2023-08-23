package main

import (
	"backend/controllers"
	"backend/middlewares"
	"backend/services"

	"github.com/gin-gonic/gin"
)

var PostRoute *gin.RouterGroup

func RunRouter() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	UserRoute := router.Group("/api/user")
	UserRoute.GET("/", controllers.AllUsers)
	UserRoute.POST("/", controllers.AddUser)
	UserRoute.POST("/login", controllers.LoginUser)
	UserRoute.GET("/one", controllers.GetOneUser)
	UserRoute.Use(middlewares.JwtAuth())
	UserRoute.POST("/update", controllers.UpdateUser)
	UserRoute.POST("/buycoins", controllers.UpdateCoins)

	ItemRoute := router.Group("/api/book")

	ItemRoute.GET("/", controllers.AllBooks)
	//ItemRoute.Use(middlewares.JwtAuth())
	ItemRoute.GET("/:id", controllers.GetOneBook)
	ItemRoute.POST("/", controllers.CreateBook)
	ItemRoute.POST("/buy", controllers.BuyBook)

	router.Run(":8080")
}
func main() {
	services.ConnectDatabase()
	RunRouter()
}
