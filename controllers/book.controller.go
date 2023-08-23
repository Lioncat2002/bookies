package controllers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"backend/utils/token"

	"github.com/gin-gonic/gin"
)

type Bookdata struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"desc"`
	Price       float32 `json:"price" binding:"required"`
	Tag         string  `json:"tag" binding:"required"`
}

type BuyBookData struct {
	/* UserToken uint `json:"user_token" binding:"required"` */
	BookID uint `json:"item_id" binding:"required"`
}

func BuyBook(c *gin.Context) {
	buyBookData := BuyBookData{}
	if err := c.ShouldBindJSON(&buyBookData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	user := models.User{}
	if err := services.DB.Where("id = ?", id).Find(&user).Error; err != nil {
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
	if err := services.DB.Where("id = ?", id).Find(&user).Update("coins", coins).Error; err != nil {
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
	var bookData Bookdata
	if err := c.ShouldBindJSON(&bookData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	id, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	book := models.Book{}
	book.UserID = id
	book.Name = bookData.Name
	book.Description = bookData.Description
	book.Tag = bookData.Tag
	//item.CurrentOwnerID = itemData.AuthorID
	book.Price = bookData.Price
	if err := services.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}
	//services.DB.Debug().Model(&models.User{}).Related()
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   book,
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
