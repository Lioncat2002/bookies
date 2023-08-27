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
	UserRoute.PUT("/", controllers.AddUser)
	UserRoute.POST("/login", controllers.LoginUser)

	UserRoute.Use(middlewares.JwtAuth())
	UserRoute.GET("/one", controllers.GetOneUser)
	UserRoute.PATCH("/update", controllers.UpdateUser)
	UserRoute.DELETE("/delete", controllers.DeleteUser)
	UserRoute.PATCH("/buycoins", controllers.UpdateCoins)

	ItemRoute := router.Group("/api/book")
	ItemRoute.GET("/", controllers.AllBooks)
	ItemRoute.GET("/:id", controllers.GetOneBook)
	ItemRoute.GET("/search/:name", controllers.SearchBook)

	ItemRoute.Use(middlewares.JwtAuth())
	ItemRoute.PUT("/", controllers.CreateBook)
	ItemRoute.PUT("/upload/:id", controllers.AddBookUrl)
	ItemRoute.POST("/buy", controllers.BuyBook)
	ItemRoute.POST("/addcart", controllers.AddToCart)
	ItemRoute.POST("/rate", controllers.RateBook)
	router.Run(":8080")
}
func main() {
	services.ConnectDatabase()
	RunRouter()
}
