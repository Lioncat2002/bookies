package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils/password"
	"backend/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserSignUpData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserData struct {
	//ID   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Role string `json:"role" binding:"required"`
}

type AddCoin struct {
	//UserID uint    `json:"user_id" binding:"required"`
	Coin float32 `json:"coin" binding:"required"`
}

func AllUsers(c *gin.Context) {
	var users []models.User
	if err := services.DB.Preload("Author").Preload("Owns").Preload("Carts").Preload("Rating").Find(&users).Error; err != nil {
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

func LoginUser(c *gin.Context) {
	data := UserLoginData{}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user := models.User{}
	if err := services.DB.Where("email = ?", data.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err := password.VerifyPassword(user.Password, data.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	token, err := token.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   token,
	})
}

func GetOneUser(c *gin.Context) {
	id, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user := models.User{}
	if err := services.DB.Where("id = ?", id).Preload("Author").Preload("Owns").Preload("Cart").First(&user).Error; err != nil {
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

func DeleteUser(c *gin.Context) {
	id, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user := models.User{}
	if err := services.DB.Where("id = ?", id).First(&user).Unscoped().Delete(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   "deleted",
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
	pass, err := password.HashPassword(data.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user := models.User{}
	user.Email = data.Email
	user.Password = pass
	user.Carts = []models.Book{}
	user.Owns = []models.Book{}

	if err := services.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}
	token, err := token.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   token,
	})
}

func UpdateCoins(c *gin.Context) {
	var data AddCoin
	//getting data from the jwt token
	id, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user := models.User{}
	if err := services.DB.Where("id = ?", id).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	coins := user.Coins + data.Coin
	if err := services.DB.Where("id = ?", id).Find(&user).Update("coins", coins).Error; err != nil {
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

func UpdateUser(c *gin.Context) {
	var data UserData
	//getting data from the jwt token
	id, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user := models.User{}
	if err := services.DB.Where("id = ?", id).Find(&user).Update("name", data.Name).Update("role", data.Role).Error; err != nil {
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
