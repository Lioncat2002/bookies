package controllers

import (
	"backend/go-catbox"
	"backend/models"
	"backend/services"
	"backend/utils/token"
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateBookdata struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"desc"`
	Price       float32 `json:"price" binding:"required"`
	Tag         string  `json:"tag" binding:"required"`
}

type BookData struct {
	BookID string `json:"book_id" binding:"required"`
}

func SearchBook(c *gin.Context) {
	name := c.Param("name")
	book := models.Book{}
	if err := services.DB.Where("name = ?", name).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   book,
	})
}

func BuyBook(c *gin.Context) {
	id, err := token.ExtractID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := models.User{}
	if err := services.DB.Preload("Carts").Preload("Owns").Where("id = ?", id).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var totalPrice float32 = 0.0
	for _, book := range user.Carts {
		totalPrice += book.Price
	}

	coins := user.Coins - totalPrice
	if coins < 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "insufficient balance",
		})
		return
	}

	// Deduct coins from user's balance
	if err := services.DB.Model(&user).Where("id = ?", id).Update("coins", coins).Error; err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Move books from cart to user's owned books
	if err := services.DB.Model(&user).Where("id = ?", id).Association("Owns").Append(&user.Carts); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Clear user's cart
	if err := services.DB.Model(&user).Association("Carts").Clear(); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
	})
}

func RateBook(c *gin.Context) {
	bookData := BookData{}
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
		return
	}
	user := models.User{}
	if err := services.DB.Where("id = ?", id).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func AddToCart(c *gin.Context) {
	cartData := BookData{}
	if err := c.ShouldBindJSON(&cartData); err != nil {
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
		return
	}
	user := models.User{}
	if err := services.DB.Where("id = ?", id).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	book := models.Book{}
	if err := services.DB.Where("id = ?", cartData.BookID).Find(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := services.DB.Where("id = ?", id).Find(&user).Association("Carts").Append(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
	})
}

func AddBookUrl(c *gin.Context) {
	f, _ := c.FormFile("file")
	file, _ := f.Open()
	defer file.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id := c.Param("id")
	user := models.User{}
	if err := services.DB.Where("id = ?", id).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if user.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "incorrect role",
		})
		return
	}
	url, err := catbox.New(nil).Upload(buf.Bytes(), string(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	book := models.Book{}
	if err := services.DB.Where("id = ?", id).Find(&book).Update("url", url).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   book,
	})
}

func CreateBook(c *gin.Context) {
	var bookData CreateBookdata
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
		return
	}
	user := models.User{}
	if err := services.DB.Where("id = ?", id).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if user.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "incorrect permissions",
		})
		return
	}
	book := models.Book{}
	book.UserID = id
	book.Name = bookData.Name
	book.Description = bookData.Description
	book.Tag = bookData.Tag
	book.Price = bookData.Price
	if err := services.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   book,
	})
}

func GetOneBook(c *gin.Context) {
	id := c.Param("id")
	book := models.Book{}
	if err := services.DB.Where("id = ?", id).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   book,
	})
}

func AllBooks(c *gin.Context) {
	var books []models.Book
	if err := services.DB.Find(&books).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   books,
	})
}
