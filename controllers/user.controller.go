package controllers

import (
	"backend/models"
	"backend/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserSignUpData struct {
	Email string `json:"email" binding:"required"`
}

type UserData struct {
	ID   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type AddCoin struct {
	UserID uint    `json:"user_id" binding:"required"`
	Coin   float32 `json:"coin" binding:"required"`
}

func AllUsers(c *gin.Context) {
	var users []models.User
	if err := services.DB.Preload("Author").Preload("Owns").Find(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   users,
	})
}

func GetOneUser(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("id", id)
	//id, _ := strconv.ParseInt(query, 10, 32)
	user := models.User{}
	if err := services.DB.Where("id = ?", id).Preload("Author").Preload("Owns").First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   user,
	})
}

func AddUser(c *gin.Context) {
	var data UserSignUpData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user := models.User{}
	user.Email = data.Email
	//user.Owned = []models.Item{}
	if err := services.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("id", id)
	//id, _ := strconv.ParseInt(query, 10, 32)
	user := models.User{}
	if err := services.DB.Where("id = ?", id).Preload("Author").Preload("Owns").First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	services.DB.Delete(&user)
	c.JSON(http.StatusAccepted, gin.H{
		"status": "deleted",
		"data":   user,
	})
}

func UpdateCoins(c *gin.Context) {
	var data AddCoin
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user := models.User{}
	if err := services.DB.Where("id = ?", data.UserID).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	coins := user.Coins + data.Coin
	if err := services.DB.Where("id = ?", data.UserID).Find(&user).Update("coins", coins).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status": "updated",
		"data":   user,
	})

}

func UpdateUser(c *gin.Context) {
	var data UserData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user := models.User{}
	if err := services.DB.Where("id = ?", data.ID).Find(&user).Update("name", data.Name).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "updated",
		"data":   user,
	})
}
