package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Bookdata struct {
	AuthorID    uint    `json:"author_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"desc"`
	Price       float32 `json:"price" binding:"required"`
	Tag         string  `json:"tag" binding:"required"`
}

type BuyBookData struct {
	UserID uint `json:"user_id" binding:"required"`
	BookID uint `json:"item_id" binding:"required"`
}

func BuyBook(c *gin.Context) {
	var buyBookData BuyBookData
	if err := c.ShouldBindJSON(&buyBookData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := models.User{}
	if err := services.DB.Where("id = ?", buyBookData.UserID).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	item := models.Book{}
	if err := services.DB.Where("id = ?", buyBookData.BookID).Find(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	coins := user.Coins - item.Price
	if coins < 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "insufficient balance",
		})
		return
	}
	if err := services.DB.Where("id = ?", buyBookData.UserID).Find(&user).Update("coins", coins).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := services.DB.Where("id = ?", buyBookData.BookID).Find(&item).Update("current_owner_id", buyBookData.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   item,
	})
}

func CreateBook(c *gin.Context) {
	var itemData Bookdata
	if err := c.ShouldBindJSON(&itemData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	item := models.Book{}
	item.UserID = itemData.AuthorID
	item.Name = itemData.Name
	item.Description = itemData.Description
	item.Tag = itemData.Tag
	//item.CurrentOwnerID = itemData.AuthorID
	item.Price = itemData.Price
	if err := services.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}
	//services.DB.Debug().Model(&models.User{}).Related()
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   item,
	})
}

func GetOneBook(c *gin.Context) {
	id := c.Param("id")
	item := models.Book{}
	if err := services.DB.Where("id = ?", id).First(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   item,
	})
}

func AllBooks(c *gin.Context) {
	var items []models.Book
	if err := services.DB.Find(&items).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   items,
	})
}
