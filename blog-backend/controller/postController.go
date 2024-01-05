package controller

import (
	"errors"
	"fmt"
	"github/mgabrielo/React-Golang-Blog/database"
	"github/mgabrielo/React-Golang-Blog/models"
	"github/mgabrielo/React-Golang-Blog/util"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreatePost(c *fiber.Ctx) error {
	var blogpost models.Blog
	if err := c.BodyParser(&blogpost); err != nil {
		fmt.Println("unable to parse body")
	}
	if err := database.DB.Create(&blogpost).Error; err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"message": "Invalid Payload"})
	}
	return c.JSON(fiber.Map{"message": "Post Created Successfuly"})
}

func AllPost(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit := 5
	offset := (page - 1) * limit
	var total int64
	var getBlog []models.Blog
	database.DB.Preload("User").Offset(offset).Limit(limit).Find(&getBlog)
	database.DB.Model(&models.Blog{}).Count(&total)
	return c.JSON(fiber.Map{
		"data": getBlog,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": float64(int(total) / limit),
		},
	})
}

func DetailPost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var blogpost models.Blog
	database.DB.Where("id=?", id).Preload("User").First(&blogpost)
	return c.JSON(fiber.Map{"data": blogpost})
}

func UpdatePost(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var blog models.Blog
	if err := database.DB.First(&blog, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blog not found"})
	}

	// Parse and bind the request body to the blog model
	if err := c.BodyParser(&blog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unable to parse request body"})
	}

	// Update the blog in the database
	if err := database.DB.Save(&blog).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to update blog"})
	}

	return c.JSON(fiber.Map{"message": "Blog updated successfully"})
}

func UniquePost(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	id, _ := util.ParseJWT(cookie)
	var blog []models.Blog
	database.DB.Model(&blog).Where("user_id=?", id).Preload("User").Find(&blog)
	return c.JSON(fiber.Map{"data": blog})
}

func DeletePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	blog := models.Blog{
		Id: uint(id),
	}
	delete_query := database.DB.Delete(&blog)
	if errors.Is(delete_query.Error, gorm.ErrRecordNotFound) {
		c.Status(400)
		return c.JSON(fiber.Map{"message": "Post Not Found"})
	}
	return c.JSON(fiber.Map{"message": "BlogPost Deleted successfully"})
}
